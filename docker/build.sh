#!/usr/bin/env bash

docker build --label "version=$1" --tag "rverst/miabdns:$1" --tag "rverst/miabdns:latest" .
docker push --tag "rverst/miabdns:$1" --tag "rverst/miabdns:latest"
