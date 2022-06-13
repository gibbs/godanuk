FROM golang:bullseye AS build
LABEL maintainer "Dan Gibbs <dev@dangibbs.co.uk>"

WORKDIR /go/src/github.com/gibbs/godanuk
ENV GODANUK_VERSION 0.5
COPY go.mod go.sum main.go /go/src/github.com/gibbs/godanuk/
RUN go get -d && \
    go build -o /usr/local/bin/godanuk && \
    rm -rf /go

FROM debian:bullseye
RUN apt-get update && apt-get install -y bash bind9-dnsutils pwgen scrypt uuid-runtime whois

COPY --from=build /usr/local/bin/godanuk /usr/local/bin/godanuk
EXPOSE 8084
ENTRYPOINT ["/usr/local/bin/godanuk"]
