package main

import "time"
import "github.com/nanit/woody.go/statsd"

func inc(c *statsd.Client) {
	for {
		c.Increment("metric")
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	cfg := &statsd.Config{Address: "localhost:8125", Prefix: "prefix", SocketTTL: 1}
	c := statsd.NewClient(cfg)

	go inc(c)
	go inc(c)
	go inc(c)
	time.Sleep(time.Second * 5)
}
