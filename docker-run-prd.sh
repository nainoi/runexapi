#!/bin/bash
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

docker run -itd --name $NAME  -p 3008:3006 -v ~/runex-api/uploads:/root/upload think-runex-api;
docker logs -f $NAME