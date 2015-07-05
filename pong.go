package main

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
		log.Println("[ERRO]", "can not read")
		fmt.Println("[ERRO]", "can not read")
		return
	}

	writer := bufio.NewWriter(conn)

	if writer == nil {
		log.Println("[ERRO]", "can not write")
		fmt.Println("[ERRO]", "can not write")
		return
	}

	for {
		// ping
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Println("[ERRO]", err)
			fmt.Println("[ERRO]", err)
			return
		}

		log.Println("[DATA]", remote, "=>", local, line)
		fmt.Println("[DATA]", remote, "=>", local, line)

		// pong
		fmt.Fprintln(writer, "[DATA]", local, "=>", remote, "ping", now)
		writer.Flush() // Don't forget to flush!

		log.Println("[DATA]", local, "=>", remote, "ping", now)
		fmt.Println("[DATA]", local, "=>", remote, "ping", now)

	}
}
