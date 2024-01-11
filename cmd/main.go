package main

import (
	"log"
	"velk2/internal/velk"
)

func main() {
	addr := "localhost:9999"
	server, err := velk.NewVelkServer(addr)
	if err != nil {
		log.Fatalf("Failed to start server at %s", addr)
	}
	log.Println("Server is running on:", server.Addr)
	server.GameLoop()
}
