# Makefile

# Переменные
PKG = ./...
LINTER = golangci-lint
LINTER_CONFIG = ../.golangci.yml

# Стадии
.PHONY: all lint install-linter clean build generate calendar calendar_scheduler calendar_storer

all: lint

# Запуск Docker Compose
up:
	@echo "Run Docker Compose..."
	docker-compose -f ./docker-compose.yml up

# Остановка Docker Compose
down:
	@echo "Stop Docker Compose..."
	docker-compose -f ./docker-compose.yml down


#LINUX
#integration-tests:
#	docker-compose -f ./docker-compose.integra.yml up --build -d
#	test_status_code=0
#	docker-compose -f docker-compose.integra.yml run integration_tests go test ./integration_test/calendar_test.go || test_status_code=$$?
#	docker-compose -f docker-compose.integra.yml down
#	echo $$test_status_code
#	exit $$test_status_code
# -- build

build-integration-tests:
	docker-compose -f .\docker-compose.integra.yml up
	@echo "Waiting for starting services..."
	sleep 30

#Для WINDOWS
integration-tests: build-integration-tests
	set test_status_code=0
	docker-compose -f .\docker-compose.integra.yml run integration_tests go test ./integration_test/calendar_test.go || set test_status_code=$$?
	docker-compose -f .\docker-compose.integra.yml down
	echo %test_status_code%
	exit /b %test_status_code%

# ДЛЯ WINDOWS
generate:
#	@if exist internal/server/grpc_local/pb ( \
#		echo "Удаление старой директории..."; \
#		rmdir /s /q internal/server/grpc_local/pb \
#	)
#	@echo "Создание новой директории..."
#	mkdir internal/server/grpc_local/pb
	protoc \
        --proto_path=internal/server/grpc_local/proto \
        --go_out=internal/server/grpc_local/pb \
        --go-grpc_out=internal/server/grpc_local/pb \
        internal/server/grpc_local/proto/*.proto

# Стадия сборки
build:
	@echo "Сборка проекта..."
	go build ./...
# Стадии сборки для микросервисов
calendar:
	@echo "API services"
	go build -o calendar ./cmd/calendar

calendar_scheduler:
	@echo "Building scheduler"
	go build -o calendar_scheduler ./cmd/calendar_scheduler

calendar_storer:
	@echo "Building storer"
	go build -o calendar_storer ./cmd/calendar_sender

# Запуск сервисов
run_calendar:
	@echo "Запуск API сервиса..."
	./calendar -config=./configs/config.yaml

run_calendar_scheduler:
	@echo "Запуск планировщика..."
	./calendar_scheduler -config=./configs/scheduler_config.yaml

run_calendar_storer:
	@echo "Запуск хранилища..."
	./calendar_storer -config=./configs/sender_config.yaml

# Стадия для очистки (если нужно)
clean:
	@echo "Очистка..."
	# Здесь можно добавить команды для очистки, если это необходимо


# Стадия генерации кода gRPC для LINUX
#generate:
#	rm -rf internal/server/grpc_local/pb
#	mkdir -p internal/server/grpc_local/pb
#
#	protoc \
#        --proto_path=internal/server/grpc_local/proto \
#        --go_out=internal/server/grpc_local/pb \
#        --go-grpc_out=internal/server/grpc_local/pb \
#        internal/server/grpc_local/proto/*.proto

# Стадия lint
#lint: install-linter build
#	@echo "Запуск линтера..."
#	golangci-lint run ./...

# Стадия для установки линтера
#install-linter:
#	@echo "Проверка установки golangci-lint..."
#	@if ! command -v $(LINTER) &> /dev/null; then \
#		echo "golangci-lint не найден. Установка..."; \
#		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.51.2; \
#	else \
#		echo "golangci-lint уже установлен."; \
#	fi

