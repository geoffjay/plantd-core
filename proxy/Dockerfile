# Image for the plantd proxy service.
#
# When building from the root of the project:
#
# create: docker build -t geoffjay/plantd-proxy:latest -f proxy/Dockerfile .
# run:    docker run -it --rm --name proxy \
#           -e PLANTD_BROKER_ENDPOINT=tcp://host.docker.internal:9797 \
#           geoffjay/plantd-proxy:latest
# push:   docker push geoffjay/plantd-proxy:latest

# Build stage
FROM golang:1.21.6-alpine3.19 AS build

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        build-dependencies \
        build-base \
        gcc \
        clang \
        pkgconf \
        zeromq-dev \
        czmq-dev \
    && apk add \
        bash

RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN make build-proxy

# Final stage
FROM alpine:3.19

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        zeromq \
        czmq

COPY --from=build /app/build/plantd-proxy /usr/bin/proxy

EXPOSE 9797
CMD [ "/usr/bin/proxy" ]
