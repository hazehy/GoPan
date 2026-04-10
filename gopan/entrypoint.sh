#!/bin/sh
set -eu

envsubst '${MYSQL_ROOT_PASSWORD} ${MYSQL_DATABASE} ${GOPAN_REDIS_ADDR}' \
  < /app/etc/gopan-api.docker.yaml \
  > /app/etc/gopan-api.yaml

exec /app/gopan -f /app/etc/gopan-api.yaml