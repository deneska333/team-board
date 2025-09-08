# Используем официальный образ Golang
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Копируем go mod и go sum файлы
COPY go.mod go.sum ./

# Устанавливаем переменные окружения для Go proxy
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Загружаем зависимости с повторными попытками
RUN go mod download || \
    (sleep 5 && go mod download) || \
    (sleep 10 && go mod download)

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