FROM golang:1.24-bullseye AS builder

LABEL authors="k3env"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем статически с отключением отладки
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ---

FROM debian:bullseye-slim

# Устанавливаем зависимости для headless Chromium
RUN apt-get update && apt-get install -y --no-install-recommends \
    libasound2 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdbus-1-3 libdrm2 \
    libxkbcommon0 libnss3 libxcomposite1 libxdamage1 libxrandr2 libgbm1 \
    libgtk-3-0 libxss1 libxcb1 libx11-xcb1 libxext6 libxfixes3 fonts-liberation \
    wget ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/app .

# Указываем порт, на котором работает Fiber-сервер
EXPOSE 8080

CMD ["./app"]