FROM golang:bullseye AS build
LABEL maintainer "Dan Gibbs <dev@dangibbs.co.uk>"

WORKDIR /go/src/github.com/gibbs/godanuk
ENV GODANUK_VERSION 0.5

RUN apt-get update && apt-get install -y curl; \
    curl -L --silent -o godanuk.tar.gz https://github.com/gibbs/godanuk/archive/${GODANUK_VERSION}.tar.gz && \
    tar -xzf godanuk.tar.gz --strip 1 &&  \
    go get -d && \
    go build -o /usr/local/bin/godanuk && \
    rm -rf /go

FROM debian:bullseye
RUN apt-get update && apt-get install -y bash uuid-runtime scrypt whois

COPY --from=build /usr/local/bin/godanuk /usr/local/bin/godanuk
EXPOSE 8084
ENTRYPOINT ["/usr/local/bin/godanuk"]
