#!/usr/bin/env bash

DATE="$(date "+%Y%m%d%H%M")"
IMAGE_NAME="container-image-sweeper"

docker build -t ${IMAGE_NAME}:"${DATE}" .
docker tag ${IMAGE_NAME}:"${DATE}" ${IMAGE_NAME}:latest
