/*
MIT License

Copyright (c) 2022 Prince Pereira

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package kafka

import (
	c "UniversalClient/config"
	"fmt"
	"strings"
	"sync"

	"github.com/confluentinc/confluent-kafka-go-dev/kafka"
)

const (
	addrPattern = "%s:%s"
)

func Init(conf *c.Config) {
	switch c.TypeAction(strings.ToLower(string(conf.Action))) {
	case c.ActionProduce:
		produce(conf)
	case c.ActionConsume:
		consume(conf)
	default:
		fmt.Println("Unknown action : ", conf.Action)
	}
}

func produce(conf *c.Config) {
	addr := fmt.Sprintf(addrPattern, conf.IpAddr, conf.Port)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"client.id":         "localhost",
		"acks":              "all"})
	if err != nil {
		fmt.Println("Kafka connection failed. Error : ", err)
		return
	}

	del_chan := make(chan kafka.Event, 10000)
	defer close(del_chan)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &conf.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(conf.Message)},
		del_chan,
	)
	if err != nil {
		fmt.Println("Kafka produce failed. Error : ", err)
		return
	}

	channel_out := <-del_chan
	message_report := channel_out.(*kafka.Message)

	if message_report.TopicPartition.Error != nil {
		fmt.Println(message_report.TopicPartition.Error)
		return
	}

	fmt.Println(c.MessageDelivered)
}

func consume(conf *c.Config) {
	topics := strings.Split(conf.Topic, ",")
	addr := fmt.Sprintf(addrPattern, conf.IpAddr, conf.Port)
	fmt.Println("Subscribed topics : ", topics, "\n")

	defer func() {
		fmt.Println("Subscription completed.")
	}()

	wg := &sync.WaitGroup{}
	wg.Add(len(topics))
	for _, topic := range topics {
		go subscribe(addr, topic, wg)
	}
	wg.Wait()
}

func subscribe(addr string, topic string, wg *sync.WaitGroup) {

	// Channel Subscriber
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		wg.Done()
		c.Close()
	}()

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			printMsg(topic, msg)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func printMsg(topic string, m *kafka.Message) {
	fmt.Println("#========================#")
	fmt.Println("Topic     : ", topic)
	fmt.Println("Partition : ", m.TopicPartition)
	fmt.Println("Header    : ", m.Headers)
	fmt.Println("Body      : ", string(m.Value))
	fmt.Println("#========================#\n")
}
