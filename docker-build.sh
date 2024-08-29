#!/usr/bin/env bash

DOCKER_REPOSITORY_URL="ccr.ccs.tencentyun.com"
DOCKER_REPOSITORY_USERNAME="100005530286"
VERSION="v1.0.0"

timestamp=$(date "+%Y%m%d.%H%M%S")
tag="${VERSION}-${timestamp}"

git pull origin

echo "登录docker仓库"
docker login ${DOCKER_REPOSITORY_URL} --username=${DOCKER_REPOSITORY_USERNAME} --password=${DOCKER_REPOSITORY_PASSWORD}



docker buildx build --platform=linux/amd64 -t ${DOCKER_REPOSITORY_URL}/cqmike/game:${tag} -f api/Dockerfile  -o type=registry .
