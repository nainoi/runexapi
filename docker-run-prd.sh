#!/bin/bash
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

PRIVATE_IP=$(hostname -I | cut -d' ' -f3)
docker run -itd --name $NAME  -p $PRIVATE_IP:3008:3006 -v /hosting/www/runex-api/uploads:/root/upload think-runex-api;
docker logs -f $NAME