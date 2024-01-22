# Image for the echo plantd module.
#
# When building from the root of the project:
#
# create: docker build -t org.plantd.module.echo:latest -f module/echo/Dockerfile .
# run:    docker run -it --rm --name echo \
#           -e PLANTD_BROKER_ENDPOINT=tcp://host.docker.internal:9797 \
#           org.plantd.module.echo:latest
# push:   docker push org.plantd.module.echo:latest

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
RUN go build -o echo ./module/echo

# Final stage
FROM alpine:3.19

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        zeromq \
        czmq

COPY --from=build /app/echo /usr/bin/echo

EXPOSE 5001
CMD [ "/usr/bin/echo" ]
