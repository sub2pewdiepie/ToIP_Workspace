#!/bin/bash
set -e

# Проверка переменных
if [ -z "$DOMAIN" ]; then
    echo "ERROR: DOMAIN not set!"
    exit 1
fi

if [ -z "$EMAIL" ]; then
    echo "ERROR: EMAIL not set!"
    exit 1
fi

# Генерация конфига Nginx
envsubst '$DOMAIN' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Запуск cron для автообновления
service cron start

# Получение сертификата (если нет)
if [ ! -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
    echo "Getting SSL certificate for $DOMAIN..."
    mkdir -p /var/www/certbot
    certbot certonly --nginx -n --agree-tos --email $EMAIL -d $DOMAIN \
        --webroot-path /var/www/certbot \
        --staging  # Убрать для продакшена
fi

# Установка cron-задачи для обновления
echo "0 3 * * * certbot renew --quiet && nginx -s reload" | crontab -

# Запуск Nginx
exec nginx -g "daemon off;"