package main

import "time"

func inc(c *Client) {
	for {
		c.Increment("metric")
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	cfg := &Config{address: "localhost:8125", prefix: "prefix", socket_ttl: 1}
	c := NewClient(cfg)

	go inc(c)
	go inc(c)
	go inc(c)
	time.Sleep(time.Second * 5)
}
