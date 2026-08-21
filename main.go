package main

import (
	"embed"
	"log"

	"github.com/selyusize/my-home/internal/app"

	"github.com/joho/godotenv"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	_ = godotenv.Load()

	home, err := app.New(assets)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := home.Close(); err != nil {
			log.Printf("close app: %v", err)
		}
	}()

	if err := home.Run(); err != nil {
		log.Fatal(err)
	}
}
