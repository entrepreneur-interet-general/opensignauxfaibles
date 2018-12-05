#!/usr/bin/env bash

#for f in *\ *; do mv "$f" "${f// /_}"; done
merge_admin_urssaf()
{
  FILES=( *.csv )
  awk 'NR==1'  "${FILES[0]}" | iconv -f ISO-8859-1 -t UTF-8 | dos2unix > admin_urssaf.csv
  echo *.csv | xargs -n 1 awk 'NR!=1' | dos2unix > admin_urssaf_temp
  awk '!a[$0]++' admin_urssaf_temp >> admin_urssaf.csv
  rm admin_urssaf_temp
  echo ../*/admin_urssaf/ | xargs -n 1 cp admin_urssaf.csv  
}
#sort -u -r */admin_urssaf/*.csv > admin_urssaf.csv


