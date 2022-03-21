[![Go Report Card](https://goreportcard.com/badge/github.com/geoffjay/plantd/state)](https://goreportcard.com/report/github.com/geoffjay/plantd/state)

---

# 🗄 Distributed State Management

Currently, modules that exist in a `plantd` network manage their own state and
there's not a good way of persisting data if the service goes down. The idea
behind this service would be to receive state updates using a PUB/SUB system,
and allow for some kind of PUSH/PULL by the modules to load and store state.

## Quickstart

* build and run a `broker`

```shell
git clone git@gitlab.com:plantd/broker.git
cd broker
make
PLANTD_BROKER_CONFIG=configs/broker.yaml ./target/plantd-broker
```

* build and run `state`

```shell
make build-state
PLANTD_STATE_LOG_LEVEL=debug ./build/plantd-state
```

* build `client`

```shell
make build-client
```

* run test commands

```shell
./build/plant state set --service="org.plantd.Client" foo bar
./build/plant state get --service="org.plantd.Client" foo
```

## Example

If `libplantd` is installed the following example can be used to demonstrate
all of the currently available calls.

```python
import gi
import json

gi.require_version("Pd", "1.0")

from gi.repository import Pd


client = Pd.Client.new("tcp://127.0.0.1:7200")

service = "org.plantd.Jupyter"
messages = [
    ("create-scope", json.dumps({"service": service})),
    ("set", json.dumps({"service": service, "key": "foo", "value": "oof"})),
    ("set", json.dumps({"service": service, "key": "bar", "value": "rab"})),
    ("get", json.dumps({"service": service, "key": "foo"})),
    ("get", json.dumps({"service": service, "key": "bar"})),
    ("delete", json.dumps({"service": service, "key": "foo"})),
    ("delete", json.dumps({"service": service, "key": "bar"})),
    ("delete-scope", json.dumps({"service": service})),
]

for message in messages:
    client.send_request("org.plantd.State", message[0], message[1])
    response = client.recv_response()
    print(f"{message[0]}: {response}")
```

The data sink can be tested with a source from `libplantd`.

```python
import time

# eg. if,
# frontend (XSUB) on @tcp://127.0.0.1:11000
# backend  (XPUB) on @tcp://127.0.0.1:11001 - plantd-state should be connected here

source = Pd.Source.new(">tcp://localhost:11000", "")
source.start()
for i in range(10):
    source.queue_message(f"{\"key\":\"foo\",\"value\":\"{i}\"}")
    time.sleep(.5)
source.stop()
```
