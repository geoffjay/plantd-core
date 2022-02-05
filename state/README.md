# 🗄 Distributed State Management

Currently, modules that exist in a `plantd` network manage their own state and
there's not a good way of persisting data if the service goes down. The idea
behind this service would be to receive state updates using a PUB/SUB system,
and allow for some kind of PUSH/PULL by the modules to load and store state.
