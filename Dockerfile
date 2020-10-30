# Build Gmxt in a stock Go builder container
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-mxt
RUN cd /go-mxt && make gmxt

# Pull Gmxt into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-mxt/build/bin/gmxt /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["gmxt"]
