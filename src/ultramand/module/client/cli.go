package client

import (
	"flag"

	//"github.com/widuu/goini"
)

func CliRunArgs() *CliOptions {
	authKey := flag.String("auth-key", "fenghuilee:8b83ae27-a952-41d2-ac86-577a399e4cc9", "Auth key")
	webSocket := flag.String("websocket", "ws.preruntime.com:4443", "Public address listening for ngrok client")
	logTo := flag.String("log-to", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	logLevel := flag.String("log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	flag.Parse()

	return &CliOptions{
		authKey:   *authKey,
		webSocket: *webSocket,
		logTo:     *logTo,
		logLevel:  *logLevel,
	}
}
