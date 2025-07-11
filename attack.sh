#!/bin/bash
casetype=${1:-"none"}

basedir=$(pwd)
export BASEDIR="$basedir/"

allcase="none exante sandwich staircase unrealized withholding selfish staircase-ii pyrrhic-victory"

echo "casetype is $casetype"
case $casetype in
        "none")
                ./v5/runtest.sh $casetype 360
                ;;
        "exante")
                ./v5/runtest.sh $casetype 1800
                ;;
        "sandwich")
                ./v5/runtest.sh $casetype 1800
                ;;
        "staircase")
                ./v4/runtest.sh $casetype 1800
                ;;
        "unrealized")
                ./v4/runtest.sh $casetype 1800
                ;;
        "withholding")
                ./v4/runtest.sh $casetype 1800
                ;;
        "selfish")
                ./v4/runtest.sh selfish 1800
                ;;
        "staircase-ii")
                ./v5/runtest.sh staircaseii 1800
                ;;
        "pyrrhic-victory")
                ./v5/runtest.sh sync 1800
                ;;
        *)
                echo "unsupported casetype $casetype, supported cases are: $allcase"
                ;;
esac
