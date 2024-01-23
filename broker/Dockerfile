# Image for the plantd broker service.
#
# When building from the root of the project:
#
# create: docker build -t geoffjay/plantd-broker:latest -f broker/Dockerfile .
# run:    docker run -it --rm --name broker \
#           -e PLANTD_BROKER_ENDPOINT=tcp://*:9797 \
#           geoffjay/plantd-broker:latest
# push:   docker push geoffjay/plantd-broker:latest

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
RUN make build-broker

# Final stage
FROM alpine:3.19

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        zeromq \
        czmq

COPY --from=build /app/build/plantd-broker /usr/bin/broker

EXPOSE 9797
CMD [ "/usr/bin/broker" ]
