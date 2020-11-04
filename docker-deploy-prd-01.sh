#!/bin/bash
SERVER_IP=128.199.66.67
docker image build -t think-runex-api -f Dockerfile .

docker save think-runex-api > runexapi.tar \
&& scp runexapi.tar root@$SERVER_IP:/root \
&& scp docker-run-prd.sh root@$SERVER_IP:/root \
&& ssh root@$SERVER_IP "sudo docker load < runexapi.tar && sudo sh docker-run-prd.sh && exit" \
&& rm runexapi.tar


# scp docker-run-on-cloud.sh ols-user@203.150.107.41:/home/ols-user \
# && ssh ols-user@203.150.107.41 "sudo docker load < kp-callsceen-note-dev.tar && sudo sh docker-run-on-cloud.sh && exit" \