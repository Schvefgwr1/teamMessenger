.PHONY: integration integration-up integration-down test-integration

COMPOSE_FILE=tests/integration/docker-compose.integration.yml
TEST_ARGS?=

integration-up:
	docker compose -f $(COMPOSE_FILE) up -d

integration-down:
	docker compose -f $(COMPOSE_FILE) down -v

test-integration:
	go test -tags=integration -v ./tests/integration/... $(TEST_ARGS)

integration: integration-up
	$(MAKE) test-integration TEST_ARGS=$(TEST_ARGS)
	$(MAKE) integration-down

