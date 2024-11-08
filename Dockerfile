FROM golang:1.23-alpine AS builder

WORKDIR /app

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for

COPY go* .

COPY ./cmd ./cmd
COPY ./config ./config
COPY ./internal ./internal
COPY ./vendor ./vendor

#собираем файл запуска миграций
RUN go build -mod vendor -o ./bin/migrate ./cmd/migrations/main.go
#собираем сервис
RUN go build -mod vendor -o ./bin/main ./cmd/service/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache curl

COPY --from=builder /usr/local/bin/wait-for /usr/local/bin/wait-for
COPY --from=builder /app/bin /bin

COPY ./scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50051 80

ENTRYPOINT [ "/entrypoint.sh" ]