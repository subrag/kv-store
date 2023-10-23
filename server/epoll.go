// https://gist.github.com/tevino/3a4f4ec4ea9d0ca66d4f
// detail on epoll: https://medium.com/@avocadi/what-is-epoll-9bbc74272f7c#:~:text=Internally%2C when the process calls,faster than poll or select
package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/subrag/kv-store/config"
	"github.com/subrag/kv-store/core"
)

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 32
)

// graceful close of kv-store and connected client on keyboard interruption.
func gracefulCloseServerfd(serverFD int) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // wait block
	log.Println("process interrupted exiting...")
	syscall.Close(serverFD)
	os.Exit(1)
}

// create async server using epoll.
func AyncServer(db *core.DB, port int) {
	const maxClients = 20000

	var event syscall.EpollEvent
	var events [maxClients]syscall.EpollEvent

	// creates server file descriptor.
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer syscall.Close(serverFD)
	// track and close OS interrupt
	go gracefulCloseServerfd(serverFD)

	// set non block for serverFD.
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		fmt.Println("setnonblock1: ", err)
		return
	}

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(config.Host).To4())

	if err = syscall.Bind(serverFD, &addr); err != nil {
		log.Printf("unable to bind %v:%v, %s", config.Host, port, err.Error())
		return
	}

	if err = syscall.Listen(serverFD, maxClients); err != nil {
		log.Print("unable to start server listen.\n")
		return
	}

	// create epoll fd.
	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		log.Printf("error in epoll_create: %s\n", e.Error())
		return
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN
	event.Fd = int32(serverFD)
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, serverFD, &event); e != nil {
		log.Println("epoll_ctl: ", e)
		return
	}
	fmt.Printf("kv-store running %v:%v\n", config.Host, port)

	for {
		nevents, e := syscall.EpollWait(epfd, events[:], -1)
		if e != nil {
			log.Println("epoll_wait: ", e)
			continue
		}

		for ev := 0; ev < nevents; ev++ {
			// Check events from serverFD Accept new connection from ServerFD.
			if int(events[ev].Fd) == serverFD {
				connFd, _, err := syscall.Accept(serverFD)
				if err != nil {
					fmt.Println("accept: ", err)
					continue
				}
				_ = syscall.SetNonblock(serverFD, true)
				event.Events = syscall.EPOLLIN | EPOLLET
				event.Fd = int32(connFd)
				// Note:
				// As process is enrolling fd to epoll instance for events EPOLLIN or EPOLLET
				// for all child epoll fd's event type remains same. Main serverFD is alredy registered above.

				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
					log.Panicf("unable to register epoll fd to epoll instance, details: %v, %v", connFd, err)
				}
			} else {
				// Process read from child epoll.
				var buf = make([]byte, 512)
				nbytes, e := syscall.Read(int(events[ev].Fd), buf)
				if e != nil {
					syscall.Close(int(events[ev].Fd))
					continue
				}
				val, err := core.HandlerQuery(buf[:nbytes], db)
				if err != nil {
					log.Printf("%v\n", err)
				}
				_, err = syscall.Write(int(events[ev].Fd), []byte(val))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
