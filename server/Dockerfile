# Билд Go-приложения
FROM golang:1.22.6 AS server-builder

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .

# Переменные для сборки
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN make build-app

# ----------------------------------------------------------------------------------------

# Образ с бинарником
FROM alpine:latest

# Установка зависимостей для запуска бинарника
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add nginx

WORKDIR /app

# Копируем серверный бинарник
COPY --from=server-builder /build/bin/sso /app
COPY --from=server-builder /build/migrations /app/migrations

RUN chmod +x /app/sso

CMD ["./sso"]