package ssdb

import (
	"os"
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
		log.Error("Failed to connect ssdb server: %v", err)
		os.Exit(1)
	}

	gossdb.Encoding = true

	client, err := pool.NewClient()
	if err != nil {
		log.Error("Failed to connect ssdb server: %v", err)
		os.Exit(1)
	}

	log.Info("Connecting ssdb server success")
	return client
}
