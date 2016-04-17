package main

import (
	"runtime"
	"ultraman/module/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	server.Run()
}
