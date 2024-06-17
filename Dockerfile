# Используем официальный образ Go для сборки
FROM golang:1.22.2-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем бинарный файл
RUN go build -o rps .

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл и конфигурационный файл из builder
COPY --from=builder /app/rps /app/rps
COPY --from=builder /app/config/items.json /app/config/items.json

# Устанавливаем команду для запуска контейнера
CMD ["./rps"]

# Открываем порт
EXPOSE 3000
