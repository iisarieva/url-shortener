# Используем официальный образ Go
FROM golang:1.24.1 as builder

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o server ./cmd/server

# Финальный минимальный образ
FROM debian:bookworm-slim

# Устанавливаем нужные утилиты (например, curl)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Копируем бинарник из билдера
COPY --from=builder /app/server /server

# Порт, на котором работает приложение
EXPOSE 8080

# Команда запуска
ENTRYPOINT ["/server"]
