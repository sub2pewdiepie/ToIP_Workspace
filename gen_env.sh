#!/bin/bash

PASSWORD=$(jq -r '.password' ./config/config.json)
DBNAME=$(jq -r '.dbname' ./config/config.json)
USERNAME=$(jq -r '.user' ./config/config.json)

# echo "POSTGRES_PASSWORD=$PASSWORD" > .env
# echo "DOMAIN=4edu.su" >> .env
# echo "EMAIL=proninpv2304@gmail.com" >> .env
# echo "POSTGRES_USER=$USERNAME" >> .env
# echo "POSTGRES_DB=$DBNAME" >> .env
# echo "LOG_FORMAT=text" >> .env
# echo "LOG_LEVEL=" >> .env