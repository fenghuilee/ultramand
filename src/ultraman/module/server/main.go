package server

import (
	"ultraman/lib/log"
)

func Run() {
	// parse options
	cliOpts := CliRunArgs()

	// init logging
	log.LogTo(cliOpts.logTo, cliOpts.logLevel)

	// start http/websocket server listeners
	StartServer(cliOpts.http, cliOpts.webSocket, cliOpts.ssdb)
}
