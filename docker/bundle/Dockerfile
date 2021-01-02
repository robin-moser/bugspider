# +++++++++++++++++++++++++++++++++++++++++
# Dockerfile: robinmoser/bugspider:bundled
# +++++++++++++++++++++++++++++++++++++++++

FROM robinmoser/bugspider:latest as source
FROM alpine

RUN apk add --no-cache \
        beanstalkd \
        netcat-openbsd \
        ca-certificates

WORKDIR /app

COPY docker/bundle/entrypoint /
COPY --from=source /bugspider /usr/local/bin/

ENTRYPOINT ["/entrypoint"]
