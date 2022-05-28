# UniversalClient
This repository contains the source code for universal client
Universal client act as a client for all message buses (Eg: Nats,Kafka) and NosQL DBS (Eg: Etcd, Redis.)
__How to use__:

```
Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World'
```

__How to use without banner__:

```
Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World' -H
```

__Bringup Nats for testing the client__:

```
Eg: sudo docker run -d --name nats-server --entrypoint /nats-streaming-server -p 4222:4222 -p 8222:8222 nats-streaming
```
