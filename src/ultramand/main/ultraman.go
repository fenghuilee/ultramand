package main

import (
	"runtime"
	"ultramand/module/client"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	client.Run()
}
