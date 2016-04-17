package main

import (
	"runtime"
	"ultraman/module/client"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	client.Run()
}
