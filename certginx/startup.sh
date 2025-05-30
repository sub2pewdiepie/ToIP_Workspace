#!/bin/bash
set -e

# Проверка переменных
[ -z "$DOMAIN" ] && { echo "ERROR: DOMAIN not set!"; exit 1; }
[ -z "$EMAIL" ] && { echo "ERROR: EMAIL not set!"; exit 1; }

# Создаём базовый конфиг для запуска Nginx (только HTTP)
cat > /etc/nginx/conf.d/temp.conf <<EOF
server {
    listen 80;
    server_name $DOMAIN;
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }
    location / {
        return 444;  # Соединение закрывается без ответа
    }
}
EOF

# Запуск Nginx в фоне для верификации домена
nginx -g "daemon on;"

# Получение сертификата (если отсутствует)
if [ ! -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
    echo "Получение SSL-сертификата для $DOMAIN..."
    mkdir -p /var/www/certbot
    
    # Для продакшена уберите --staging
    certbot certonly --webroot -n --agree-tos --email $EMAIL -d $DOMAIN \
        --webroot-path /var/www/certbot \
        --staging || { echo "Ошибка получения сертификата"; exit 1; }
fi

# Останавливаем временный Nginx
nginx -s quit
sleep 2  # Даём время для завершения

# Генерация основного конфига Nginx
envsubst '$DOMAIN' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Настройка cron для автообновления
echo "0 3 * * * certbot renew --quiet && nginx -s reload" | crontab -

# Запуск Nginx с полной конфигурацией
exec nginx -g "daemon off;"