package server

import (
	"flag"
	//"fmt"
	"os"

	//"github.com/widuu/goini"
)

func CliRunArgs() *CliOptions {

	//conf := flag.String("conf", "./conf/server.ini", "Path to config file")

	domain := flag.String("domain", "tunnel.preruntime.com", "Domain where the tunnels are hosted")
	http := flag.String("http", ":8000", "Public address for HTTP connections, empty string to disable")
	webSocket := flag.String("websocket", ":4443", "Public address listening for ngrok client")
	logTo := flag.String("log-to", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	logLevel := flag.String("log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	flag.Parse()

	return &CliOptions{
		domain:    *domain,
		http:      *http,
		webSocket: *webSocket,
		logTo:     *logTo,
		logLevel:  *logLevel,
	}

	//	if PathExist(*conf) == false {
	//		fmt.Printf("Error: config file no found")
	//		os.Exit(1)
	//	}
	//	ini := goini.SetConfig(*conf)

	//	domain := ini.GetValue("main", "domain")
	//	http := ini.GetValue("http", "host") + ":" + ini.GetValue("http", "port")
	//	webSocket := ini.GetValue("websocket", "host") + ":" + ini.GetValue("websocket", "port")
	//	logTo := ini.GetValue("log", "to")
	//	logLevel := ini.GetValue("log", "level")

	//	return &CliOptions{
	//		domain:    domain,
	//		http:      http,
	//		webSocket: webSocket,
	//		logTo:     logTo,
	//		logLevel:  logLevel,
	//	}
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
