package main

import (
	"github.com/subrag/kv-store/core"
	"github.com/subrag/kv-store/server"
)

func main() {
	// server.Server()
	db := core.KV{}
	server.AyncServer(&db)
}
