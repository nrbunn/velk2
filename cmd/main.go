package main

import (
	"log"
	"velk2/pkg/structs"
)

func main() {
	addr := "localhost:9999"
	server, err := structs.NewVelkServer(addr)
	if err != nil {
		log.Fatalf("Failed to start server at %s", addr)
	}
	log.Println("Server is running on:", server.Addr)
	server.GameLoop()
}
