#!/usr/bin/env bash

while getopts r:esd:  option; do
  case "$option" in
    r) REGION="$OPTARG";;
    e) EFFECTIF=true;;
    s) SIREN=true;;
  esac
done
shift $(($OPTIND -1))

if [ $# -eq 0 ]; then
  echo "No arguments supplied"
  exit 1
fi

cat "$@" |
  if [ -n "$REGION" ]; then csvgrep --regex "$REGION" --columns 24; else cat; fi |
    if [ -n "$EFFECTIF" ]; then csvgrep --invert-match --regex "(NN|00|01|02|03)" --columns 46; else cat; fi |
      if [ -n "$SIREN" ]; then csvcut --quoting 3 --columns 1; else cat; fi 


