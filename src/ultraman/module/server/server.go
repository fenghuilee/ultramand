package server

import (
	"sync"
	"ultraman/lib/log"

	httpserv "ultraman/lib/conn/http"
	websocketserv "ultraman/lib/conn/websocket"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup
var httpServer *httpserv.Server
var websockServer *websocketserv.Server

func StartServer(httpAddr, webSocketAddr string) {
	wg.Add(1)
	go buildHttpServer(&wg, httpAddr)

	wg.Add(1)
	go buildWebSocketServer(&wg, webSocketAddr)

	wg.Wait()
}

func buildHttpServer(wg *(sync.WaitGroup), addr string) {
	defer (*wg).Done()

	httpServer = httpserv.New(addr)

	httpServer.OnNewClient(func(c *(httpserv.Client)) {
		// new client connected
		// lets send some message
		log.Debug("New http connection connected: %v", c.Conn.RemoteAddr().String())
		httpServer.Clients[c.Conn.RemoteAddr().String()] = c
		log.Debug("Total %d http connection(s) connected", len(httpServer.Clients))
	})

	httpServer.OnNewRequest(func(c *(httpserv.Client), message []byte) {
		// new http request message received
		log.Debug("Received host %s new http request:\n<------Request Message------>\n%v\n<------Request Message------>", c.Conn.RemoteAddr(), string(message))
		ProxyHttpRequest(&message)
	})

	httpServer.OnClientClosed(func(c *(httpserv.Client), err error) {
		// connection with client lost
		log.Debug("Http connection disconnected: %v", c.Conn.RemoteAddr().String())
		delete(httpServer.Clients, c.Conn.RemoteAddr().String())
		log.Debug("Total %d http connection(s) connected", len(httpServer.Clients))
	})

	httpServer.Listen()
}

// Handles a new http connection from the public internet
func ProxyHttpRequest(message *([]byte)) {
	for _, wsc := range websockServer.Clients {
		wsc.Conn.WriteMessage(websocket.BinaryMessage, *message)
	}
}

func buildWebSocketServer(wg *(sync.WaitGroup), addr string) {
	defer (*wg).Done()

	websockServer = websocketserv.New(addr)

	websockServer.OnNewClient(func(c *(websocketserv.Client)) {
		// new client connected
		// lets send some message
		log.Debug("New websocket connection connected: %v", c.Conn.RemoteAddr().String())
		websockServer.Clients[c.Conn.RemoteAddr().String()] = c
		log.Debug("Total %d websocket connection(s) connected", len(websockServer.Clients))
	})

	websockServer.OnNewRequest(func(c *(websocketserv.Client)) {
		// new http request message received
		//HandleHttpRequest(c)
	})

	websockServer.OnNewRespone(func(c *(websocketserv.Client), message []byte) {
		// new http request message received
		HandleHttpRespone(&message)
	})

	websockServer.OnClientClosed(func(c *(websocketserv.Client), err error) {
		// connection with client lost
		log.Debug("Websocket connection disconnected: %v", c.Conn.RemoteAddr().String())
		delete(websockServer.Clients, c.Conn.RemoteAddr().String())
		log.Debug("Total %d http connection(s) connected", len(websockServer.Clients))
	})

	websockServer.Listen()
}

func HandleHttpRespone(message *([]byte)) {
	log.Debug("HandleHttpRespone: %s", *message)
	for _, hc := range httpServer.Clients {
		hc.Conn.Write(*message)
	}
}
