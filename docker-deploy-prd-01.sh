#!/bin/bash
SERVER_IP=128.199.66.67
docker image build -t think-runex-api -f Dockerfile .

docker save think-runex-api > runexapi.tar \
&& scp runexapi.tar root@$SERVER_IP:/root \
&& scp docker-run-prd.sh root@$SERVER_IP:/root \
&& ssh root@$SERVER_IP "docker load < runexapi.tar && sh docker-run-prd.sh && exit" \
&& rm runexapi.tar
