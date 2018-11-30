#!/usr/bin/env bash
while getopts d:  option; do
  case "$option" in
    d) MIN_DATE="$OPTARG";;
  esac
done
shift $(($OPTIND -1))

AWK_COMMAND='
BEGIN { # Semi-column separated csv as input and output
FS = ","
OFS = ","
} 
NR==1 || ($17 != "N" && (($4 == "" && $6 != "F") || $4 > min_date ))
'


awk -F "," -v min_date="${MIN_DATE:-'2014-01-01'}" "$AWK_COMMAND" "$@"
