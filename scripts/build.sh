IMAGE="registry.xuelangyun.com/shuzhi-amd64/suanpan_go_sdk_development"
VERSION="0"
VERSION_TAIL="_0.0.1"
docker build -t ${IMAGE}:${VERSION} . -f ./Dockerfile
docker tag ${IMAGE}:${VERSION} ${IMAGE}:${VERSION}${VERSION_TAIL}
docker tag ${IMAGE}:${VERSION} ${IMAGE}:latest

docker login registry.cn-shanghai.aliyuncs.com/shuzhi
docker push ${IMAGE}:${VERSION}
docker push ${IMAGE}:${VERSION}${VERSION_TAIL}
docker push ${IMAGE}:latest

