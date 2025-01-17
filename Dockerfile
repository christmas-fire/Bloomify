# Используем официальный образ Golang для сборки приложения
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Используем минимальный образ Alpine для финального образа
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Создаем директорию для конфигов
RUN mkdir -p /root/configs

# Копируем скомпилированное приложение из промежуточного образа
COPY --from=builder /app/main .

# Копируем конфигурационный файл в контейнер в правильную директорию
COPY configs/config.yaml /root/configs/config.yaml

# Экспонируем порт, который будет использоваться приложением
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]