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

$ wget https://github.com/princepereira/Binaries/raw/main/UniversalClientLib/client
$ chmod +x client
$ ./client -c <Type of server. Eg: Nats/Kafka/Etcd> -a <produce/consume> -i <Server IP> -p <Server Port> -t <Topics/Subjects> -m <PUT/Produce Message>

Nats Producer Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World'

Nats Consumer Eg: ./client -c Nats -a consume -i 127.0.0.1 -p 4222 -t "test1,test2"

Kafka Producer Eg: ./client -c Kafka -a produce -i 127.0.0.1 -p 9092 -t test -m 'Hello World'

Kafka Consumer Eg: ./client -c Kafka -a consume -i 127.0.0.1 -p 9092 -t "test1,test2"

```
__Bringup Nats for testing the client__:
```
Start Nats Server Docker
$ sudo docker run -d --name nats-server --entrypoint /nats-streaming-server -p 4222:4222 -p 8222:8222 nats-streaming

Publish data to Nats Server Subject
$ ./client -c Nats -a Produce -i 127.0.0.1 -p 4222 -t subject1 -m 'Hello World'

Retrieve data from Nats Server Subject
$ ./client -c Nats -a Consume -i 127.0.0.1 -p 4222 -t subject1 
```


__Bringup Kafka for testing the client__:

```
# File:docker-compose.yml
---
version: '3'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.0.1
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  broker:
    image: confluentinc/cp-kafka:7.0.1
    container_name: broker
    ports:
    # To learn about configuring Kafka for access across networks see
    # https://www.confluent.io/blog/kafka-client-cannot-connect-to-broker-on-aws-on-docker-etc/
      - "9092:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_INTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092,PLAINTEXT_INTERNAL://broker:29092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
```
```
Start Kafka and Zookeeper Server Dockers
$ docker-compose up -d

Publish data to Kafka Server Topic
$ ./client -c Kafka -a Produce -i 127.0.0.1 -p 9092 -t topic1 -m 'Hello World'

Retrieve data from Kafka Server Topic
$ ./client -c Kafka -a Consume -i 127.0.0.1 -p 9092 -t topic1 
```

__Bringup Redis for testing the client__:

```
Start Redis Docker
$ docker run --name redis-server -p 6379:6379 -d redis

PUT data to Redis Server
$ ./client -c Redis -a Put -i 127.0.0.1 -p 6379 -t test -m 'Hello World'

GET data from Redis Server
$ ./client -c Redis -a Get -i 127.0.0.1 -p 6379 -t test

DEL data from Redis Server
$ ./client -c Redis -a Del -i 127.0.0.1 -p 6379 -t test

```
