FROM golang:1.22.6 AS builder

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .

# Переменные для сборки
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN make build-app

FROM alpine:latest

# Установка зависимостей для запуска бинарника
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/bin/sso /app
COPY --from=builder /build/migrations /app/migrations

RUN chmod +x /app/sso

CMD ["./sso"]