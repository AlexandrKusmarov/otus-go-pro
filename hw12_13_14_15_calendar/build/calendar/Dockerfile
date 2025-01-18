# Собираем в golang
# docker buildx build -t calendar-app2 -f build/calendar/Dockerfile .      - сборка img севриса
FROM golang:1.22.2 as build

ENV BIN_FILE /opt/calendar/calendar-app2
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кешируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем код из директории src
COPY build/calendar ${CODE_DIR}

# Собираем статический бинарник Go
ARG LDFLAGS
COPY . /go/src/
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/calendar-app2"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

ENV CONFIG_MIGRATION /etc/migrations
COPY ./migrations ${CONFIG_MIGRATION}

CMD ${BIN_FILE} -config ${CONFIG_FILE}