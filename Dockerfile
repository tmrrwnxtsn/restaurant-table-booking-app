FROM golang:alpine AS builder

RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            dos2unix \
            ca-certificates && \
    rm -rf /var/cache/apk/*

# установка утилиты migrate, которая будет использоваться в скрипте scripts/entrypoint.sh, чтобы запустить миграции к БД
ARG MIGRATE_VERSION=4.15.2
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN mkdir -p /usr/local/bin/migrate
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin/migrate

WORKDIR /app/

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./apiserver ./cmd/apiserver

RUN dos2unix ./scripts/entrypoint.sh
RUN chmod +x ./scripts/entrypoint.sh


FROM alpine:latest

RUN apk --no-cache add ca-certificates bash
RUN mkdir -p /var/log/app

WORKDIR /app/

COPY --from=builder /usr/local/bin/migrate /usr/local/bin
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/website ./website/
COPY --from=builder /app/apiserver .
COPY --from=builder /app/scripts/entrypoint.sh .
COPY --from=builder /app/configs/*.yml ./configs/

RUN ls -la

ENTRYPOINT ["./entrypoint.sh"]