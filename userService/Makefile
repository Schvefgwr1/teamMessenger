.PHONY: up-integration down-integration migrate-integration run-tests-integration test-integration logs-integration integration

COMPOSE=docker-compose -f tests/integration/docker-compose.integration.yml
# Extra args for go test (default -v). Override with: make integration TEST_ARGS="-v -run TestX"
TEST_ARGS?=-v
# Cross-platform env injection: use PowerShell so that it works on Windows shells too.
TEST_CMD=powershell -Command "$$env:DB_HOST='localhost'; $$env:DB_PORT='5433'; $$env:DB_USER='postgres'; $$env:DB_PASSWORD='postgres'; $$env:DB_NAME='team_messenger_test'; $$env:KAFKA_BROKERS='localhost:9092'; $$env:KAFKA_TOPIC_NOTIFICATIONS='notifications_test'; $$env:KAFKA_TOPIC_KEY_UPDATES='key_updates_test'; go test -tags=integration $(TEST_ARGS) ./tests/integration/..."

up-integration:
	@echo "==> Bringing up infra (Postgres, Kafka, Zookeeper)..."
	$(COMPOSE) up -d --wait
	@echo "==> Infra is up."

down-integration:
	@echo "==> Stopping infra and cleaning up..."
	$(COMPOSE) down -v
	@echo "==> Infra stopped."

migrate-integration:
	@echo "==> Applying migrations..."
	$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U postgres -d team_messenger_test -c "DROP SCHEMA IF EXISTS user_service CASCADE;"
	$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U postgres -d team_messenger_test -f /migrations/000001_init_db_v1.up.sql
	$(COMPOSE) exec -T postgres psql -U postgres -d team_messenger_test -f /migrations/000002_add_new_permissions.up.sql
	@echo "==> Migrations applied."

run-tests-integration:
	@echo "==> Running integration tests (cmd: $(TEST_CMD))"
	$(TEST_CMD)

test-integration: up-integration migrate-integration run-tests-integration

logs-integration:
	$(COMPOSE) logs -f postgres kafka

# Single entrypoint: up -> migrate -> test -> down
integration:
	@$(MAKE) up-integration
	@$(MAKE) migrate-integration
	@$(MAKE) run-tests-integration
	@echo "==> Tearing down infra..."
	@$(MAKE) down-integration
	@echo "==> Tests finished."

