package client

import (
	"strings"
	"sync"
	"ultramand/lib/log"

	httpclient "ultramand/lib/conn/http"
	websocketclient "ultramand/lib/conn/websocket"

	"github.com/gorilla/websocket"
)

var domainList map[string]string

var wg sync.WaitGroup

func startClient(authKey, webSocketAddr string) {
	domainList = make(map[string]string)
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
				// CTL
				if mt == websocket.TextMessage {
					dlhps := ""
					for _, dlhp := range strings.Split(string(msg), "\n") {
						d := (strings.Split(dlhp, "|"))[0]
						lhp := (strings.Split(dlhp, "|"))[1]
						domainList[d] = lhp
						dlhps += "\thttp://" + d + " -> " + lhp + "\n"
					}

					log.Info("\nForwarding list:\n%s", dlhps)
				}
				// DATA
				if mt == websocket.BinaryMessage {
					headers := strings.Split(string(msg), "\n")

					id := headers[0]

					host := strings.Split(headers[2], ":")
					domain := strings.TrimSpace(host[1])
					localHostPort := domainList[domain]
					host[2] = (strings.Split(localHostPort, ":"))[1]
					headers[2] = strings.Join(host, ":")

					requestHeaders := headers[1:]
					message := []byte(strings.Join(requestHeaders, "\n"))

					log.Debug("\n<--Http request message-->\n%s\n<--Http request message-->", message)
					resp := httpClient.OpenUrl(&message, &localHostPort)

					newResp := string(id) + "\n" + string(resp)
					err = wsClient.Conn.WriteMessage(mt, []byte(newResp))
					if err != nil {
						log.Warn("Failed to write websocket: %v", err)
						break
					}
				}
			}
		} else {
			log.Error("Auth failed")
		}
	} else {
		log.Error("Failed to connect server")
	}
}
