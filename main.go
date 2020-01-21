package main

import (
	"fmt"
	"io"
	"net"
	"log"
)

var (
	counter int

	// TODO configurable
	listenAddress = "localhost:8080"

	// TODO configurable
	server = []string {
		"localhost:5000",
		"localhost:5001",
		"localhost:5002",
	}
)

func main()  {
	listener, err := net.Listen("tcp", listenAddress)

	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	defer listener.Close()

	for   {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
		}

		backend := chooseBackend()
		fmt.Printf("counter: %d backend:%s\n", counter, backend)

		go func() {
			err := proxy(backend, conn)
			if err != nil {
				log.Printf("WARNING: proxying failed: %v", err)
			}
		}()
	}
}

func proxy(backend string, connection net.Conn) error {
	backendConnection, err := net.Dial("tcp", backend)

	if err != nil {
		return fmt.Errorf("Failed to connect to backend %s: %v", backend, err)
	}

	// connection => backendConnection
	go io.Copy(backendConnection, connection)

	// backendConnection => connection
	go io.Copy(connection, backendConnection)

	return nil
}

func chooseBackend()  string {
	//robin
	s := server[counter % len(server)]
	counter++
	return s
}
