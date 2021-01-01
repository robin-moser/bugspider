# +++++++++++++++++++++++++++++++++++++++++
# Dockerfile: robinmoser/bugspider:bundled
# +++++++++++++++++++++++++++++++++++++++++

FROM golang:1.15-alpine as build

RUN mkdir /build
ADD . /build
WORKDIR /build

RUN CGO_ENABLED=0 go build -o bugspider .

FROM alpine as final

RUN apk add --no-cache \
        beanstalkd \
        netcat-openbsd \
        ca-certificates

WORKDIR /app

COPY --from=build /build/bugspider /usr/local/bin/
COPY scripts/ /

ENTRYPOINT ["/entrypoint"]
