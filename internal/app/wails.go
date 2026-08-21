package app

import (
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func newWails(assets embed.FS, svc *services) *application.App {
	app := application.New(application.Options{
		Name:        "myapp",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewServiceWithOptions(svc.github, application.ServiceOptions{
				Name: "GitHubService",
			}),
			application.NewServiceWithOptions(svc.gitlab, application.ServiceOptions{
				Name: "GitLabService",
			}),
			application.NewServiceWithOptions(svc.bitrix, application.ServiceOptions{
				Name: "BitrixService",
			}),
			application.NewServiceWithOptions(svc.docker, application.ServiceOptions{
				Name: "DockerService",
			}),
			application.NewServiceWithOptions(svc.dl, application.ServiceOptions{
				Name: "DLService",
			}),
			application.NewServiceWithOptions(svc.window, application.ServiceOptions{
				Name: "WindowService",
			}),
			application.NewServiceWithOptions(svc.settings, application.ServiceOptions{
				Name: "SettingsService",
			}),
			application.NewServiceWithOptions(svc.localRepos, application.ServiceOptions{
				Name: "LocalReposService",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	} else {
		file := menu.AddSubmenu("File")
		file.AddRole(application.Quit)
	}
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.ViewMenu)
	menu.AddRole(application.WindowMenu)
	app.Menu.Set(menu)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:              "Window 1",
		Width:              1000,
		Height:             618,
		UseApplicationMenu: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(20, 20, 20),
		URL:              "/",
	})

	return app
}
