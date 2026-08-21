package app

import (
	"github.com/selyusize/my-home/internal/bitrix"
	"github.com/selyusize/my-home/internal/dl"
	"github.com/selyusize/my-home/internal/docker"
	"github.com/selyusize/my-home/internal/git"
	"github.com/selyusize/my-home/internal/settings"
	"github.com/selyusize/my-home/internal/window"
	"github.com/selyusize/my-home/pkg/db"
)

type services struct {
	github     *git.GitHubService
	gitlab     *git.GitLabService
	bitrix     *bitrix.BitrixService
	docker     *docker.DockerService
	dl         *dl.DLService
	window     *window.WindowService
	settings   *settings.SettingsService
	localRepos *git.LocalReposService
}

func newServices(database db.Client) (*services, error) {
	gitCreds, err := git.NewCredentialStore(database)
	if err != nil {
		return nil, err
	}
	bitrixCreds, err := bitrix.NewCredentialStore(database)
	if err != nil {
		return nil, err
	}
	settingsService, err := settings.NewSettingsService(database)
	if err != nil {
		return nil, err
	}
	localReposService, err := git.NewLocalReposService(database)
	if err != nil {
		return nil, err
	}
	githubService, err := git.NewGitHubService(gitCreds)
	if err != nil {
		return nil, err
	}
	gitlabService, err := git.NewGitLabService(gitCreds)
	if err != nil {
		return nil, err
	}
	bitrixService, err := bitrix.NewBitrixService(bitrixCreds)
	if err != nil {
		return nil, err
	}
	dockerService, err := docker.NewDockerService()
	if err != nil {
		return nil, err
	}
	dlService, err := dl.NewDLService(dockerService)
	if err != nil {
		return nil, err
	}

	return &services{
		github:     githubService,
		gitlab:     gitlabService,
		bitrix:     bitrixService,
		docker:     dockerService,
		dl:         dlService,
		window:     window.NewWindowService(),
		settings:   settingsService,
		localRepos: localReposService,
	}, nil
}
