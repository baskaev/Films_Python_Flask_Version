# Используем официальный образ Go
FROM golang:1.23.4-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum, чтобы установить зависимости
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Копируем исходный код проекта
COPY . .

# Собираем приложение
RUN go build -o main .

# Указываем порт, который будет слушать приложение
EXPOSE 8081

# Запускаем приложение
CMD ["./main"]
