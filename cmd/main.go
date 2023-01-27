package main

import (
	"log"
	"os"
	"netcat/server"
)

func main() {
	port := "9090" // порт по дефолту
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if !server.PortCheck(port) {
		log.Println("[USAGE]: ./TCPChat $port")
		return
	}
	client := server.NewServer()
	client.Run(port)
}
