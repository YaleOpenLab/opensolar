package main

import (
	"log"

	rpc "github.com/YaleOpenLab/opensolar/rpc"
)

func main() {
	port := 8002
	insecure := true

	log.Println("Starting opensolar")
	rpc.StartServer(port, insecure)
}
