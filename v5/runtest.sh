#!/bin/bash
casetype=${1:-"1"}
testduration=${2:-"1800"}

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

testnormal() {
        # start mysql
        docker compose -f $casedir/mysql.yml up -d 
	      caseduration=9000

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

testrl() {
        # start mysql
        docker compose -f $casedir/mysql.yml up -d
	      testcasenocollect ext-staircase 3600
	      # collect the results from mysql.
	      resultA=$(bash $basedir/tool/query_best.sh $basedir)

	      docker compose -f $casedir/mysql.yml down
	      sudo mv ${basedir}/database ${basedir}/results/ext-staircase/

	      # start mysql
	      docker compose -f $casedir/mysql.yml up -d
	      testcasenocollect rlstaircase 3600
	      # collect the results from mysql.
	      resultB=$(bash $basedir/tool/query_best.sh $basedir)

	      docker compose -f $casedir/mysql.yml down
	      sudo mv ${basedir}/database ${basedir}/results/rlstaircase/

	      # dump result.
	      echo "Results for ext-staircase: $resultA"
	      echo "Results for rl-staircase: $resultB"
}

testcasenocollect() {
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
  mysqlfile=$casedir/mysql.yml
  project=$docase
  echo "docker compose -p $project -f $file down" > /tmp/_stop.sh
  echo "docker compose -f $mysqlfile down" >> /tmp/_stop.sh
  docker compose -p $project -f $file up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -p $project -f $file down
  sudo mv data $resultdir/data
  echo "$docase test finished"
}


testattack() {
          casename=$1
          caseduration=$2
          # check if database dir exists, if exist, backup it.
          dbdir="${basedir}/database"
          if [ -d $dbdir ]; then
            echo "database dir exist, backup it to $dbdir-$(date +%Y%m%d%H%M%S)"
            mv $dbdir $dbdir-$(date +%Y%m%d%H%M%S)
          fi
          # start mysql
          docker compose -f $casedir/mysql.yml up -d

          testcase $casename $caseduration
          # stop mysql
          docker compose -f $casedir/mysql.yml down
          # mv the database dir to results
          resultdir="${basedir}/results/${docase}"
          sudo mv ${dbdir} ${resultdir}/
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
  mysqlfile=$casedir/mysql.yml
  project=$docase
  echo "docker compose -p $project -f $file down" > /tmp/_stop.sh
  echo "docker compose -f $mysqlfile down" >> /tmp/_stop.sh
  docker compose -p $project -f $file up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -p $project -f $file down
  echo "result collect"

  if [ "$docase" != "none" ]; then
    # fetch reorg log and format output.
    grep "reorg" data/beacon2/d.log | awk -v attacktype="$docase" -F ' ' '{for(i=1;i<=NF;i++){if($i ~ /newSlot=/){split($i,a,"="); newSlot=a[2]} if($i ~ /oldSlot=/){split($i,b,"="); oldSlot=b[2]}} if (newSlot && oldSlot) {print attacktype " attack occurs reorganize blocks in slot " newSlot "-" oldSlot "."}}'
  fi
  sudo mv data $resultdir/data

  echo "test finished and all nodes data in $resultdir"
}

echo "casetype is $casetype"
case $casetype in
        "normal")
                testnormal
                ;;
        "strategy")
                teststrategy
                ;;
        "rl")
                testrl
                ;;
        *)
                testattack $casetype $testduration
                ;;
esac
