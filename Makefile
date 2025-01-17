# Установите переменную окружения для Go
GO ?= go

# Путь к исходникам
SRC_DIRS := ./...

# Цель по умолчанию
default: run

# Цель для сборки приложения
build:
	@echo "Building Bloomify..."
	@mkdir -p bin
	@$(GO) build -o ./bin/bloomify ./cmd
	@echo "Build complete."

# Цель для запуска приложения
run:
	@echo "Generating documentation with swag..."
	@swag init -g cmd/main.go
	@echo "Documentation generated."
	@echo "Running Bloomify..."
	@$(GO) run $(SRC_DIRS)

# Цель для проверки стиля кода
lint:
	@echo "Linting code..."
	@$(GO) install golang.org/x/lint/golint@latest
	@golint $(SRC_DIRS)

# Цель для форматирования кода
fmt:
	@echo "Formatting code..."
	@$(GO) fmt $(SRC_DIRS)

# Цель для очистки временных файлов
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@echo "Cleanup complete."