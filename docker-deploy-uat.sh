#!/bin/bash
docker image build -t registry.thinkdev.app/think/runex/runexapi:dev -f Dockerfile.uat .

# docker push registry.thinkdev.app/think/runex/runexapi:dev;
# scp docker-run-uat.sh root@128.199.163.81:/root \
# && ssh root@128.199.163.81 "sudo sh docker-run-uat.sh && exit" \
docker save registry.thinkdev.app/think/runex/runexapi:dev > runexapi.tar \
&& scp runexapi.tar root@128.199.163.81:/root \
&& scp docker-run-uat.sh root@128.199.163.81:/root \
&& ssh root@128.199.163.81 "sudo docker load < runexapi.tar && sudo sh docker-run-uat.sh && exit" \
&& rm runexapi.tar


# scp docker-run-on-cloud.sh ols-user@203.150.107.41:/home/ols-user \
# && ssh ols-user@203.150.107.41 "sudo docker load < kp-callsceen-note-dev.tar && sudo sh docker-run-on-cloud.sh && exit" \