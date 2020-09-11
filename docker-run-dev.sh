#!/bin/bash
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

docker run -itd --name $NAME -v $(pwd):/app -p 3008:3006 registry.thinkdev.app/think/runex/runexapi:dev /root/swag init;
docker logs -f $NAME