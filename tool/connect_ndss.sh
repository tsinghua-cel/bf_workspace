#!/bin/bash
source .env.ndss
set -e

# Verify environment variables (must be set by user)
: "${MYSQL_HOST:?MYSQL_HOST must be set (remote MySQL server IP)}"
: "${MYSQL_PORT:=3306}"
: "${MYSQL_USER:?MYSQL_USER must be set (MySQL username)}"
: "${MYSQL_DATABASE:?MYSQL_DATABASE must be set (target database name)}"
: "${MYSQL_PASSWORD:?MYSQL_PASSWORD must be set (target database password)}"



# Execute commands using temporary MySQL container
echo -e "\nConnecting ...\n"
docker run --rm -it \
  mysql:8.0 \
  mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" 

