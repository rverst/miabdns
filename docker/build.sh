#!/usr/bin/env bash

while getopts ":u:p:v:t:" opt; do
  case $opt in
    u) DOCKER_USERNAME="$OPTARG"
    ;;
    p) DOCKER_PASSWORD="$OPTARG"
    ;;
    v) VERSION="$OPTARG"
    ;;
    t) GITHUB_TOKEN="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

docker build --no-cache --build-arg VERSION=${VERSION} --build-arg TOKEN=${GITHUB_TOKEN} -t rverst/miabdns:${VERSION} -t rverst/miabdns:latest .

docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
docker push "rverst/miabdns:$VERSION"
docker push "rverst/miabdns:latest"
