.PHONY: lint lint-go lint-frontend dev build

.DEFAULT_GOAL := help

help:
	@echo "make lint   - Go vet + frontend ESLint/FSD"
	@echo "make dev    - run the app in development mode"
	@echo "make build  - production build"

lint: lint-go lint-frontend

lint-go:
	mkdir -p frontend/dist
	touch frontend/dist/.gitkeep
	go vet .

lint-frontend:
	pnpm --dir frontend lint

dev:
	wails3 dev -config ./build/config.yml -port $${WAILS_VITE_PORT:-9245}

build:
	PACKAGE_MANAGER=pnpm task build
