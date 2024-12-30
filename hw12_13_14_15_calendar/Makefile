# Makefile

# Переменные
PKG = ./...
LINTER = golangci-lint
LINTER_CONFIG = ../.golangci.yml

# Стадии
.PHONY: all lint install-linter clean build

all: lint

# Стадия сборки
build:
	@echo "Сборка проекта..."
	go build ./...

# Стадия lint
lint: install-linter build
	@echo "Запуск линтера..."
	golangci-lint run --disable exhaustivestruct ./...

# Стадия для установки линтера
install-linter:
	@echo "Проверка установки golangci-lint..."
	@if ! command -v $(LINTER) &> /dev/null; then \
		echo "golangci-lint не найден. Установка..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.51.2; \
	else \
		echo "golangci-lint уже установлен."; \
	fi

# Стадия для очистки (если нужно)
clean:
	@echo "Очистка..."
	# Здесь можно добавить команды для очистки, если это необходимо