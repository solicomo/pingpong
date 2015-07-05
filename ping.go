package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type Ping struct {
	App
	conn net.Conn
}

func (self *Ping) Run() {

	err := self.connect()

	if err != nil {
		log.Println("[ERRO]", err)
		fmt.Println("[ERRO]", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		self.read()
	}()

	go func() {
		defer wg.Done()
		v.ping()
	}()
}

func (self *Ping) connect() (err error) {

	self.conn, err = net.Dial(self.config.Type, self.config.Host+":"+self.config.Port)

	return
}

func (self *Ping) read() {

	reader := bufio.NewReader(self.conn)

	if reader == nil {
		log.Println("[ERRO]", "can not read")
		fmt.Println("[ERRO]", "can not read")
		return
	}

	remote := self.conn.RemoteAddr().String()
	local := self.conn.LocalAddr().String()

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Println("[ERRO]", err)
			fmt.Println("[ERRO]", err)
			return
		}

		log.Println("[DATA]", remote, "=>", local, line)
		fmt.Println("[DATA]", remote, "=>", local, line)
	}
}

func (self *Ping) ping() {

	writer := bufio.NewWriter(self.conn)

	if writer == nil {
		log.Println("[ERRO]", "can not write")
		fmt.Println("[ERRO]", "can not write")
		return
	}

	remote := self.RemoteAddr().String()
	local := self.LocalAddr().String()

	for now := range time.Tick(1 * time.Second) {

		fmt.Fprintln(writer, "[DATA]", local, "=>", remote, "ping", now)
		writer.Flush() // Don't forget to flush!

		log.Println("[DATA]", local, "=>", remote, "ping", now)
		fmt.Println("[DATA]", local, "=>", remote, "ping", now)
	}
}
