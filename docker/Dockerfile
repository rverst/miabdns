FROM alpine:3

ARG VERSION
ARG TOKEN

LABEL version="$VERSION"
LABEL maintainer="robert@verst.eu"
LABEL source="https://github.com/rverst/miabdns"

ENV GH_API="https://api.github.com"
ENV GH_REPO="$GH_API/repos/rverst/miabdns"
ENV GH_TAGS="$GH_REPO/releases/tags/v$VERSION"
ENV ASSET_NAME="miabdns_linux-amd64.tar.gz"
ENV AUTH="Authorization: token $TOKEN"
ENV USERAGENT="User-Agent: rverst/miabdns ($VERSION)"


RUN apk add --update --no-cache jq curl ca-certificates && \
    mkdir -p /miabdns

RUN curl -o /dev/null -s -H "$AUTH" -H "$USERAGENT" $GH_REPO || exit 99 && \
    resp=$(curl -s -H "$AUTH" -H "$USERAGENT" "$GH_TAGS") && \
    id=$(echo $resp | jq -e '.assets[] | select(.name == "miabdns_linux-amd64.tar.gz") | .id') && \
    curl -L -H "$AUTH" -H "$USERAGENT" -H 'Accept: application/octet-stream' "$GH_REPO/releases/assets/$id" | tar xvz -C /miabdns && \
    apk del jq curl && \
    rm /miabdns/example_config.json && \
    ls -la /miabdns

WORKDIR /miabdns

VOLUME /config
EXPOSE 8080/tcp

ENV MIABDNS_CONFIG=/config/config.json
ENV MIABDNS_PORT=8080
ENV MIABDNS_BIND=0.0.0.0

ENTRYPOINT ["/miabdns/miabdns"]
