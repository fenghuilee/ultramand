package websocket

import (
	"net/http"
	"net/url"
	"ultramand/lib/log"

	"github.com/gorilla/websocket"
)

// Creates new server instance
func New(addr string) *Server {

	log.Info("Creating websocket server with address %v", addr)

	server := &Server{
		Addr:    addr,
		Clients: make(map[string](*Client)),
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewRequest(func(c *Client) {})
	server.OnNewRespone(func(c *Client, message []byte) {})
	server.OnClientClosed(func(c *Client, err error) {})

	return server
}

// Listens for new websocket connections from the public internet
func (s *Server) Listen() {

	log.Info("Listening for public websocket connections on %v", s.Addr)

	http.HandleFunc("/", s.handleWebsocket)

	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		panic(log.Error("Failed to listen public websocket address: %v", err))
	}
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{} // use default options
	rawConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to Upgrade: %v", err)
		panic(err)
	}

	client := &Client{
		Conn:   rawConn,
		Server: s,
	}

	s.onNewClient(client)

	go client.Serve()

}

// Read client data from channel
func (c *Client) Serve() {

	log.Debug("Now serve for %s", c.Conn.RemoteAddr().String())

	defer func() {
		c.Conn.Close()
		c.Server.onClientClosed(c, nil)
	}()
	//c.Server.onNewRequest(c)

	for {
		mt, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Warn("Failed to read websocket: %v", err)
			return
		}
		// DATA
		if mt == websocket.BinaryMessage {

			c.Server.onNewRespone(c, msg)
			if err != nil {
				log.Warn("Failed to write websocket: %v", err)
				break
			}
		}
	}

}

// Listens for new websocket connections from the public internet
func (c *ClientClient) Dial() bool {
	url := url.URL{Scheme: "ws", Host: c.Addr, Path: "/"}
	log.Info("Connecting to server %s", url.String())

	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		panic(log.Error("Failed to connect server: %v", err))
	}

	c.Conn = ws

	return true
}

func (c *ClientClient) Auth() bool {

	c.Conn.WriteMessage(websocket.TextMessage, []byte(c.AuthKey))

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Warn("Failed to read websocket: %v", err)
			return false
		}
		// CTL
		if string(msg) == "ok" {
			log.Info("Auth success")
			return true
		}
	}

	return false
}

// Called right after server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClient = callback
}

// Called when Client receives new message
func (s *Server) OnNewRequest(callback func(c *Client)) {
	s.onNewRequest = callback
}

func (s *Server) OnNewRespone(callback func(c *Client, message []byte)) {
	s.onNewRespone = callback
}

// Called right after connection closed
func (s *Server) OnClientClosed(callback func(c *Client, err error)) {
	s.onClientClosed = callback
}

// Send text message to client
func (c *Client) Send(message string) error {
	return c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// Send bytes to client
func (c *Client) SendBinary(b []byte) error {
	return c.Conn.WriteMessage(websocket.BinaryMessage, b)
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
