# Используем последнюю версию Go 1.22 на базе Alpine
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN go build -o main ./cmd/server/main.go

# Финальный минимальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарник из builder-образа
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]
