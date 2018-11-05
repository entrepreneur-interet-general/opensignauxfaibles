#!/bin/bash
# Creates a base export_to and copies its content to a feature.csv file
# batch and algo are given as parameters e.g:
# bash export.sh 1808 algo2
batch=$1
algo=$2
cd /home/pierre/Documents/opensignauxfaibles/dbmongo/export/
mongo --eval "var batch = \"$batch\", algo = \"$algo\"" export_aggregate.js
mongoexport --db opensignauxfaibles --collection to_export --out ../../output/features/features.csv --type=csv --fieldFile export_fields.txt
