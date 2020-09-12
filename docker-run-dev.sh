#!/bin/bash
NAME=runex-api
docker stop $NAME;
docker rm -f $NAME;

<<<<<<< HEAD
docker run -itd --name $NAME -v $(pwd):/app -p 3008:3006 --memory="1g" --cpus="1" registry.thinkdev.app/think/runex/runexapi:dev;
=======
docker run -itd --name $NAME -v $(pwd):/app -p 3008:3006 registry.thinkdev.app/think/runex/runexapi:dev /root/swag init;
>>>>>>> 400a98dede8c520540a9b79a094dd4a3290208a2
docker logs -f $NAME