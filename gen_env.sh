#!/bin/bash

PASSWORD=$(jq -r '.password' ./config/config.json)

echo "POSTGRES_PASSWORD=$PASSWORD" > .env
echo "DOMAIN=4edu.su" >> .env
echo "EMAIL=proninpv2304@gmail.com" >> .env
