#!/bin/bash
casetype=${1:-"1"}
caseduration=${2:-"9000"}

basedir=$(pwd)
casedir="${basedir}/case"
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

testcase() {
  docase=$1
  targetdir="${casedir}/${docase}"
  resultdir="${basedir}/results/${docase}"

  if [ -d $resultdir ]; then
    # backup the resultdir
    echo "resultdir $resultdir exist, backup it to $resultdir-$(date +%Y%m%d%H%M%S)"
    mv $resultdir $resultdir-$(date +%Y%m%d%H%M%S)
  fi
  mkdir -p $resultdir
  echo "Running testcase $docase"
  echo "docker compose -p $project -f $file down" > /tmp/_stop.sh
  updategenesis
  file=$casedir/attack-$docase.yml
  project=$docase
  docker compose -p $project -f $file up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -p $project -f $file down
  echo "result collect"
  sudo mv data $resultdir/data

  echo "test done and result in $resultdir"
}

echo "casetype is $casetype"
case $casetype in
        1)
                testcase basic
                ;;
        2)
                testcase mix
                ;;
        3)
                testcase exante
                ;;
        4)
                testcase sandwich
                ;;
        5)
                testcase staircase
                ;;
        6)
                testcase unrealized
                ;;
        7)
                testcase withholding
                ;;
        8)
                testcase ext-exante
                ;;
        9)
                testcase ext-sandwich
                ;;
        10)
                testcase ext-staircase
                ;;
        11)
                testcase ext-unrealized
                ;;
        12)
                testcase ext-withholding
                ;;
        13)
                testcase sync
                ;;
        *)
                testcase $casetype
                ;;
esac
