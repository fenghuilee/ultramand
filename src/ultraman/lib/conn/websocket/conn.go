package websocket

import (
	//"io"
	"net/http"
	"net/url"
	"os"
	//"time"
	"html/template"
	"sync"
	"ultraman/lib/log"

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

	http.HandleFunc("/echo", s.handleWebsocket)
	http.HandleFunc("/", home)

	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		log.Error("Failed to listen public websocket address: %v", err)
		os.Exit(1)
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

	var wg sync.WaitGroup
	wg.Add(1)
	client.Auth(&wg)
	wg.Wait()

	go client.Serve()

}

func (c *Client) Auth(wg *(sync.WaitGroup)) {
	defer (*wg).Done()

	log.Debug("Wait client %v login", c.Conn.RemoteAddr().String())

	c.Conn.WriteMessage(websocket.TextMessage, []byte("Please login!"))

	for {

		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Warn("Failed to read websocket: %v", err)
			return
		}
		log.Debug("Websocket recv: %s", msg)
		// CTL
		c.Conn.WriteMessage(websocket.TextMessage, []byte("ok"))

		break

	}
}

// Read client data from channel
func (c *Client) Serve() {

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
		log.Debug("Websocket recv: %v,%s", mt, msg)
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
	url := url.URL{Scheme: "ws", Host: c.Addr, Path: "/echo"}
	log.Info("Connecting to server %s", url.String())

	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Error("Failed to connect server: %v", err)
		os.Exit(1)
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
		log.Debug("Websocket recv: %s", msg)
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

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}
