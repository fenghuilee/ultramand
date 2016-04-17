package client

import (
	"sync"
	"ultraman/lib/log"

	httpclient "ultraman/lib/conn/http"
	websocketclient "ultraman/lib/conn/websocket"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

func StartClient(authKey, webSocketAddr string) {
	wg.Add(1)
	buildWebSocketClient(&wg, authKey, webSocketAddr)
	wg.Wait()
}

func buildWebSocketClient(wg *(sync.WaitGroup), auth, addr string) {

	defer (*wg).Done()

	wsClient := &websocketclient.ClientClient{
		Addr:    addr,
		AuthKey: auth,
	}

	httpClient := &httpclient.ClientClient{}

	if wsClient.Dial() == true {
		if wsClient.Auth() == true {
			for {
				mt, msg, err := wsClient.Conn.ReadMessage()
				if err != nil {
					log.Warn("Failed to read websocket: %v", err)
					return
				}
				log.Debug("Websocket recv: %v,%s", mt, msg)
				// DATA
				if mt == websocket.BinaryMessage {

					resp := httpClient.OpenUrl(&msg)

					log.Debug("%s", resp)

					err = wsClient.Conn.WriteMessage(mt, resp)
					if err != nil {
						log.Warn("Failed to write websocket: %v", err)
						break
					}
				}
			}
		}
	}

}
