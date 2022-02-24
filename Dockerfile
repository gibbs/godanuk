FROM golang:alpine3.15 AS build
LABEL maintainer "Dan Gibbs <dev@dangibbs.co.uk>"

WORKDIR /go/src/github.com/gibbs/godanuk
ENV GODANUK_VERSION 0.5

RUN apk add --update -t build-deps curl libc-dev gcc libgcc; \
  curl -L --silent -o godanuk.tar.gz https://github.com/gibbs/godanuk/archive/${GODANUK_VERSION}.tar.gz && \
  tar -xzf godanuk.tar.gz --strip 1 &&  \
  go get -d && \
  go build -o /usr/local/bin/godanuk && \
  apk del --purge build-deps && \
  rm -rf /var/cache/apk/* && \
  rm -rf /go

FROM alpine:3.15
RUN apk add --update bash uuidgen

COPY --from=build /usr/local/bin/godanuk /usr/local/bin/godanuk
EXPOSE 8084
ENTRYPOINT ["/usr/local/bin/godanuk"]
