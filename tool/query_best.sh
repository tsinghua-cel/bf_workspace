#!/bin/bash
source .env.local
set -e

# Verify environment variables (must be set by user)
: "${MYSQL_HOST:?MYSQL_HOST must be set (remote MySQL server IP)}"
: "${MYSQL_PORT:=3306}"
: "${MYSQL_USER:?MYSQL_USER must be set (MySQL username)}"
: "${MYSQL_DATABASE:?MYSQL_DATABASE must be set (target database name)}"
: "${MYSQL_PASSWORD:?MYSQL_PASSWORD must be set (target database password)}"

function queryBest() {
  SQL_COMMANDS=$(cat <<'SQL'
  SELECT uuid, honest_lose_rate_avg, attacker_lose_rate_avg from t_strategy where is_end=1 and honest_lose_rate_avg > 0 and attacker_lose_rate_avg < 1 order by (honest_lose_rate_avg - attacker_lose_rate_avg) desc limit 1;
SQL
  )

  info=$(docker run --rm -i \
    --env MYSQL_PWD="$MYSQL_PASSWORD" \
    mysql:8.0 \
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" "$MYSQL_DATABASE" \
    -e "$SQL_COMMANDS" \
    --skip-column-names)

    UUID=$(echo "$info" | awk '{print $1}')
#    echo "Best strategy UUID: $UUID"

    HONEST_LOSE_RATE_AVG=$(echo "$info" | awk '{print $2}' | awk '{printf "%.4f", $1 }')
    # HONEST_LOSE_RATE_AVG = 1-HONEST_LOSE_RATE_AVG
    HONEST_LOSE_RATE_AVG=$(echo "1 - $HONEST_LOSE_RATE_AVG" | bc)
#    echo "Honest lose rate average: $HONEST_LOSE_RATE_AVG"
    HONEST_LOSE_RATE_ECHO=$(echo "$HONEST_LOSE_RATE_AVG" | awk '{printf "%.2f%%", $1 * 100}')
#    echo "Honest lose rate echo: $HONEST_LOSE_RATE_ECHO"


    ATTACKER_LOSE_RATE_AVG=$(echo "$info" | awk '{print $3}' | awk '{printf "%.4f", $1 }')
    # ATTACKER_LOSE_RATE_AVG = 1-ATTACKER_LOSE_RATE_AVG
    ATTACKER_LOSE_RATE_AVG=$(echo "1 - $ATTACKER_LOSE_RATE_AVG" | bc)
#    echo "Bazantine lose rate average: $ATTACKER_LOSE_RATE_AVG"
    ATTACKER_LOSE_RATE_ECHO=$(echo "$ATTACKER_LOSE_RATE_AVG" | awk '{printf "%.2f%%", $1 * 100}')
#    echo "Bazantine lose rate echo: $ATTACKER_LOSE_RATE_ECHO"


    # success rate
    SQL_COMMANDS_SR=$(cat <<'SQL'
      SELECT CONCAT(TRUNCATE(
        (SELECT COUNT(1) FROM t_strategy WHERE reorg_count > 0) / (SELECT COUNT(1) FROM t_strategy WHERE is_end = 1) * 100,
        2
      ), '%') AS successRate;
SQL
      )

    SR=$(docker run --rm -i \
      --env MYSQL_PWD="$MYSQL_PASSWORD" \
      mysql:8.0 \
      mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" "$MYSQL_DATABASE" \
      -e "$SQL_COMMANDS_SR" \
      --skip-column-names)

#    echo "Success rate: $SR"

    byzantine_advantage=$(echo "$ATTACKER_LOSE_RATE_AVG - $HONEST_LOSE_RATE_AVG" | bc)
    byzantine_advantage=$(echo "$byzantine_advantage" | awk '{printf "%.2f%%", $1*100 }')

#    echo "Calculating byzantine advantage: $byzantine_advantage"

    echo "honest lose rate: ${HONEST_LOSE_RATE_ECHO}, byzantine lose rate: ${ATTACKER_LOSE_RATE_ECHO}, byzantine advantage: ${byzantine_advantage} , success rate: ${SR}"
}

queryBest