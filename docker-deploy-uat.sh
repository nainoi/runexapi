SERVER_IP=128.199.163.81
docker image build -t registry.thinkdev.app/think/runex/runexapi:dev -f Dockerfile.uat .

docker save registry.thinkdev.app/think/runex/runexapi:dev > runexapi.tar \
&& scp runexapi.tar root@$SERVER_IP:/root \
&& scp docker-run-uat.sh root@$SERVER_IP:/root \
&& ssh root@$SERVER_IP "docker load < runexapi.tar && sh docker-run-uat.sh && exit" \
&& rm runexapi.tar