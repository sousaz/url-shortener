package main

import (
	"fmt"
	"net"

	"github.com/sousaz/urlshortener/database"
	"github.com/sousaz/urlshortener/socket"
)

func main() {
	database.InitDB()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error initializing server: %v\n", err)
		return
	}

	fmt.Println("Server running in port: 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Inicializa uma goroutine
		go socket.HandleConnection(conn)
	}
}
