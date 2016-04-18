package main

import (
	"runtime"
	"ultramand/module/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	server.Run()
}
