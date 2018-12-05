#!/usr/bin/env bash
while getopts d:  option; do
  case "$option" in
    d) MIN_DATE="$OPTARG";;
  esac
done
shift $(($OPTIND -1))

AWK_COMMAND='
NR == 1 || ($17 != "N" && ($4 > min_date || ($4=="" && $6 != "F")))
'


awk -F "," -v min_date="${MIN_DATE:-'2015-01-01'}" "$AWK_COMMAND" "$@"
