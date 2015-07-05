package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Pong struct {
	App
	listener net.Listener
	clients  []net.Conn
}

func (self *Pong) Run() {

	err := self.listen()

	if err != nil {
		log.Println("[ERRO]", err)
		fmt.Println("[ERRO]", err)
		return
	}

	self.clients = make([]net.Conn, 0)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		self.accept()
	}()

	wg.Wait()

	for _, conn := range self.clients {

		conn.Close()
	}
}

func (self *Pong) listen() (err error) {

	self.listener, err = net.Listen(self.config.Type, self.config.Host+":"+self.config.Port)

	return
}

func (self *Pong) accept() {

	for {

		conn, err := self.listener.Accept()
		if err != nil {
			log.Println("[ERRO]", err)
			fmt.Println("[ERRO]", err)
			return
		}

		self.clients = append(self.clients, conn)

		go self.pong(conn)
	}
}

func (self *Pong) pong(conn net.Conn) {

	remote := conn.RemoteAddr().String()
	local := conn.LocalAddr().String()

	reader := bufio.NewReader(conn)

	if reader == nil {
		log.Println("[ERRO]", "can not read from", remote)
		fmt.Println("[ERRO]", "can not read from", remote)
		return
	}

	writer := bufio.NewWriter(conn)

	if writer == nil {
		log.Println("[ERRO]", "can not write to", remote)
		fmt.Println("[ERRO]", "can not write to", remote)
		return
	}

	for {
		// ping
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Println("[INFO]", remote, "=>", local, ":", "closed")
			fmt.Println("[INFO]", remote, "=>", local, ":", "closed")
			return
		}

		if err != nil {
			log.Println("[ERRO]", remote, "=>", local, ":", err)
			fmt.Println("[ERRO]", remote, "=>", local, ":", err)
			return
		}

		log.Println("[RECV]", remote, "=>", local, ":", line)
		fmt.Println("[RECV]", remote, "=>", local, ":", line)

		// pong
		fmt.Fprint(writer, line)
		writer.Flush() // Don't forget to flush!

		log.Println("[SEND]", local, "=>", remote, ":", line)
		fmt.Println("[SEND]", local, "=>", remote, ":", line)

	}
}
