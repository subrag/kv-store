package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/subrag/kv-store/core"
)

// graceful close of kv-store and connected client on keyboard interruption.
func gracefulShutdown(listen net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // wait block
	log.Println("process interrupted exiting...")
	listen.Close()
}

// serve cli client.
func handleRequest(conn net.Conn, client int, db *core.DB) {
	for {
		var buf = make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		// command := string(buf[:n])

		val, err := core.HandlerQuery(buf[:n], db)

		if err != nil {
			log.Printf("%v\n", err)
		}
		err = writeResponse(conn, val)
		if err != nil {
			fmt.Printf("error while writing response clinet [%v]", client)
		}
	}
}

func writeResponse(c net.Conn, d string) error {
	resStr := fmt.Sprintf("%v", d)
	if _, err := c.Write([]byte(resStr)); err != nil {
		return err
	}
	return nil
}
