package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8587")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			log.Println("accept1")
			if err != nil {
				//handle error
				panic(err)
			}
			go handleConection1(conn)
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			log.Println("accept")
			if err != nil {
				//handle error
				panic(err)
			}
			go handleConection(conn)
		}
	}()

	for {
	}
}
func handleConection1(conn net.Conn) {
	log.Println("handleConection1")
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("read error", err)
			return
		}
		if n == 0 {
			return
		}
		log.Printf("wtf!!!! received from %v: %s", conn.RemoteAddr(), string(buf[:n]))
	}
}

func handleConection(conn net.Conn) {
	log.Println("handleConection")
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("read error", err)
			return
		}
		if n == 0 {
			return
		}
		log.Printf("received from %v: %s", conn.RemoteAddr(), string(buf[:n]))
	}
}
