#!/bin/bash
casetype=${1:-"none"}

basedir=$(pwd)
casedir="${basedir}/v5/case"
export BASEDIR="$basedir/"


PYTHON=$(which python3)

updategenesis() {
        docker run -it --rm -v "${basedir}/v5/config:/root/config" --entrypoint /usr/bin/prysmctl tscel/prysmctl:v5.2.0 \
                testnet \
                generate-genesis \
                --fork=deneb \
                --num-validators=256 \
                --genesis-time-delay=15 \
                --output-ssz=/root/config/genesis.ssz \
                --chain-config-file=/root/config/config.yml \
                --geth-genesis-json-in=/root/config/genesis.json \
                --geth-genesis-json-out=/root/config/genesis.json
}

runattack() {
        caseduration=1800 # 30 minutes
        # start mysql
        docker compose -f $casedir/mysql.yml up -d
	      testcase $1 $caseduration
        # stop mysql
        docker compose -f $casedir/mysql.yml down
}

testcase() {
  docase=$1
  caseduration=$2
  targetdir="${casedir}/${docase}"
  resultdir="${basedir}/results/${docase}"

  if [ -d $resultdir ]; then
    # backup the resultdir
    echo "resultdir $resultdir exist, backup it to $resultdir-$(date +%Y%m%d%H%M%S)"
    mv $resultdir $resultdir-$(date +%Y%m%d%H%M%S)
  fi
  mkdir -p $resultdir
  echo "run strategy $docase"
  updategenesis
  file=$casedir/attack-$docase.yml
  mysqlfile=$casedir/attack-$docase.yml
  project=$docase
  echo "docker compose -p $project -f $file down" > /tmp/_stop.sh
  echo "docker compose -f $mysqlfile down" >> /tmp/_stop.sh
  docker compose -p $project -f $file up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -p $project -f $file down
  echo "result collect"
  $basedir/tool/query_local.sh
  sudo mv data $resultdir/data

  echo "test done and test data in $resultdir"
}

echo "casetype is $casetype"
case $casetype in
        "1")
                runattack exante
                ;;
        "2")
                runattack sandwich
                ;;
        "3")
                runattack staircase
                ;;
        "4")
                runattack unrealized
                ;;
        "5")
                runattack withholding
                ;;
        "6")
                runattack ext-exante
                ;;
        "7")
                runattack staircaseii
                ;;
        "8")
                runattack sync
                ;;
        *)
                echo "unsupported casetype $casetype"
                ;;
esac
