package main

import "fmt"
import "os"
import "time"

func inc(c *Client) {
	for {
		c.Increment("metric")
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	cfg := &Config{address: "localhost:8125", prefix: "prefix"}
	c, err := NewClient(cfg)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	go inc(c)
	time.Sleep(time.Second * 5)
}
