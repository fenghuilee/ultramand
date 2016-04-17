package http

import (
	"net"
	"time"
)

const (
	NotAuthorized = `HTTP/1.0 401 Not Authorized
WWW-Authenticate: Basic realm="ultraman"
Content-Length: 23

Authorization required
`

	NotFound = `HTTP/1.0 404 Not Found
Content-Length: %d

Tunnel %s not found
`

	BadRequest = `HTTP/1.0 400 Bad Request
Content-Length: 12

Bad Request
`
)

const (
	// Time allowed to read/write the tcp connection to the client.
	TimeKeepAlive = 60 * time.Second
)

type Server struct {
	Addr           string
	Clients        map[string](*Client)
	onNewClient    func(c *Client)
	onNewRequest   func(c *Client, message []byte)
	onClientClosed func(c *Client, err error)
}

type Client struct {
	Conn     net.Conn
	Server   *Server
	ReadBuf  []byte
	WriteBuf []byte
}

type ClientClient struct {
	Conn     map[string](*net.Conn) //[domain] *net.Conn
	ReadBuf  []byte
	WriteBuf []byte
}
