package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/subrag/kv-store/config"
)

func Server() {
	Addr := fmt.Sprintf("%v:%v", config.Host, config.Port)

	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("Welcome to kv-store!\n")
	fmt.Printf("kv-store running at %v\n", Addr)
	client := 0

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Process interrupted exiting!!!")
		listen.Close()
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("Error: %v\nClosing tcp listener...", err)
			listen.Close()
		}
		client += 1
		log.Printf("creating new connection with client [%v]!!!\n", client)

		go handleRequest(conn, client)

	}
}

func handleRequest(conn net.Conn, client int) {
	for {
		var buf []byte = make([]byte, 512)
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("Client %v got handler. Error: %v]\n", client, err)
			return
		}
		_ = writeResponse(conn, string(buf[:n]))

	}

}

func writeResponse(c net.Conn, d string) error {
	resStr := fmt.Sprintf("%v", d)
	if _, err := c.Write([]byte(resStr)); err != nil {
		return err
	}
	return nil
}
