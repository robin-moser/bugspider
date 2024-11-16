# +++++++++++++++++++++++++++++++++++++++++
# Dockerfile: robinmoser/bugspider:latest
# +++++++++++++++++++++++++++++++++++++++++

FROM golang:1.20-alpine as build

# downlaod latest ca-certificates for in-app ssl validation
RUN apk --no-cache add ca-certificates

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN CGO_ENABLED=0 go build -o bugspider .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/bugspider /bugspider

ENTRYPOINT ["/bugspider"]
