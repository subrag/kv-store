package main

import (
	"flag"

	"github.com/subrag/kv-store/core"
	"github.com/subrag/kv-store/server"
)

func main() {
	var Port int
	// server.Server()
	db := core.DB{KV: core.KV{}}
	go db.ScheduledEvict()
	flag.IntVar(&Port, "p", 9089, "Specify the port to run kv-store")
	server.AyncServer(&db, Port)
}
