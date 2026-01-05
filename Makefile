.PHONY: help up down restart logs shell mysql-shell build run clean

help: ## Show this help
	@echo Available commands:
	@echo   make up          - Start all services
	@echo   make down        - Stop all services
	@echo   make shell       - Enter app container
	@echo   make mysql-shell - Enter MySQL shell
	@echo   make build       - Build the Go binary
	@echo   make run         - Run the app
	@echo   make logs        - Show logs
	@echo   make clean       - Clean everything

up: ## Start all services
	docker-compose up -d
	@echo Starting services...
	@echo MySQL 5.5 is loading 1 million shuffled records...
	@echo This will take 5-10 minutes on first startup!
	@echo.
	@echo Watch progress with: make logs-mysql
	@echo Or check when ready: make check-mysql

down: ## Stop all services
	docker-compose down

restart: down up ## Restart all services

logs: ## Show logs
	docker-compose logs -f

logs-app: ## Show app logs
	docker-compose logs -f app

logs-mysql: ## Show MySQL logs and watch insertion progress
	docker-compose logs -f mysql

shell: ## Enter app container shell
	docker-compose exec app bash

mysql-shell: ## Enter MySQL shell
	docker-compose exec mysql mysql -uroot test

build: ## Build the Go binary
	docker-compose exec app go build -o /bin/goParseInnodb ./cmd/goParseInnodb

run: ## Run the app
	docker-compose exec app /bin/goParseInnodb

quick: build run ## Build and run

test: ## Test the setup
	@echo Testing Go app...
	@docker-compose exec app /bin/goParseInnodb
	@echo.
	@echo Testing MySQL connection...
	@docker-compose exec mysql mysql -uroot test -e "SELECT COUNT(*) as row_count FROM t;"

check-mysql: ## Check if MySQL is ready and data is loaded
	@echo Checking MySQL status...
	@docker-compose exec mysql mysqladmin ping || echo MySQL not ready yet
	@echo.
	@echo Checking row count...
	@docker-compose exec mysql mysql -uroot test -e "SELECT COUNT(*) as rows_inserted FROM t;" || echo Table not ready yet

wait-for-mysql: ## Wait for MySQL to finish loading data
	@echo Waiting for MySQL to complete data loading...
	@echo This checks every 30 seconds until 1 million rows are present.
	@:; while [ "$$(docker-compose exec mysql mysql -uroot test -se 'SELECT COUNT(*) FROM t' 2>/dev/null || echo 0)" != "1000000" ]; do \
		echo "Still loading... Current count: $$(docker-compose exec mysql mysql -uroot test -se 'SELECT COUNT(*) FROM t' 2>/dev/null || echo 0)"; \
		sleep 30; \
	done
	@echo Data loading complete!

clean: ## Clean everything
	docker-compose down -v
	if exist output rmdir /s /q output
	mkdir output

status: ## Show container status
	docker-compose ps

datadir: ## Show MySQL data directory
	docker-compose exec app ls -lah /mysql-data/test/

sample-data: ## Show sample data from table
	docker-compose exec mysql mysql -uroot test -e "SELECT * FROM t LIMIT 20;"
	@echo.
	@echo Showing first 20 rows. Total rows:
	docker-compose exec mysql mysql -uroot test -e "SELECT COUNT(*) as total FROM t;"