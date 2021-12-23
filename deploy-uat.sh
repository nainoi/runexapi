SERVER_IP=ip.thinkdev.app
docker image build -t registry.thinkdev.app/think/runex/runexapi:dev -f Dockerfile.uat .

docker save registry.thinkdev.app/think/runex/runexapi:dev > runexapi.tar \
&& scp runexapi.tar root@$SERVER_IP:/root \
&& scp docker-run-uat.sh root@$SERVER_IP:/root \
&& ssh root@$SERVER_IP "docker load < runexapi.tar && sh docker-run-uat.sh && exit" \
&& rm runexapi.tar