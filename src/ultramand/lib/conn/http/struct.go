package http

import (
	"net"
	"time"
)

const (
	// Time allowed to read/write the tcp connection to the client.
	TimeKeepAlive = 60 * time.Second
	ReadTimeOut   = 30 * time.Second
)

type Server struct {
	Addr           string
	Clients        map[string](*Client)
	onNewClient    func(c *Client)
	onNewRequest   func(c *Client, message []byte)
	onClientClosed func(c *Client)
}

type Client struct {
	Conn     *net.Conn
	Server   *Server
	ReadBuf  []byte
	WriteBuf []byte
}

type ClientClient struct {
	Conn     map[string](*net.Conn) //[domain] *net.Conn
	ReadBuf  []byte
	WriteBuf []byte
}
