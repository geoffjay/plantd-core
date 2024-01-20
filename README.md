[![codecov](https://codecov.io/gh/geoffjay/plantd/graph/badge.svg?token=sHiAEpWC7e)](https://codecov.io/gh/geoffjay/plantd)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

---

# Plantd

Core services for building distributed control systems.

Pre-alpha state, might burn your house down, don't use.

## Purpose

So much of `plantd` related tooling has been spread apart, the purpose of this
project is to attempt to bring together the `go` services that are actually
used in the hopes that all the rest can be archived one day.

## Quickstart

```shell
make
sudo make install
# eg. to test one of the services
PLANTD_PROXY_LOG_LEVEL=debug PORT=8080 plantd-proxy
```

For additional service control steps see the contents of `init/` for
`launchctl` and `systemd` options.

## Projects

### 🏚 Existing

The list of projects that should be brought into this one:

* [Broker][broker]
* [Command Line Client][plantcli]
* [Control Tool][plantctl]

### 🏠 Planned

* [Core](core/README.md)
* [Identity](identity/README.md)
* [Proxy](proxy/README.md)
* [State](state/README.md)

## Contributing

It's recommended that some common tooling and commit hooks be installed.

```shell
make setup
```

Once complete you can start everything with `docker` and `overmind`.

```shell
docker compose up -d
overmind start
```

<!-- links -->

[broker]: https://gitlab.com/plantd/broker
[plantctl]: https://gitlab.com/plantd/plantctl
[plantcli]: https://gitlab.com/plantd/plantcli
