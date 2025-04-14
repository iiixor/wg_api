#!/bin/bash

# Генерация SSL-сертификатов, если нужно
if [ ! -f nginx/ssl/cert.pem ]; then
    echo "Generating SSL certificates..."
    bash ./generate-ssl.sh
fi

# Инициализация WireGuard
echo "Initializing WireGuard..."
bash ./init-wireguard.sh

# Запуск с Docker Compose
echo "Starting services..."
docker-compose up -d

echo "Services started! Access Swagger UI at https://localhost/swagger/index.html"
