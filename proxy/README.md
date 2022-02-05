# 📨 Proxy Service

Underlying services in a `plantd` network use a ZeroMQ API that isn't a lot of
fun to work with. The intention behind this is to make it easier to communicate
with a service using other protocols.

Some ideas of protocols to translate calls for are:

* REST
* GraphQL
* gRPC
