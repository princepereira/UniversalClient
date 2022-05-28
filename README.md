# UniversalClient
```
This repository contains the source code for universal client
Universal client act as a client for all message buses (Eg: Nats,Kafka) and NosQL DBS (Eg: Etcd, Redis.). 
You will find the binary in bin directory.
After generating the client, it will be copied to repo : https://github.com/princepereira/UniversalClientLib
```

__How to generate binary__:
```
$ go build -o bin/client
```
__How to use__:
```

$ wget https://github.com/princepereira/UniversalClientLib/raw/main/client
$ chmod +x client
$ ./client -c <Type of server. Eg: Nats/Kafka/Etcd> -a <produce/consume> -i <Server IP> -p <Server Port> -t <Topics/Subjects> -m <PUT/Produce Message>

Producer Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World'

Consumer Eg: ./client -c Nats -a consume -i 127.0.0.1 -p 4222 -t "test1,test2"

```
__How to use without banner__:
```
Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World' -H
```
__Bringup Nats for testing the client__:
```
Eg: sudo docker run -d --name nats-server --entrypoint /nats-streaming-server -p 4222:4222 -p 8222:8222 nats-streaming
```

