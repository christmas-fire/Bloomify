# ================= Этап 1: Сборка и генерация документации =================
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые утилиты для сборки
RUN apk add --no-cache git

WORKDIR /app

# Копируем только файлы зависимостей для кеширования этого слоя
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем утилиту swag ВНУТРИ контейнера-сборщика
# Это добавит бинарник swag в $GOPATH/bin, который есть в PATH
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем весь исходный код
COPY . .

# Генерируем документацию Swag ДО сборки приложения
# Утилита swag теперь доступна, т.к. мы ее установили выше.
# Она создаст директорию /app/docs
RUN swag init -g ./cmd/main.go

# Собираем статически слинкованный бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# ================= Этап 2: Финальный легковесный образ =================
FROM alpine:3.21.3

WORKDIR /app

# Создаем пользователя для безопасного запуска
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем ТОЛЬКО бинарник из стадии сборки
COPY --from=builder /app/server .

# Копируем ТОЛЬКО сгенерированную документацию
# Эта команда создаст директорию /app/docs в финальном образе
COPY --from=builder /app/docs ./docs

# Устанавливаем правильные права на файлы
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

EXPOSE 8080

# Запускаем приложение
CMD ["./server"]