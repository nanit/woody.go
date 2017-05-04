package main

import "time"
import "./statsd"

func inc(c *statsd.Client) {
	for {
		c.Gauge("go.metric", 111)
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	cfg := &statsd.Config{Address: "localhost:8125", Prefix: "test", SocketTTL: 10}
	c := statsd.NewClient(cfg)

	go inc(c)
	go inc(c)
	go inc(c)
	time.Sleep(time.Second * 500)
}
