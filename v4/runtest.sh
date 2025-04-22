#!/bin/bash
casetype=${1:-"1"}

basedir=$(pwd)
casedir="${basedir}/v4/case"
export BASEDIR="$basedir/"


PYTHON=$(which python3)

updategenesis() {
        docker run -it --rm -v "${basedir}/v4/config:/root/config" --entrypoint /usr/bin/prysmctl tscel/prysmctl:v4.2.1 \
                testnet \
                generate-genesis \
                --fork=capella \
                --num-validators=256 \
                --output-ssz=/root/config/genesis.ssz \
                --chain-config-file=/root/config/config.yml \
                --geth-genesis-json-in=/root/config/genesis.json \
                --geth-genesis-json-out=/root/config/genesis.json
}

testnormal() {
        # start mysql
        docker compose -f $casedir/mysql.yml up -d 
	caseduration=18000

        # loop 10 times to run testcase basic
        for i in $(seq 1 10); do
                testcase basic $caseduration
        done

        # stop mysql
        docker compose -f $casedir/mysql.yml down

}

teststrategy() {
        # start mysql
        docker compose -f $casedir/mysql.yml up -d 
	caseduration=18000

        # loop 10 times to run ext testcase
        for i in $(seq 1 20); do
                switch=$(($i % 5))
                if [ $switch -eq 0 ]; then
                        testcase ext-exante $caseduration
                elif [ $switch -eq 1 ]; then
                        testcase ext-sandwich $caseduration
                elif [ $switch -eq 2 ]; then
                        testcase ext-staircase $caseduration
                elif [ $switch -eq 3 ]; then
                        testcase ext-unrealized $caseduration
                else
                        testcase ext-withholding $caseduration
                fi
        done
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
        "normal")
                testnormal
                ;;
        "strategy")
                teststrategy
                ;;
        *)
                echo "unsupported casetype $casetype"
                ;;
esac
