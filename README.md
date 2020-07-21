# runex farmme
docker image build -t suthisakch/runex-api:0.9.0 -f Dockerfile .
# runex uat
docker image build -t suthisakch/runex-api:0.9.6 -f Dockerfile .
# runex production
docker image build -t suthisakch/runex-api:0.9.5 -f Dockerfile .

- omise
- merchant
# run docker
docker run --name runex-api -p 3006:3006 -v ~/runex-api/uploads:/root/upload --network mongo-network -e "TZ=Asia/Bangkok" -d --restart=always docker.io/suthisakch/runex-api:0.9.5

// local
docker run --name runex-api -v $PWD/runex:/root/upload -p 3006:3006 --network mongo-net -d runex-api
docker network connect mongo-network runex-api

//db
docker run -d --network mongo-network --name mongodb -e MONGO_INITDB_ROOT_USERNAME=idever -e MONGO_INITDB_ROOT_PASSWORD=idever@987 -e TZ=Asia/Bangkok -p 27017:27017 -v /data/db:/data/db -v /data/mongo/config:/etc/mongo/mongodb.conf  --restart=always  mongo:4.2

//อย่าลืมเปลี่ยน url ก่อน build