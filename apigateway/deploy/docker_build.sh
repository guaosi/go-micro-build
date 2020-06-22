#!/bin/bash
DOCKER_IMAGE_HOST="registry.cn-shenzhen.aliyuncs.com"

DOCKER_IMAGE_NAMESPACE="go_micro"

DOCKER_IMAGE_HUB="apigw"

IMAGE_TAG="v1.0"

WORK_PATH=$(dirname $0)

# 当前位置跳到脚本位置
cd ./${WORK_PATH}

# 取到脚本目录
WORK_PATH=$(pwd)

mkdir bin

# 跨平台 Mac编译Linux 需要交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${WORK_PATH}/bin/ ${WORK_PATH}/../

echo -e "\033[32m编译完成: \033[0m ${WORK_PATH}/bin/"

# 容器制作
docker build -t ${DOCKER_IMAGE_HOST}/${DOCKER_IMAGE_NAMESPACE}/${DOCKER_IMAGE_HUB}:${IMAGE_TAG} -f ./Dockerfile .

echo -e "\033[32m镜像打包完成，请推送: \033[0m ${DOCKER_IMAGE_HOST}/${DOCKER_IMAGE_NAMESPACE}/${DOCKER_IMAGE_HUB}:${IMAGE_TAG}\n"

# 删除原二进制文件以及所在目录
rm -rf bin

echo -e "\033[32m残留二进制文件清理成功"