# Build Carrier in a stock Go builder container
FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /carrier-go
RUN cd /carrier-go && make carrier

# Pull Carrier into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /carrier-go/build/bin/carrier /usr/local/bin/

EXPOSE 6060 4000 3500 13000 12000/udp
ENTRYPOINT ["carrier"]
