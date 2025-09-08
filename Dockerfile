# Используем официальный образ Golang
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod и go sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированное приложение
COPY --from=builder /app/main .

# Копируем статические файлы фронтенда
COPY --from=builder /app/frontend ./frontend

# Открываем порт
EXPOSE 3000

# Запускаем приложение
CMD ["./main"]