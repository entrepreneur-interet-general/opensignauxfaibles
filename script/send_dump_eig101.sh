#!/usr/bin/env bash

VERSION=$1

scp "/home/pierre/Documents/opensignauxfaibles/output/mongodump/v${VERSION}/test_signauxfaibles/Features.bson.gz" pierre@eig101:/var/lib/mongodb
