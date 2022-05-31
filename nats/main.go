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

package nats

import (
	c "UniversalClient/config"
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	addrPattern = "%s:%s"
	clusterId   = "test-cluster"
	appName     = "Nats-Client-Prince"
)

func Init(conf *c.Config) {
	switch c.TypeAction(strings.ToLower(string(conf.Action))) {
	case c.ActionProduce:
		produce(conf)
	case c.ActionConsume:
		consume(conf)
	default:
		fmt.Println("Unsupported action : ", conf.Action, " for Nats. Supported actions are : 'Produce/Consume'")
	}
}

func produce(conf *c.Config) {
	addr := fmt.Sprintf(addrPattern, conf.IpAddr, conf.Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		fmt.Printf("Connecting to Nats topic : %s failed. Error : %v", conf.Topic, err)
		return
	}

	defer nc.Close()
	// Simple Publisher
	err = nc.Publish(conf.Topic, []byte(conf.Message))
	if err != nil {
		fmt.Printf("Producing to Nats topic : %s failed. Error : %v", conf.Topic, err)
		return
	}

	fmt.Printf(c.MessageDelivered)
}

func consume(conf *c.Config) {
	topics := strings.Split(conf.Topic, ",")
	addr := fmt.Sprintf(addrPattern, conf.IpAddr, conf.Port)
	var opts []nats.Option

	opts = append(opts, nats.Name("test"))

	opts = append(opts, nats.UserInfo("", ""))

	nc, err := nats.Connect(addr, opts...)
	if err != nil {
		fmt.Printf("Connecting to Nats topic : %s failed. Error : %v", conf.Topic, err)
		return
	}

	fmt.Println("Subscribed topics : ", topics, "\n")
	sc, err := stan.Connect(clusterId, appName, stan.NatsConn(nc))

	if err != nil || sc == nil {
		fmt.Printf("Stan connection to Nats topic : %s failed. Error : %v", conf.Topic, err)
		return
	}

	defer func() {
		fmt.Println("Subscription completed.")
		sc.Close()
		nc.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(len(topics) * 2)
	for _, topic := range topics {
		go subscribeStan(sc, topic, wg)
		go subscribeNats(nc, topic, wg)
	}
	wg.Wait()
}

func subscribeStan(sc stan.Conn, topic string, wg *sync.WaitGroup) {
	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err := sc.Subscribe(topic, func(m *stan.Msg) {
		printStanMsg(m)
	})

	if err != nil {
		fmt.Printf("Subscribing topic : %s failed. Error : %v\n", topic, err)
		return
	}

	defer func() {
		// Unsubscribe
		sub.Unsubscribe()
		close(ch)
		wg.Done()
	}()

	<-ch
}

func subscribeNats(nc *nats.Conn, topic string, wg *sync.WaitGroup) {
	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe(topic, ch)
	defer func() {
		// Unsubscribe
		sub.Unsubscribe()
		close(ch)
		wg.Done()
	}()

	if err != nil {
		fmt.Printf("Subscribing to Nats topic : %s failed. Error : %v", topic, err)
		return
	}

	for msg := range ch {
		printNatsMsg(msg)
	}

}

func printStanMsg(m *stan.Msg) {
	fmt.Println("#========================#")
	fmt.Println("Topic     : ", m.Subject)
	fmt.Println("Timestamp : ", m.Timestamp)
	fmt.Println("Body      : ", string(m.Data))
	fmt.Println("#========================#\n")
}

func printNatsMsg(m *nats.Msg) {
	fmt.Println("#========================#")
	fmt.Println("Topic     : ", m.Subject)
	fmt.Println("Header    : ", m.Header)
	fmt.Println("Body      : ", string(m.Data))
	fmt.Println("#========================#\n")
}
