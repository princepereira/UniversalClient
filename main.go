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

package main

import (
	c "UniversalClient/config"
	"UniversalClient/kafka"
	"UniversalClient/nats"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var argKeys = map[string]string{
	c.AtribHelp:    "Eg: ./client -c Nats -a produce -i 127.0.0.1 -p 4222 -t test -m 'Hello World'",
	c.AtribIpAddr:  "Server IP. Eg: -i 127.0.0.1",
	c.AtribPort:    "Server Port Number. Eg: -p 4222",
	c.AtribAction:  "Action. Eg. -a produce/consume",
	c.AtribTopic:   "Subject/Topic. Eg: -t test",
	c.AtribClient:  "Type of client. EG: -c Nats/Kafka/Etcd",
	c.AtribMessage: "Message in case of produce/Database. Eg: -m 'Hello World'",
	c.AtribHide:    "Hide Banner",
}

func dialCheck(ip, port string) error {
	addr := fmt.Sprintf("%s:%s", ip, port)
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("unreachable address : %s", addr)
	}
	return nil
}

func validateVals(atrib, val string, args map[string]string) error {
	switch atrib {
	case c.AtribHelp:
		return fmt.Errorf("no values should follow '-h' attribute")
	case c.AtribIpAddr:
		if port, ok := args[c.AtribPort]; ok {
			return dialCheck(val, port)
		}
	case c.AtribPort:
		if ip, ok := args[c.AtribIpAddr]; ok {
			return dialCheck(ip, val)
		}
	case c.AtribAction:
		switch strings.ToLower(val) {
		case "produce", "consume":
			return nil
		default:
			return fmt.Errorf("unsupported values for action: %s. Supported values are : produce/consume", val)
		}
	case c.AtribClient:
		switch strings.ToLower(val) {
		case "nats", "kafka", "etcd":
			return nil
		default:
			return fmt.Errorf("unsupported values for '-c' : %s. Supported values are : nats/kafka/etcd", val)
		}
	}
	return nil
}

func validateArgs() (map[string]string, error) {

	var args = make(map[string]string)

	for i := 1; i < len(os.Args); i++ {

		if i%2 == 1 {

			// Then argument
			if _, ok := argKeys[os.Args[i]]; !ok {
				return nil, fmt.Errorf("unsupported attribute : %s, supported format : %s", os.Args[i], argKeys[c.AtribHelp])
			}
			if _, ok := args[os.Args[i]]; ok {
				return nil, fmt.Errorf("repeated attribute : %s, supported format : %s", os.Args[i], argKeys[c.AtribHelp])
			}

			if i == 1 && os.Args[i] == c.AtribHelp {
				args[os.Args[i]] = "true"
				return args, nil
			}

			if os.Args[i] == c.AtribHide {
				args[os.Args[i]] = "true"
				i++
				return args, nil
			}

		} else {

			// Value
			if _, ok := argKeys[os.Args[i]]; ok {
				return nil, fmt.Errorf("value : %s should not match attribute, supported format : %s", os.Args[i], argKeys[c.AtribHelp])
			}

			if err := validateVals(os.Args[i-1], os.Args[i], args); err != nil {
				return nil, err
			}

			args[os.Args[i-1]] = os.Args[i]

		}
	}

	return args, nil
}

func dealHelp(args map[string]string) bool {
	if _, ok := args[c.AtribHelp]; ok {
		if len(args) > 1 {
			fmt.Println("unsupported command format. Follow ", argKeys[c.AtribHelp])
		} else {
			fmt.Println(argKeys[c.AtribHelp])
		}
		return true
	}
	return false
}

func finalValidation(args map[string]string) error {
	var p string
	var ok bool

	if _, ok = args[c.AtribIpAddr]; !ok {
		return fmt.Errorf("unsupported command format. Missing ip address. %s", argKeys[c.AtribHelp])
	}
	if _, ok = args[c.AtribPort]; !ok {
		return fmt.Errorf("unsupported command format. Missing port number. %s", argKeys[c.AtribHelp])
	}
	if p, ok = args[c.AtribAction]; !ok {
		return fmt.Errorf("unsupported command format. Missing argument for produce/consume. %s", argKeys[c.AtribHelp])
	}
	if _, ok = args[c.AtribTopic]; !ok {
		return fmt.Errorf("unsupported command format. Missing argument for Subject/Topic. %s", argKeys[c.AtribHelp])
	}
	if _, ok = args[c.AtribClient]; !ok {
		return fmt.Errorf("unsupported command format. Missing argument for type of clients nats/kafka/etcd. %s", argKeys[c.AtribHelp])
	}
	if _, ok = args[c.AtribMessage]; !ok && strings.ToLower(p) == "produce" {
		return fmt.Errorf("unsupported command format. Missing argument for message as action is 'produce'. %s", argKeys[c.AtribHelp])
	}
	return nil
}

func main() {
	args, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	if dealHelp(args) {
		return
	}

	err = finalValidation(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	conf := c.NewConfig(args)

	if conf.Hide != "" {

	} else if conf.Action == c.ActionConsume {
		c.PrintBanner()
	}

	switch c.TypeClient(strings.ToLower(string(conf.Client))) {
	case c.ClientNats:
		nats.Init(conf)
	case c.ClientKafka:
		kafka.Init(conf)
	case c.ClientEtcd:
		fmt.Println("Etcd not yet supported")
	default:
		fmt.Println("Unknown client : ", conf.Client)
	}
}
