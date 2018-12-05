#!/bin/bash

# 
#https://owncloud.data.gouv.fr/remote.php/webdav/ /home/pierre/owncloud davfs user,rw,auto 0 0

mount -t davfs https://owncloud.data.gouv.fr/remote.php/webdav/ /home/pierre/owncloud
BATCH="$1"
cd /home/pierre/owncloud/SignauxFaibles/data/urssaf/
find  /home/pierre/owncloud/SignauxFaibles/data/urssaf/   -ctime -20 -name '*.csv' -print | \
 cut -sd "/" -f 8- |\
 awk '$0 = "./"$0' |\
 xargs -d '\n' cp -p --parents -t "/home/pierre/Documents/opensignauxfaibles/data-raw/$BATCH/"

