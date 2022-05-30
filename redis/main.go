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

package redis

import (
	c "UniversalClient/config"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

const (
	addrPattern = "%s:%s"
)

func Init(conf *c.Config) {
	switch c.TypeAction(strings.ToLower(string(conf.Action))) {
	case c.ActionPut:
		put(conf)
	case c.ActionGet:
		get(conf)
	case c.ActionDel:
		del(conf)
	default:
		fmt.Println("Unsupported action : ", conf.Action, " for Redis. Supported actions are : 'Put/Get/Del'")
	}
}

func redisClient(conf *c.Config) (client *redis.Client) {
	addr := fmt.Sprintf(addrPattern, conf.IpAddr, conf.Port)
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return
}

func put(conf *c.Config) {
	client := redisClient(conf)
	err := client.Set(conf.Topic, conf.Message, 0).Err()
	if err != nil {
		fmt.Println("Write to Redis DB Failed. Error : ", err)
		return
	}
	fmt.Printf(c.DataStored)
}

func get(conf *c.Config) {
	client := redisClient(conf)
	val, err := client.Get(conf.Topic).Result()
	if err != nil {
		fmt.Println("Read from Redis DB Failed. Error : ", err)
		return
	}
	fmt.Println("Key :", conf.Topic, " - Value : ", val)
}

func del(conf *c.Config) {
	client := redisClient(conf)
	_, err := client.Del(conf.Topic).Result()
	if err != nil {
		fmt.Println("Delete from Redis DB Failed. Error : ", err)
		return
	}
	fmt.Println("Delete Key :", conf.Topic, " successful.")
}
