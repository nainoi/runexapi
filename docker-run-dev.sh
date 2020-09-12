#!/bin/bash
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

docker run -itd --name $NAME -v $(pwd):/app -p 3008:3006 --memory="1g" --cpus="1" registry.thinkdev.app/think/runex/runexapi:dev;
docker logs -f $NAME