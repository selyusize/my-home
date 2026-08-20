.PHONY: lint lint-go lint-frontend dev build sonar-up sonar-down sonar-token sonar-coverage sonar-scan sonar

.DEFAULT_GOAL := help

help:
	@echo "make lint            - Go vet + frontend ESLint/FSD (+ sonarjs, без токена)"
	@echo "make dev             - run the app in development mode"
	@echo "make build           - production build"
	@echo "make sonar-up        - start free local SonarQube (Docker)"
	@echo "make sonar-down      - stop SonarQube"
	@echo "make sonar-token     - create free local token → .env.sonar"
	@echo "make sonar-coverage  - generate Go coverage.out"
	@echo "make sonar-scan      - scan Go + frontend (сервер должен быть up)"
	@echo "make sonar           - up + scan + down (контейнеры гасятся после)"

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

sonar-up:
	docker compose -f docker-compose.sonarqube.yml up -d

sonar-down:
	docker compose -f docker-compose.sonarqube.yml down

sonar-token:
	chmod +x scripts/sonar-token.sh
	./scripts/sonar-token.sh

sonar-coverage:
	mkdir -p frontend/dist
	touch frontend/dist/.gitkeep
	go test ./... -coverprofile=coverage.out -covermode=atomic || true

sonar-scan: sonar-coverage
	@if [ -z "$${SONAR_TOKEN}" ] && [ -f .env.sonar ]; then set -a; . ./.env.sonar; set +a; fi; \
	if [ -z "$${SONAR_TOKEN}" ]; then \
	  echo "Нет токена. Запусти: make sonar-token"; \
	  exit 1; \
	fi; \
	docker run --rm \
		--network host \
		-e SONAR_HOST_URL=$${SONAR_HOST_URL:-http://localhost:9000} \
		-e SONAR_TOKEN=$${SONAR_TOKEN} \
		-v "$(CURDIR):/usr/src" \
		sonarsource/sonar-scanner-cli

# Always tear down Sonar containers after the run (success or failure).
sonar:
	@set -e; \
	$(MAKE) sonar-up; \
	$(MAKE) sonar-token; \
	status=0; \
	$(MAKE) sonar-scan || status=$$?; \
	$(MAKE) sonar-down; \
	exit $$status
