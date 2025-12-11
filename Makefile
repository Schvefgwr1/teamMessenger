.PHONY: help up build-up debug debug-build test test-full down debug-down

# Variables
COMPOSE_FILE=docker-compose.yml
COMPOSE_DEBUG=docker-compose.debug.yml
ENV_FILE=compose.env
SERVICES=apiService chatService fileService taskService userService notificationService

# Show help
help:
	@echo "Team Messenger - Launch Commands:"
	@echo ""
	@echo "Basic commands:"
	@echo "  make up          - Start system without tests"
	@echo "  make build-up    - Rebuild and start system"
	@echo ""
	@echo "Debug:"
	@echo "  make debug       - Start in debug mode"
	@echo "  make debug-build - Rebuild and start in debug mode"
	@echo ""
	@echo "Testing:"
	@echo "  make test        - Run unit tests for all services and start system"
	@echo "  make test-full   - Full testing (unit + integration tests + start)"
	@echo ""
	@echo "Stop:"
	@echo "  make down       - Stop all containers"
	@echo "  make debug-down - Stop all containers (debug mode)"

# 1. Basic system startup without tests
up:
	@echo "Starting system..."
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d
	@echo "System started!"

# 2. Start with rebuild
build-up:
	@echo "Rebuilding images..."
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) build
	@echo "Starting system..."
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d
	@echo "System rebuilt and started!"

# 3. Start in debug mode
debug:
	@echo "Starting in debug mode..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEBUG) --env-file $(ENV_FILE) up -d
	@echo "System started in debug mode!"

# 4. Start in debug mode with rebuild
debug-build:
	@echo "Rebuilding images for debug..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEBUG) --env-file $(ENV_FILE) build
	@echo "Starting in debug mode..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEBUG) --env-file $(ENV_FILE) up -d
	@echo "System rebuilt and started in debug mode!"

# Helper targets for unit tests per service
test-apiService:
	@echo "Running unit tests for apiService..."
	-cd apiService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for apiService completed"

test-chatService:
	@echo "Running unit tests for chatService..."
	-cd chatService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for chatService completed"

test-fileService:
	@echo "Running unit tests for fileService..."
	-cd fileService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for fileService completed"

test-taskService:
	@echo "Running unit tests for taskService..."
	-cd taskService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for taskService completed"

test-userService:
	@echo "Running unit tests for userService..."
	-cd userService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for userService completed"

test-notificationService:
	@echo "Running unit tests for notificationService..."
	-cd notificationService && go test -short ./tests/controllers/... ./tests/handlers/...
	@echo "Unit tests for notificationService completed"

# 5. Start with unit tests
test: test-apiService test-chatService test-fileService test-taskService test-userService test-notificationService
	@echo ""
	@echo "Starting system..."
	@$(MAKE) up
	@echo "Unit tests passed, system started!"

# Helper targets for integration tests per service
integration-apiService:
	@echo "Running integration tests for apiService..."
	-cd apiService && $(MAKE) integration

integration-chatService:
	@echo "Running integration tests for chatService..."
	-cd chatService && $(MAKE) integration

integration-fileService:
	@echo "Running integration tests for fileService..."
	-cd fileService && $(MAKE) integration

integration-taskService:
	@echo "Running integration tests for taskService..."
	-cd taskService && $(MAKE) integration

integration-userService:
	@echo "Running integration tests for userService..."
	-cd userService && $(MAKE) integration

integration-notificationService:
	@echo "Running integration tests for notificationService..."
	-cd notificationService && $(MAKE) integration

# 6. Full testing (recommended production option)
test-full: test-apiService test-chatService test-fileService test-taskService test-userService test-notificationService
	@echo ""
	@echo "Step 2: Integration tests for all services..."
	@$(MAKE) integration-apiService integration-chatService integration-fileService integration-taskService integration-userService integration-notificationService
	@echo ""
	@echo "Step 3: Starting system..."
	@$(MAKE) up
	@echo ""
	@echo "Full testing completed, system started!"

# Stop system
down:
	@echo "Stopping system..."
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down
	@echo "System stopped!"

# Stop system in debug mode
debug-down:
	@echo "Stopping system in debug mode..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEBUG) --env-file $(ENV_FILE) down
	@echo "System stopped!"
