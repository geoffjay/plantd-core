# Image for the plantd logger service.
#
# When building from the root of the project:
#
# create: docker build -t geoffjay/plantd-logger:latest -f logger/Dockerfile .
# run:    docker run -it --rm --name logger \
#           -e PLANTD_BROKER_ENDPOINT=tcp://host.docker.internal:9797 \
#           geoffjay/plantd-logger:latest
# push:   docker push geoffjay/plantd-logger:latest

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
RUN make build-logger

# Final stage
FROM alpine:3.19

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        zeromq \
        czmq

COPY --from=build /app/build/plantd-logger /usr/bin/logger

EXPOSE 9797
CMD [ "/usr/bin/logger" ]
