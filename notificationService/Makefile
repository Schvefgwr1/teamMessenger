.PHONY: up-integration down-integration run-tests-integration test-integration logs-integration integration

COMPOSE=docker-compose -f tests/integration/docker-compose.integration.yml
# Extra args for go test (default -v). Override with: make integration TEST_ARGS="-v -run TestX"
TEST_ARGS?=-v
# Cross-platform env injection: use PowerShell so that it works on Windows shells too.
TEST_CMD=powershell -Command "$$env:KAFKA_BROKERS='localhost:9092'; $$env:KAFKA_TOPIC_NOTIFICATIONS='notifications_test'; $$env:SMTP_HOST='localhost'; $$env:SMTP_PORT='1025'; $$env:SMTP_USERNAME='test@example.com'; $$env:SMTP_PASSWORD='testpassword'; $$env:TEMPLATE_PATH='../../templates'; go test -tags=integration $(TEST_ARGS) ./tests/integration/..."

up-integration:
	@echo "==> Bringing up infra (Kafka, Zookeeper, MailHog)..."
	$(COMPOSE) up -d --wait
	@echo "==> Infra is up."

down-integration:
	@echo "==> Stopping infra and cleaning up..."
	$(COMPOSE) down -v
	@echo "==> Infra stopped."

run-tests-integration:
	@echo "==> Running integration tests (cmd: $(TEST_CMD))"
	$(TEST_CMD)

test-integration: up-integration run-tests-integration

logs-integration:
	$(COMPOSE) logs -f kafka mailhog

# Single entrypoint: up -> test -> down
integration:
	@$(MAKE) up-integration
	@$(MAKE) run-tests-integration
	@echo "==> Tearing down infra..."
	@$(MAKE) down-integration
	@echo "==> Tests finished."

