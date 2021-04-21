package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

const host = "127.0.0.1:20203"

func handleClient(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()

	log.Print("Client at ", addr, " connected")

	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}

		log.Print("Server received packet from ", addr, ": ", msg)

		if _, err = conn.Write([]byte("Hello client\n")); err != nil {
			log.Print(err)
			return
		}
	}
}

func serve(ready chan struct{}) {
	crt, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Print(err)
		return
	}

	conf := &tls.Config{
		Certificates:             []tls.Certificate{crt},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites:             []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	l, err := tls.Listen("tcp", host, conf)
	if err != nil {
		log.Print(err)
		return
	}
	defer l.Close()

	log.Print("Server listening on ", host)

	close(ready)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleClient(conn)
	}
}

func connect(ready chan struct{}) {
	<-ready

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", host, conf)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	log.Print("Client connection successful")

	if _, err = conn.Write([]byte("Hello server\n")); err != nil {
		log.Print(err)
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(err)
	}
	buf = buf[:n]

	log.Print("Client received packet: ", string(buf))
}

func main() {
	ready := make(chan struct{})

	go serve(ready)
	connect(ready)
}
