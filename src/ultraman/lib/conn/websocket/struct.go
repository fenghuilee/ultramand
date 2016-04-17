package websocket

import (
	"github.com/gorilla/websocket"
)

type Server struct {
	Addr           string
	Clients        map[string](*Client)
	onNewClient    func(c *Client)
	onNewRequest   func(c *Client)
	onNewRespone   func(c *Client, message []byte)
	onClientClosed func(c *Client, err error)
}

type Client struct {
	Conn     *websocket.Conn
	Server   *Server
	ReadBuf  []byte
	WriteBuf []byte
}

type ClientClient struct {
	Addr     string
	AuthKey  string
	Conn     *websocket.Conn
	ReadBuf  []byte
	WriteBuf []byte
}
