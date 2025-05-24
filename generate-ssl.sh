#!/bin/bash

mkdir -p nginx/ssl

# Генерация самоподписанного сертификата
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem \
  -subj "/C=RU/ST=State/L=City/O=Organization/OU=Unit/CN=localhost"

chmod 600 nginx/ssl/key.pem
