package main

import "time"
import "github.com/nanit/woody.go/statsd"

func inc(c *statsd.Client) {
	for {
		c.Increment("go.metric")
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	cfg := &statsd.Config{Address: "statsd:8125", Prefix: "test", SocketTTL: 10}
	c := statsd.NewClient(cfg)

	go inc(c)
	go inc(c)
	go inc(c)
	time.Sleep(time.Second * 500)
}
