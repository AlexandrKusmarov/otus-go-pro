# Собираем в golang
FROM golang:1.22.2 as build

ENV BIN_FILE /opt/calendar/calendar-sender
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кешируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем код из директории src
COPY build/calendar_sender ${CODE_DIR}

# Собираем статический бинарник Go
ARG LDFLAGS
COPY . /go/src/
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/calendar-sender"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/sender_config.yaml
COPY ./configs/sender_config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}