package http

import (
	"io"
	"net"
	"strings"
	"time"
	"ultraman/lib/log"
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
	server.OnClientClosed(func(c *Client, err error) {})

	return server
}

// Listens for new http connections from the public internet
func (s *Server) Listen() {
	log.Info("Listening for public http connections on %v", s.Addr)
	listener, err := net.Listen("tcp", s.Addr)
	defer listener.Close()

	if err != nil {
		panic(log.Error("Failed to listen public http address: %v", err))
	}

	for {
		rawConn, err := listener.Accept()
		if err != nil {
			log.Warn("Failed to accept new http connection: %v", err)
			continue
		}

		rawConn.SetReadDeadline(time.Now().Add(TimeKeepAlive))

		client := &Client{
			Conn:   &rawConn,
			Server: s,
		}

		s.onNewClient(client)

		go client.Serve()
	}
}

// Read client data from channel
func (c *Client) Serve() {
	log.Debug("Now serve for %s", (*(c.Conn)).RemoteAddr().String())

	var err error

	defer func() {
		c.Close(err)
	}()

	n := 0
	buf := make([]byte, 512)

	message := ""

	for {
		n, err = (*(c.Conn)).Read(buf)

		if err == io.EOF {
			message = ""
			continue
		}

		if err != nil {
			log.Debug("Failed to read http request message: %v", err)
			break
		}

		message += string(buf[0:n])

		if n > 0 && n < 512 {
			(*(c.Conn)).SetReadDeadline(time.Now().Add(TimeKeepAlive))
			go c.Server.onNewRequest(c, []byte(message))
			message = ""
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
func (s *Server) OnClientClosed(callback func(c *Client, err error)) {
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

func (c *Client) Close(err error) error {
	c.Server.onClientClosed(c, err)
	return (*(c.Conn)).Close()
}

func (c *ClientClient) OpenUrl(message *([]byte)) []byte {

	conn, err := net.Dial("tcp", ":80")

	defer func() {
		conn.SetReadDeadline(time.Now().Add(TimeKeepAlive))
	}()

	if err != nil {
		log.Warn("Failed to dial local http connection: %v", err)
		return []byte{}
	}

	*message = []byte(strings.Replace(string(*message), "8000", "80", 1))
	conn.Write(*message)

	n := 0
	buf := make([]byte, 512)

	respMessage := ""

	for {
		n, err = conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Debug("Failed to read http request message: %v", err)
			break
		}
		respMessage += string(buf[0:n])
		if n > 0 && n < 512 {
			break
		}
	}

	return []byte(respMessage)
}
