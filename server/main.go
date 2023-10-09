package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/subrag/kv-store/config"
	"github.com/subrag/kv-store/core"
)

func Server() {
	Addr := fmt.Sprintf("%v:%v", config.Host, config.Port)

	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("kv-store running: %v\n", Addr)

	db := core.KV{}
	db["a"] = "a"
	client := 0
	go gracefullShutdown(listen)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("error: %v\n\nclosing tcp listener...", err)
			listen.Close()
		}
		client += 1
		log.Printf("client %v connected.\n", client)

		go handleRequest(conn, client, &db)

	}
}
