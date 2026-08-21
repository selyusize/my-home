package app

import (
	"embed"
	"fmt"
	"time"

	homedb "github.com/selyusize/my-home/internal/db"
	"github.com/selyusize/my-home/pkg/db"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	application.RegisterEvent[string]("time")
}

type App struct {
	database db.Client
	wails    *application.App
}

func New(assets embed.FS) (*App, error) {
	database, err := homedb.New()
	if err != nil {
		return nil, err
	}

	services, err := newServices(database)
	if err != nil {
		_ = database.Close()
		return nil, err
	}

	wailsApp := newWails(assets, services)
	app := &App{database: database, wails: wailsApp}
	app.startClock()
	return app, nil
}

func (a *App) Close() error {
	if a == nil || a.database == nil {
		return nil
	}
	return a.database.Close()
}

func (a *App) Run() error {
	if a == nil || a.wails == nil {
		return fmt.Errorf("app: not initialized")
	}
	return a.wails.Run()
}

func (a *App) startClock() {
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			a.wails.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()
}
