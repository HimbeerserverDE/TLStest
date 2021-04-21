package main

import (
	"log"
	"net"
)

const host = "127.0.0.1:20202"

func handleClient(conn net.Conn) {
	log.Print("Server accepted client at ", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
	}
	buf = buf[:n]

	log.Print("Server received packet from ", conn.RemoteAddr().String(), ": ", string(buf))
}

func main() {
	log.Print("Starting server on ", host)

	l, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Print(err)
				continue
			}

			go handleClient(c)
		}
	}()

	log.Print("Server started successfully")
	log.Print("Starting client")

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Print("Client connection successful")

	if _, err = conn.Write([]byte("hello")); err != nil {
		log.Print(err)
	}

	for {
	}
}
