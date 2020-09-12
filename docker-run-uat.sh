#!/bin/bash
GIT_COMMIT=$(git log -1 --format=%h)
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

docker run -itd --name $NAME  -p 3008:3006 registry.thinkdev.app/think/runex/runexapi:$GIT_COMMIT;
docker logs -f $NAME