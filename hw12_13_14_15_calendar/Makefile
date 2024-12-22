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
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --disable exhaustivestruct ./...

# Стадия для очистки (если нужно)
clean:
	@echo "Очистка..."
	# Здесь можно добавить команды для очистки, если это необходимо