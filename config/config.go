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

package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	AtribHelp    = "-h"
	AtribIpAddr  = "-i"
	AtribPort    = "-p"
	AtribAction  = "-a"
	AtribTopic   = "-t"
	AtribClient  = "-c"
	AtribMessage = "-m"
	AtribHide    = "-H"
)

type TypeAction string

const (
	ActionProduce TypeAction = "produce"
	ActionConsume TypeAction = "consume"
	ActionPut     TypeAction = "put"
	ActionGet     TypeAction = "get"
	ActionDel     TypeAction = "del"
)

type TypeClient string

const (
	ClientNats  TypeClient = "nats"
	ClientKafka TypeClient = "kafka"
	ClientEtcd  TypeClient = "etcd"
	ClientRedis TypeClient = "redis"
)

const (
	MessageDelivered = "\nMessage Delivered\n\n"
	DataStored       = "\nData Stored\n"
)

type Config struct {
	IpAddr  string     `json:"-i"`
	Port    string     `json:"-p"`
	Action  TypeAction `json:"-a"`
	Topic   string     `json:"-t"`
	Client  TypeClient `json:"-c"`
	Message string     `json:"-m"`
	Hide    string     `json:"-H"`
}

func NewConfig(args map[string]string) (config *Config) {
	config = &Config{}
	b, _ := json.Marshal(args)
	json.Unmarshal(b, config)
	return
}

func PrintBanner(conf *Config) {
	if conf.Hide != "" {

	} else {
		switch TypeAction(strings.ToLower(string(conf.Action))) {
		case ActionConsume:
			PrintBanner2()
		}
	}
}

func PrintBanner2() {
	fmt.Println(" ")
	fmt.Println("#===========================================#")
	fmt.Println("#         Title   : Universal Client        #")
	fmt.Println("#         Version : v28.05.2022             #")
	fmt.Println("#         Author  : Prince Pereira          #")
	fmt.Println("#===========================================#")
	fmt.Println(" ")
}
