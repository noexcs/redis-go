package client

import (
	"net"
)

type Client struct {
	Conn net.Conn

	Authenticated bool
}

func (c *Client) Write(bytes []byte) error {
	_, err := c.Conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func New(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

func (c *Client) Close() {
	c.Conn.Close()
}
