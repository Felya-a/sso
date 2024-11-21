# Шаг 1: Билд React-приложения
FROM node:22 AS client-builder
WORKDIR /build

COPY ./client/package.json ./
RUN npm i --force --loglevel=error

COPY ./client/ .
RUN npm run build

# ----------------------------------------------------------------------------------------

# Шаг 2: Билд Go-приложения
FROM golang:1.22.6 AS server-builder

WORKDIR /build

COPY ./server/go.mod ./server/go.sum .
RUN go mod download

COPY ./server/ .

# Переменные для сборки
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN make build-app

# ----------------------------------------------------------------------------------------

# Шаг 3: Финальный образ с приложением
FROM alpine:latest

# Установка зависимостей для запуска бинарника
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add nginx

WORKDIR /app
RUN mkdir server

# Копируем серверный бинарник
COPY --from=server-builder /build/bin/sso /app/server
COPY --from=server-builder /build/migrations /app/server/migrations

# Копируем статические файлы React
COPY --from=client-builder /build/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=client-builder /build/build/ /usr/share/nginx/html/

RUN chmod +x /app/server/sso

CMD ["./server/sso"]