# Makefile for Backuper Project

.PHONY: up down migrate clean restore help

# Start containers
up:
	docker-compose up -d --build

# Stop containers
down:
	docker-compose down

# Create a test table and insert some dummy data
# Usage: make migrate
migrate:
	@echo "Migrating test data to database..."
	docker exec -i test_db psql -U user -d test_db -c "CREATE TABLE IF NOT EXISTS v1_statistic (id SERIAL PRIMARY KEY, data TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
	docker exec -i test_db psql -U user -d test_db -c "INSERT INTO v1_statistic (data) VALUES ('test data 1'), ('test data 2'), ('test data 3');"
	docker exec -i test_db psql -U user -d test_db -c "\dt"
	@echo "Migration completed."

# Restore a dump from a file
# Usage: make restore file=test_db/your_file.sql.gz
restore:
	@if [ -z "$(file)" ]; then \
		echo "Error: Please specify the file parameter."; \
		echo "Usage: make restore file=test_db/filename.sql[.gz]"; \
		exit 1; \
	fi
	@echo "Restoring from ./backups/$(file)..."
	@if echo "$(file)" | grep -q "\.gz$$"; then \
		gunzip -c ./backups/$(file) | docker exec -i test_db psql -U user -d test_db; \
	else \
		cat ./backups/$(file) | docker exec -i test_db psql -U user -d test_db; \
	fi
	@echo "Restore completed successfully."

# Remove logs and backups
clean:
	rm -rf ./logs/* ./backups/*

help:
	@echo "Available commands:"
	@echo "  make up      - Start the backuper and postgres containers"
	@echo "  make down    - Stop the containers"
	@echo "  make migrate - Run test migrations (inject dummy data)"
	@echo "  make restore file=path/to/file - Restore database from a dump"
	@echo "  make clean   - Clear logs and backups"
