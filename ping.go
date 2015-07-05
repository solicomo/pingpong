package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
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
		self.ping()
	}()

	wg.Wait()
}

func (self *Ping) connect() (err error) {

	self.conn, err = net.Dial(self.config.Type, self.config.Host+":"+self.config.Port)

	return
}

func (self *Ping) read() {

	reader := bufio.NewReader(self.conn)

	if reader == nil {
		log.Println("[ERRO]", "can not read from", self.config.Host+":"+self.config.Port)
		fmt.Println("[ERRO]", "can not read from", self.config.Host+":"+self.config.Port)
		return
	}

	remote := self.conn.RemoteAddr().String()
	local := self.conn.LocalAddr().String()

	for {
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
	}
}

func (self *Ping) ping() {

	writer := bufio.NewWriter(self.conn)

	if writer == nil {
		log.Println("[ERRO]", "can not write to", self.config.Host+":"+self.config.Port)
		fmt.Println("[ERRO]", "can not write to", self.config.Host+":"+self.config.Port)
		return
	}

	remote := self.conn.RemoteAddr().String()
	local := self.conn.LocalAddr().String()

	for now := range time.Tick(1 * time.Second) {

		line := strconv.FormatInt(now.Unix(), 10)

		fmt.Fprintln(writer, line)
		writer.Flush() // Don't forget to flush!

		log.Println("[SEND]", local, "=>", remote, ":", line)
		fmt.Println("[SEND]", local, "=>", remote, ":", line)
	}
}
