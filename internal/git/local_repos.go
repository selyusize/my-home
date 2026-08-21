package git

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	pkgdb "github.com/selyusize/my-home/pkg/db"
	pkggit "github.com/selyusize/my-home/pkg/git"
	pkglocal "github.com/selyusize/my-home/pkg/git/local"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
)

const maxLocalPathLength = 4096

// LocalRepoSettings is the persisted folder configuration from the profile page.
type LocalRepoSettings struct {
	SharedPath     string `json:"sharedPath"`
	GitHubSeparate bool   `json:"githubSeparate"`
	GitHubPath     string `json:"githubPath"`
	GitLabSeparate bool   `json:"gitlabSeparate"`
	GitLabPath     string `json:"gitlabPath"`
}

// ScanReport is how many clones Check stored per provider.
type ScanReport struct {
	GitHub int `json:"github"`
	GitLab int `json:"gitlab"`
}

type localRepoSettingsRow struct {
	ID             uint `gorm:"primaryKey"`
	SharedPath     string
	GitHubSeparate bool
	GitHubPath     string
	GitLabSeparate bool
	GitLabPath     string
	UpdatedAt      time.Time
}

func (localRepoSettingsRow) TableName() string {
	return "local_repo_settings"
}

type localCloneRow struct {
	ID       uint   `gorm:"primaryKey"`
	Provider string `gorm:"size:32;index"`
	FullName string `gorm:"size:512;index"`
	Owner    string
	Name     string
	Path     string
}

func (localCloneRow) TableName() string {
	return "local_clones"
}

// LocalReposService stores local clone folders and scan results.
type LocalReposService struct {
	db pkgdb.Client
}

func NewLocalReposService(database pkgdb.Client) (*LocalReposService, error) {
	if database == nil {
		return nil, fmt.Errorf("local repos: database is required")
	}
	service := &LocalReposService{db: database}
	if err := database.AutoMigrate(&localRepoSettingsRow{}, &localCloneRow{}); err != nil {
		return nil, fmt.Errorf("local repos: migrate: %w", err)
	}
	return service, nil
}

func (s *LocalReposService) GetSettings() (LocalRepoSettings, error) {
	row, err := s.loadSettingsRow()
	if err != nil {
		return LocalRepoSettings{}, err
	}
	return settingsFromRow(row), nil
}

func (s *LocalReposService) SaveSettings(settings LocalRepoSettings) error {
	normalized, err := normalizeSettings(settings)
	if err != nil {
		return err
	}
	return s.persistSettings(normalized)
}

func (s *LocalReposService) SelectDirectory(title string) (string, error) {
	app := application.Get()
	if app == nil {
		return "", fmt.Errorf("local repos: app is not running")
	}
	if strings.TrimSpace(title) == "" {
		title = "Выберите папку"
	}

	dialog := app.Dialog.OpenFile().
		SetTitle(title).
		CanChooseDirectories(true).
		CanChooseFiles(false)
	if win := app.Window.Current(); win != nil {
		dialog.AttachToWindow(win)
	}

	path, err := dialog.PromptForSingleSelection()
	if err != nil || strings.TrimSpace(path) == "" {
		return "", nil
	}
	return filepath.Clean(path), nil
}

func (s *LocalReposService) Check(settings LocalRepoSettings) (ScanReport, error) {
	normalized, err := normalizeSettings(settings)
	if err != nil {
		return ScanReport{}, err
	}
	if err := s.persistSettings(normalized); err != nil {
		return ScanReport{}, err
	}

	clones, err := collectClones(normalized)
	if err != nil {
		return ScanReport{}, err
	}
	if err := s.replaceClones(clones); err != nil {
		return ScanReport{}, err
	}

	report := ScanReport{}
	for _, clone := range clones {
		switch clone.Provider {
		case string(pkggit.ProviderGitHub):
			report.GitHub++
		case string(pkggit.ProviderGitLab):
			report.GitLab++
		}
	}
	return report, nil
}

func (s *LocalReposService) ListClones(provider string) ([]pkglocal.Clone, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	query := s.db.DB().Order("name").Order("path")
	switch provider {
	case "":
	case string(pkggit.ProviderGitHub), string(pkggit.ProviderGitLab):
		query = query.Where("provider = ?", provider)
	default:
		return nil, fmt.Errorf("local repos: unknown provider")
	}

	var rows []localCloneRow
	if err := query.Find(&rows).Error; err != nil {
		if provider == "" {
			return nil, fmt.Errorf("local repos: list: %w", err)
		}
		return nil, fmt.Errorf("local repos: list %s: %w", provider, err)
	}

	clones := make([]pkglocal.Clone, 0, len(rows))
	for _, row := range rows {
		clones = append(clones, cloneFromRow(row))
	}
	return clones, nil
}

func (s *LocalReposService) loadSettingsRow() (localRepoSettingsRow, error) {
	var row localRepoSettingsRow
	if err := s.db.DB().Where("id = ?", 1).Limit(1).Find(&row).Error; err != nil {
		return localRepoSettingsRow{}, fmt.Errorf("local repos: load settings: %w", err)
	}
	return row, nil
}

func (s *LocalReposService) persistSettings(settings LocalRepoSettings) error {
	row := localRepoSettingsRow{
		ID:             1,
		SharedPath:     settings.SharedPath,
		GitHubSeparate: settings.GitHubSeparate,
		GitHubPath:     settings.GitHubPath,
		GitLabSeparate: settings.GitLabSeparate,
		GitLabPath:     settings.GitLabPath,
	}
	if err := s.db.DB().Save(&row).Error; err != nil {
		return fmt.Errorf("local repos: save settings: %w", err)
	}
	return nil
}

func (s *LocalReposService) replaceClones(clones []pkglocal.Clone) error {
	rows := make([]localCloneRow, 0, len(clones))
	for _, clone := range clones {
		rows = append(rows, localCloneRow{
			Provider: clone.Provider,
			FullName: clone.FullName,
			Owner:    clone.Owner,
			Name:     clone.Name,
			Path:     clone.Path,
		})
	}

	err := s.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&localCloneRow{}).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		return tx.Create(&rows).Error
	})
	if err != nil {
		return fmt.Errorf("local repos: save clones: %w", err)
	}
	return nil
}

func cloneFromRow(row localCloneRow) pkglocal.Clone {
	return pkglocal.Clone{
		Provider: row.Provider,
		Owner:    row.Owner,
		Name:     row.Name,
		FullName: row.FullName,
		Path:     row.Path,
		HTMLURL:  pkglocal.HTMLURL(row.Provider, row.FullName),
	}
}

func collectClones(settings LocalRepoSettings) ([]pkglocal.Clone, error) {
	roots := settingsRoots(settings)
	found := make([]pkglocal.Clone, 0)
	for _, root := range roots {
		clones, err := pkglocal.Scan(root)
		if err != nil {
			return nil, err
		}
		found = append(found, clones...)
	}

	uniq := make(map[string]pkglocal.Clone)
	order := make([]string, 0)
	for _, clone := range found {
		if !belongsToAny(clone.Path, roots) {
			continue
		}
		key := cloneKey(clone)
		if _, exists := uniq[key]; exists {
			continue
		}
		uniq[key] = clone
		order = append(order, key)
	}

	result := make([]pkglocal.Clone, 0, len(order))
	for _, key := range order {
		result = append(result, uniq[key])
	}
	return result, nil
}

func settingsRoots(settings LocalRepoSettings) []string {
	githubRoot := settings.SharedPath
	if settings.GitHubSeparate && settings.GitHubPath != "" {
		githubRoot = settings.GitHubPath
	}
	gitlabRoot := settings.SharedPath
	if settings.GitLabSeparate && settings.GitLabPath != "" {
		gitlabRoot = settings.GitLabPath
	}

	roots := make([]string, 0, 2)
	seen := make(map[string]struct{})
	for _, root := range []string{githubRoot, gitlabRoot} {
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}
	return roots
}

func belongsToAny(path string, roots []string) bool {
	for _, root := range roots {
		if pathBelongsTo(path, root) {
			return true
		}
	}
	return false
}

func cloneKey(clone pkglocal.Clone) string {
	if clone.Provider != "" && clone.FullName != "" {
		return clone.Provider + "\x00" + strings.ToLower(clone.FullName)
	}
	return clone.Path
}

func pathBelongsTo(child, root string) bool {
	rel, err := filepath.Rel(root, child)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, "..")
}

func settingsFromRow(row localRepoSettingsRow) LocalRepoSettings {
	return LocalRepoSettings{
		SharedPath:     row.SharedPath,
		GitHubSeparate: row.GitHubSeparate,
		GitHubPath:     row.GitHubPath,
		GitLabSeparate: row.GitLabSeparate,
		GitLabPath:     row.GitLabPath,
	}
}

func normalizeSettings(settings LocalRepoSettings) (LocalRepoSettings, error) {
	shared, err := normalizePath(settings.SharedPath)
	if err != nil {
		return LocalRepoSettings{}, err
	}
	githubPath, err := normalizePath(settings.GitHubPath)
	if err != nil {
		return LocalRepoSettings{}, err
	}
	gitlabPath, err := normalizePath(settings.GitLabPath)
	if err != nil {
		return LocalRepoSettings{}, err
	}
	return LocalRepoSettings{
		SharedPath:     shared,
		GitHubSeparate: settings.GitHubSeparate,
		GitHubPath:     githubPath,
		GitLabSeparate: settings.GitLabSeparate,
		GitLabPath:     gitlabPath,
	}, nil
}

func normalizePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", nil
	}
	if strings.ContainsRune(path, 0) {
		return "", fmt.Errorf("local repos: invalid path")
	}
	if utf8.RuneCountInString(path) > maxLocalPathLength {
		return "", fmt.Errorf("local repos: path is too long")
	}
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("local repos: path must be absolute")
	}
	return filepath.Clean(path), nil
}
