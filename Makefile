SHELL := /bin/bash
.PHONY: all api web desktop migrate dev-api dev-web dev-api-web dev-desktop clean test deps

# Build output directories
BUILD_DIR := build
API_BIN := $(BUILD_DIR)/asset-manager-api
MIGRATE_BIN := $(BUILD_DIR)/asset-manager-migrate
HASHPW_BIN := $(BUILD_DIR)/asset-manager-hashpw

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod

# Default target
all: deps api migrate web desktop

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	cd web && npm install

# Build API server
api:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(API_BIN) ./cmd/api

# Build migration tool
migrate:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(MIGRATE_BIN) ./cmd/migrate

# Build password hash tool
hashpw:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(HASHPW_BIN) ./cmd/hashpw

# Build web frontend
web:
	cd web && npm run build

# Build desktop app
desktop:
	@cd desktop && \
	DISTRO=$$(lsb_release -is 2>/dev/null || echo "Unknown"); \
	VERSION=$$(lsb_release -rs 2>/dev/null || echo "0"); \
	if [[ ("$$DISTRO" == "Ubuntu" && "$$(printf '%s\n' "$$VERSION" "24.00" | sort -V | head -n1)" == "24.00") || \
	  ("$$DISTRO" == "Elementary" && "$$(printf '%s\n' "$$VERSION" "8" | sort -V | head -n1)" == "8") || \
	  ("$$DISTRO" == "Linuxmint" && "$$(printf '%s\n' "$$VERSION" "22.0" | sort -V | head -n1)" == "22.0") ]]; then \
	  echo "Detected Ubuntu >= 24, Elementary >= 8, or Linux Mint >= 22. Using webkit2_41..."; \
	  wails build -tags webkit2_41; \
	else \
	  echo "Using standard webkit build..."; \
	  wails build; \
	fi

# Run database migrations
run-migrate: migrate
	$(MIGRATE_BIN) -config config.yaml -migrations migrations

# Development: run API server
dev-api:
	$(GOCMD) run ./cmd/api -config config.yaml

# Development: run web frontend
dev-web:
	cd web && npm run dev

# Development: run API and web together
dev-api-web:
	@echo "Starting API and Web servers..."
	@$(GOCMD) run ./cmd/api -config config.yaml & API_PID=$$!; \
	cd web && npm run dev -- --host; \
	kill $$API_PID 2>/dev/null || true

# Development: run desktop app
dev-desktop:
	@cd desktop && \
	DISTRO=$$(lsb_release -is 2>/dev/null || echo "Unknown"); \
	VERSION=$$(lsb_release -rs 2>/dev/null || echo "0"); \
	if [[ ("$$DISTRO" == "Ubuntu" && "$$(printf '%s\n' "$$VERSION" "24.00" | sort -V | head -n1)" == "24.00") || \
	  ("$$DISTRO" == "Elementary" && "$$(printf '%s\n' "$$VERSION" "8" | sort -V | head -n1)" == "8") || \
	  ("$$DISTRO" == "Linuxmint" && "$$(printf '%s\n' "$$VERSION" "22.0" | sort -V | head -n1)" == "22.0") ]]; then \
	  echo "Detected Ubuntu >= 24, Elementary >= 8, or Linux Mint >= 22. Using webkit2_41..."; \
	  wails dev -tags webkit2_41; \
	else \
	  echo "Using standard webkit build..."; \
	  wails dev; \
	fi

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	cd web && rm -rf dist node_modules
	cd desktop && rm -rf build
	cd desktop/frontend && rm -rf dist node_modules

# Install Wails CLI (if not installed)
install-wails:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Generate a secure JWT secret
gen-secret:
	@openssl rand -base64 32

# Help
help:
	@echo "Asset Manager Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          Build everything (api, migrate, web, desktop)"
	@echo "  deps         Download and install dependencies"
	@echo "  api          Build API server"
	@echo "  migrate      Build migration tool"
	@echo "  web          Build web frontend"
	@echo "  desktop      Build desktop app"
	@echo "  run-migrate  Run database migrations"
	@echo "  dev-api      Run API server in development mode"
	@echo "  dev-web      Run web frontend in development mode"
	@echo "  dev-api-web  Run API and web frontend together"
	@echo "  dev-desktop  Run desktop app in development mode"
	@echo "  test         Run tests"
	@echo "  clean        Clean build artifacts"
	@echo "  gen-secret   Generate a secure JWT secret"
	@echo "  help         Show this help message"
