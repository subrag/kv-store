package main

import (
	"github.com/subrag/kv-store/core"
	"github.com/subrag/kv-store/server"
)

func main() {
	// server.Server()
	db := core.DB{KV: core.KV{}}
	go db.ScheduledEvict()
	server.AyncServer(&db)
}
