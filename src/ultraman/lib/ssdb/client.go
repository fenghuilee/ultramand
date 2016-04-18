package ssdb

import (
	"ultraman/lib/log"

	"github.com/seefan/gossdb"
)

func Run(host string, port int) *gossdb.Client {
	log.Info("Connecting ssdb server %v:%v", host, port)

	pool, err := gossdb.NewPool(&gossdb.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		panic(log.Error("Failed to connect ssdb server: %v", err))
	}

	gossdb.Encoding = true

	client, err := pool.NewClient()
	if err != nil {
		panic(log.Error("Failed to connect ssdb server: %v", err))
	}

	log.Info("Connecting ssdb server success")
	return client
}
