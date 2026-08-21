package git

import (
	"fmt"
	"strings"
	"time"

	pkgdb "github.com/selyusize/my-home/pkg/db"
	pkggit "github.com/selyusize/my-home/pkg/git"

	"gorm.io/gorm"
)

type gitCredential struct {
	Provider  string `gorm:"primaryKey;size:32"`
	Token     string `gorm:"not null"`
	UpdatedAt time.Time
}

func (gitCredential) TableName() string {
	return "git_credentials"
}

type CredentialStore struct {
	db pkgdb.Client
}

func NewCredentialStore(database pkgdb.Client) (*CredentialStore, error) {
	if database == nil {
		return nil, fmt.Errorf("git credentials: database is required")
	}
	store := &CredentialStore{db: database}
	if err := migrateGitCredentials(database.DB()); err != nil {
		return nil, fmt.Errorf("git credentials: migrate: %w", err)
	}
	return store, nil
}

func migrateGitCredentials(gdb *gorm.DB) error {
	m := gdb.Migrator()
	if m.HasTable(&gitCredential{}) && m.HasColumn(&gitCredential{}, "updated_at") {
		cols, err := m.ColumnTypes(&gitCredential{})
		if err != nil {
			return err
		}
		for _, col := range cols {
			if !strings.EqualFold(col.Name(), "updated_at") {
				continue
			}
			switch strings.ToLower(col.DatabaseTypeName()) {
			case "bigint", "int8", "integer", "int4", "int":
				if err := m.DropColumn(&gitCredential{}, "UpdatedAt"); err != nil {
					return err
				}
			}
			break
		}
	}
	return gdb.AutoMigrate(&gitCredential{})
}

func (s *CredentialStore) Load(provider pkggit.Provider) (string, error) {
	var row gitCredential
	if err := s.db.DB().Where("provider = ?", string(provider)).Limit(1).Find(&row).Error; err != nil {
		return "", fmt.Errorf("git credentials: load %s: %w", provider, err)
	}
	return row.Token, nil
}

func (s *CredentialStore) Save(provider pkggit.Provider, token string) error {
	row := gitCredential{Provider: string(provider), Token: token}
	if err := s.db.DB().Save(&row).Error; err != nil {
		return fmt.Errorf("git credentials: save %s: %w", provider, err)
	}
	return nil
}

func (s *CredentialStore) Delete(provider pkggit.Provider) error {
	if err := s.db.DB().Where("provider = ?", string(provider)).Delete(&gitCredential{}).Error; err != nil {
		return fmt.Errorf("git credentials: delete %s: %w", provider, err)
	}
	return nil
}
