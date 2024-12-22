# Makefile

# Переменные
PKG = ./...
LINTER = golangci-lint
LINTER_CONFIG = .golangci.yml

# Стадии
.PHONY: all lint

all: lint

# Стадия lint
lint:
	@echo "Установка линтера..."
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo "Запуск линтера..."
	$(LINTER) run $(PKG) --config=$(LINTER_CONFIG)

# Стадия для установки линтера (если необходимо)
install-linter:
	@echo "Установка golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.51.2

# Стадия для очистки (если нужно)
clean:
	@echo "Очистка..."
	# Здесь можно добавить команды для очистки, если это необходимо