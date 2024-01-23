# Image for the plantd identity service.
#
# When building from the root of the project:
#
# create: docker build -t geoffjay/plantd-identity:latest -f identity/Dockerfile .
# run:    docker run -it --rm --name identity \
#           -e PLANTD_BROKER_ENDPOINT=tcp://host.docker.internal:9797 \
#           geoffjay/plantd-identity:latest
# push:   docker push geoffjay/plantd-identity:latest

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
RUN make build-identity

# Final stage
FROM alpine:3.19

RUN apk update \
    && apk upgrade \
    && apk add --virtual \
        zeromq \
        czmq

COPY --from=build /app/build/plantd-identity /usr/bin/identity

EXPOSE 9797
CMD [ "/usr/bin/identity" ]
