# Используем официальный образ Golang для сборки и запуска
FROM golang:1.21-alpine


WORKDIR /app


COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 8082

CMD ["./main"]