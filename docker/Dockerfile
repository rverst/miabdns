FROM alpine:3.11
LABEL maintainer="robert@verst.eu"
LABEL source="https://github.com/rverst/miabdns"

RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' appuser

ADD ./miabdns /

RUN chown appuser:appuser /miabdns

USER appuser
VOLUME /config
EXPOSE 8080/tcp

ENV MIABDNS_CONFIG=/config/config.json
ENV MIABDNS_PORT=8080
ENV MIABDNS_BIND=0.0.0.0

CMD ["/miabdns"]
