package main

import (
	"flag"

	"github.com/subrag/kv-store/core"
	"github.com/subrag/kv-store/server"
)

func main() {
	var Port int
	db := core.DB{KV: core.KV{}}
	go db.ScheduledEvict()
	flag.IntVar(&Port, "p", 8987, "Specify the port to run kv-store")
	flag.Parse()
	// server.Server(Port)
	server.AyncServer(&db, Port)
}
