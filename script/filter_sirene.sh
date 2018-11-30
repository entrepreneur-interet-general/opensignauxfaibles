#!/usr/bin/env bash

while getopts r:esd:  option; do
  case "$option" in
    r) REGION="$OPTARG";;
    e) EFFECTIF=true;;
    s) SIREN=true;;
    d) DIFF_FILE="$OPTARG";;
  esac
done
shift $(($OPTIND -1))

cat "$@" |
	if [ -n "$REGION" ]; then csvgrep --regex "$REGION" --columns 24; else cat; fi |
  if [ -n "$EFFECTIF" ]; then csvgrep --invert-match --regex "(NN|00|01|02|03)" --columns 46; else cat; fi |
  if [ -n "$SIREN" ]; then csvcut --quoting 3 --columns 1; else cat; fi |
  if [ -n "$DIFF_FILE"]; then sort | comm -2 -3 <(cat "$DIFF_FILE" | sort); else cat; fi


