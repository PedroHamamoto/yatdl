# Database connection
DB_URL=postgres://postgres:postgres@db:5432/yatdl_db?sslmode=disable
DOCKER_NETWORK=yatdl_default
MIGRATIONS_PATH=$(PWD)/migrations

# Migration base command
MIGRATE=docker run --rm -v $(MIGRATIONS_PATH):/migrations --network $(DOCKER_NETWORK) migrate/migrate -path=/migrations/ -database "$(DB_URL)"

# Run all pending migrations
migrate-up:
	$(MIGRATE) up

# Rollback the last migration
migrate-down:
	$(MIGRATE) down 1

# Rollback all migrations
migrate-down-all:
	$(MIGRATE) down -all

# Check current migration version
migrate-version:
	$(MIGRATE) version

# Create a new migration file
# Usage: make migrate-create NAME=create_todos_table
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=create_todos_table"; \
		exit 1; \
	fi
	@timestamp=$$(date +%s); \
	touch $(MIGRATIONS_PATH)/$${timestamp}_$(NAME).up.sql; \
	touch $(MIGRATIONS_PATH)/$${timestamp}_$(NAME).down.sql; \
	echo "Created migration files:"; \
	echo "  $(MIGRATIONS_PATH)/$${timestamp}_$(NAME).up.sql"; \
	echo "  $(MIGRATIONS_PATH)/$${timestamp}_$(NAME).down.sql"

# Drop everything in the database (DANGEROUS - use with caution)
migrate-drop:
	@echo "WARNING: This will drop all tables in the database!"
	@echo "Press Ctrl+C to cancel or Enter to continue..."; \
	read confirm; \
	$(MIGRATE) drop

api-start:
	go run cmd/api/main.go

.PHONY: migrate migrate-up migrate-down migrate-down-all migrate-version migrate-force migrate-create migrate-drop