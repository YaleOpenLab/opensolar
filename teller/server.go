package main

import (
	//"fmt"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// HCHeaderResponse defines the hash chain header's response
type HCHeaderResponse struct {
	Hash string
}

// hashChainHeaderHandler returns the header of the ipfs hash chain for
// clients who want historical record of all activities. This avoids a need for
// a direct endpoint that will serve data directly.
func hashChainHeaderHandler() {
	http.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}
		var x HCHeaderResponse
		x.Hash = HashChainHeader
		erpc.MarshalSend(w, x)
	})
}

func setupRoutes() {
	erpc.SetupDefaultHandler()
	hashChainHeaderHandler()
}

// curl https://localhost/ping --insecure {"Code":200,"Status":""}
// generate your own ssl certificate from letsencrypt or something to make sure the teller API calls
// are accessible frmo outside localhost
func startServer(port int) {
	setupRoutes()

	portString, err := utils.ToString(port)
	if err != nil {
		log.Println("couldn't parse port, setting it to 80 by default")
		portString = "80"
	}

	log.Fatal(http.ListenAndServeTLS(":"+portString, "ssl/server.crt", "ssl/server.key", nil))
}
