package statsd

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Config struct {
	Address   string
	Prefix    string
	SocketTTL int64
}

type Client struct {
	socket        *net.UDPConn
	cfg           *Config
	socket_expiry int64
	channel       chan string
}

func NewClient(cfg *Config) *Client {
	channel := make(chan string)
	if cfg.SocketTTL == 0 {
		cfg.SocketTTL = 600
	}
	client := &Client{cfg: cfg, channel: channel}
	go client.udpPublisher()
	return client
}

func (c *Client) udpPublisher() {
	for {
		metric := <-c.channel
		err := c.ensureSocket()
		if err == nil {
			fmt.Printf("published %s\n", metric)
			fmt.Fprintln(c.socket, metric)
		}
	}
}

func (c *Client) closeSocket() {
	if c.socket != nil {
		c.socket.Close()
	}
}

func (c *Client) ensureSocket() error {
	if c.socket == nil || c.socketExpired() {
		return c.createSocket()
	}
	return nil
}

func (c *Client) socketExpired() bool {
	return time.Now().Unix() > c.socket_expiry
}

func (c *Client) createSocket() error {
	c.closeSocket()

	ra, err := net.ResolveUDPAddr("udp", c.cfg.Address)
	if err != nil {
		fmt.Printf("Error resolving address %s: %v", c.cfg.Address, err)
		return err
	}
	conn, err := net.DialUDP("udp", nil, ra)

	if err != nil {
		fmt.Printf("Error creating socket: ", err)
		return err
	}

	c.socket = conn
	c.socket_expiry = time.Now().Unix() + c.cfg.SocketTTL
	return nil
}

func (c *Client) publish(s string) error {
	select {
	case c.channel <- s:
		return nil
	default:
		return errors.New("Cannot publish to channel")
	}
}

func (c *Client) prefix(metric string) string {
	if len(c.cfg.Prefix) > 0 {
		return fmt.Sprintf("%s.%s", c.cfg.Prefix, metric)
	} else {
		return metric
	}
}

func (c *Client) Increment(metric string, val int) error {
	s := fmt.Sprintf("%s:%v|c", c.prefix(metric), val)
	return c.publish(s)
}

func (c *Client) Inc(metric string) error {
	return c.Increment(metric, 1)
}

func (c *Client) Gauge(metric string, val int) error {
	s := fmt.Sprintf("%s:%v|g", c.prefix(metric), val)
	return c.publish(s)
}

func (c *Client) Timing(metric string, ms int) error {
	s := fmt.Sprintf("%s:%v|ms", c.prefix(metric), ms)
	return c.publish(s)
}

func (c *Client) Close() {
	c.closeSocket()
}
