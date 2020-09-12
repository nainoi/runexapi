#!/bin/bash
GIT_COMMIT=$(git log -1 --format=%h)

docker push registry.thinkdev.app/think/runex/runexapi:$GIT_COMMIT;
scp docker-run-uat.sh root@128.199.163.81:/root \
&& ssh root@128.199.163.81 "sudo sh docker-run-uat.sh && exit" \


