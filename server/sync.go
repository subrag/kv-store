package server

import (
	"fmt"
	"log"
	"net"

	"github.com/subrag/kv-store/config"
	"github.com/subrag/kv-store/core"
)

func Server(port int) {
	// Server spinsup n goroutines for n tcp requests.
	Addr := fmt.Sprintf("%v:%v", config.Host, port)

	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("kv-store running: %v\n", Addr)

	db := core.DB{}
	db.KV["a"] = core.Item{}
	client := 0
	go gracefulShutdown(listen)

	for {
		conn, err := listen.Accept()
		if err != nil {
			listen.Close()
			log.Fatalf("error: %v\n\nclosing tcp listener...", err)
		}
		client++
		log.Printf("client %v connected.\n", client)

		go handleRequest(conn, client, &db)
	}
}
