# Установите переменную окружения для Go
GO ?= go
# Путь к исходникам
SRC_DIRS := ./...
# Цель по умолчанию
default: run

# Цель для запуска приложения
run:
	@echo "Generating documentation with swag..."
	@swag init -g cmd/main.go
	@echo "Documentation generated."
	@echo "You can see it on http://localhost:8080/swagger/index.html#/"
	@echo "Running Bloomify..."
	@docker-compose up --build

# Цель для проверки стиля кода
lint:
	@echo "Linting code..."
	@$(GO) install golang.org/x/lint/golint@latest
	@golint $(SRC_DIRS)

# Цель для форматирования кода
fmt:
	@echo "Formatting code..."
	@$(GO) fmt $(SRC_DIRS)

# Цель для тестирования кода
test:
	@echo "Testing code..."
	@$(GO) install gotest.tools/gotestsum@latest
	@gotestsum --format testname ./...
	
# Цель для остановки приложения
clean:
	@echo "Stopping Bloomify..."
	@docker-compose down

