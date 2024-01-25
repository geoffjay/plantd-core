[![Go Report Card](https://goreportcard.com/badge/github.com/geoffjay/plantd/broker)](https://goreportcard.com/report/github.com/geoffjay/plantd/broker)

---

# Reliable Message Broker

This message broker is based on version 2 of the Majordomo Protocol (MDP/2),
it allows for clients and workers to connect and disconnect in a reliable way
where messages are kept until delivered.

## Quickstart

* build the `broker` service

```shell
make build-broker
PLANTD_BROKER_LOG_LEVEL=debug ./build/plantd-broker
```
