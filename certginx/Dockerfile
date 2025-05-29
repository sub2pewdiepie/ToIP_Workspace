FROM nginx:1.23

RUN apt-get update && \
    apt-get install -y certbot python3-certbot-nginx cron gettext-base && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY startup.sh /startup.sh

RUN chmod +x /startup.sh

EXPOSE 80 443

CMD ["/startup.sh"]