FROM golang:1.22.2 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кешируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}
ENV CGO_ENABLED=0

CMD go test -v -tags integration ./integration_test