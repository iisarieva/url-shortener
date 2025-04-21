# Этап сборки
FROM golang:1.24.1 AS builder

# Рабочая директория
WORKDIR /app

# Устанавливаем swag для генерации Swagger-доков
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем go.mod и go.sum и качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Генерируем Swagger-документацию
RUN swag init -g cmd/server/main.go

# Собираем бинарник
RUN go build -o server ./cmd/server

# Финальный минимальный образ
FROM debian:bookworm-slim

# Устанавливаем сертификаты
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Копируем бинарник
COPY --from=builder /app/server /server

# Копируем Swagger-доки (если ты хочешь отдавать их как часть приложения)
COPY --from=builder /app/docs /app/docs

# Открываем порт
EXPOSE 8080

# Запускаем сервер
ENTRYPOINT ["/server"]
