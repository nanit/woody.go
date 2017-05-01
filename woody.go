package main

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	address    string
	prefix     string
	socket_ttl int64
}

type Client struct {
	socket        *net.UDPConn
	cfg           *Config
	socket_expiry int64
	channel       chan string
}

func NewClient(cfg *Config) (*Client, error) {
	channel := make(chan string)
	if cfg.socket_ttl == 0 {
		cfg.socket_ttl = 600
	}
	client := &Client{cfg: cfg, channel: channel}
	go client.sendUDP()
	return client, nil
}

func (c *Client) sendUDP() {
	for {
		metric := <-c.channel
		c.ensureSocket()
		fmt.Fprintln(c.socket, metric)
	}
}

func (c *Client) ensureSocket() {
	if c.socket == nil || c.socketExpired() {
		c.createSocket()
	}
}

func (c *Client) socketExpired() bool {
	return time.Now().Unix() > c.socket_expiry
}

func (c *Client) createSocket() {
	if c.socket != nil {
		c.socket.Close()
	}

	ra, err := net.ResolveUDPAddr("udp", c.cfg.address)

	if err != nil {
		fmt.Printf("Some error %v", err)
	}
	conn, err := net.DialUDP("udp", nil, ra)

	if err != nil {
		fmt.Printf("Some error %v", err)
	}
	c.socket = conn
	c.socket_expiry = time.Now().Unix() + c.cfg.socket_ttl
}

func (c *Client) Increment(metric string) {
	s := fmt.Sprintf("%s.%s", c.cfg.prefix, metric)
	c.channel <- s
}
