#!/bin/bash
source .env.ndss
set -e

# Verify environment variables (must be set by user)
: "${MYSQL_HOST:?MYSQL_HOST must be set (remote MySQL server IP)}"
: "${MYSQL_PORT:=3306}"
: "${MYSQL_USER:?MYSQL_USER must be set (MySQL username)}"
: "${MYSQL_DATABASE:?MYSQL_DATABASE must be set (target database name)}"
: "${MYSQL_PASSWORD:?MYSQL_PASSWORD must be set (target database password)}"


# Define SQL query statements with formatted output
SQL_COMMANDS=$(cat <<'SQL'
SELECT 
    'Metric 1' AS Metric,
    COUNT(1) AS Value
FROM t_strategy 
WHERE 
    (honest_lose_rate_avg > 0 AND honest_lose_rate_avg < 1) 
    AND (attacker_lose_rate_avg > 0 AND attacker_lose_rate_avg < 1)
UNION ALL
SELECT 
    'Metric 2' AS Metric,
    COUNT(1) AS Value
FROM t_strategy 
WHERE 
    honest_lose_rate_avg > 1 
    AND attacker_lose_rate_avg > 1
UNION ALL
SELECT 
    'Metric 3' AS Metric,
    COUNT(1) AS Value
FROM t_strategy 
WHERE 
    (honest_lose_rate_avg > 1) 
    AND (attacker_lose_rate_avg > 0 AND attacker_lose_rate_avg < 1)
UNION ALL
SELECT 
    'Metric 4' AS Metric,
    COUNT(1) AS Value
FROM t_strategy 
WHERE 
    (honest_lose_rate_avg > 0 AND honest_lose_rate_avg < 1) 
    AND (attacker_lose_rate_avg > 1) ;
SQL
)

# Execute commands using temporary MySQL container
echo -e "\nExecuting metrics queries...\n"
docker run --rm -i \
  mysql:8.0 \
  mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" \
  -e "$SQL_COMMANDS" \
  -t  # Enable table formatting

echo -e "\nExecution completed - container has been automatically removed"

