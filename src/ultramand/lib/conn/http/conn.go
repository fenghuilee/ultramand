package http

import (
	"io"
	"net"
	"time"
	"ultramand/lib/log"
)

// Creates new server instance
func New(addr string) *Server {
	log.Info("Creating http server with address %v", addr)
	server := &Server{
		Addr:    addr,
		Clients: make(map[string](*Client)),
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewRequest(func(c *Client, message []byte) {})
	server.OnClientClosed(func(c *Client) {})

	return server
}

// Listens for new http connections from the public internet
func (s *Server) Listen() {
	l, err := net.Listen("tcp", s.Addr)

	defer func() {
		l.Close()
	}()

	if err != nil {
		panic(log.Error("Failed to listen public http address: %v", err))
	}
	log.Info("Listening for public http connections on %v", s.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Debug("Failed to accept new http connection: %v", err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(ReadTimeOut))

		client := &Client{
			Conn:   &conn,
			Server: s,
		}

		s.onNewClient(client)

		go client.Serve()
	}
}

// Read client data from channel
func (c *Client) Serve() {

	defer func() {
		c.Close()
	}()

	var err error
	n := 0
	buf := make([]byte, 4)
	message := ""
	for {
		n, err = (*(c.Conn)).Read(buf)
		if err != nil {
			break
		}
		if n != 0 {
			message += string(buf[0:n])
		}
		if n > 0 && n < 4 {
			(*(c.Conn)).SetReadDeadline(time.Now().Add(ReadTimeOut))
			go c.Server.onNewRequest(c, []byte(message))
			message = ""
			continue
		}

	}

}

// Called right after server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClient = callback
}

// Called when Client receives new message
func (s *Server) OnNewRequest(callback func(c *Client, message []byte)) {
	s.onNewRequest = callback
}

// Called right after connection closed
func (s *Server) OnClientClosed(callback func(c *Client)) {
	s.onClientClosed = callback
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := (*(c.Conn)).Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := (*(c.Conn)).Write(b)
	return err
}

func (c *Client) Close() error {
	c.Server.onClientClosed(c)
	return (*(c.Conn)).Close()
}

func (c *ClientClient) OpenUrl(message *([]byte), addr *string) []byte {
	conn, err := net.Dial("tcp", *addr)
	defer func() {
		conn.SetReadDeadline(time.Now().Add(ReadTimeOut))
	}()

	if err != nil {
		log.Warn("Failed to dial local http connection: %v", err)
		return []byte{}
	}

	conn.Write(*message)

	n := 0
	buf := make([]byte, 1024)

	respMessage := ""

	for {
		n, err = conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warn("Failed to read http request message: %v", err)
			break
		}
		respMessage += string(buf[0:n])
		if n > 0 && n < 1024 {
			break
		}
	}

	return []byte(respMessage)
}
