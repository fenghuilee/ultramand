package client

import (
	"ultraman/lib/log"
)

func Run() {
	// parse options
	cliOpts := CliRunArgs()

	// init logging
	log.LogTo(cliOpts.logTo, cliOpts.logLevel)

	// start http/websocket server listeners
	startClient(cliOpts.authKey, cliOpts.webSocket)
}
