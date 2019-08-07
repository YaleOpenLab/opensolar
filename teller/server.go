package main

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// server starts a local server which would inform us about the uptime of the teller and provide a data endpoint

func checkGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func responseHandler(w http.ResponseWriter, r *http.Request, status int) {
	var response erpc.StatusResponse
	response.Code = status
	switch status {
	case erpc.StatusOK:
		response.Status = "OK"
	case erpc.StatusNotFound:
		response.Status = "404 Error Not Found!"
	case erpc.StatusInternalServerError:
		response.Status = "Internal Server Error"
	}
	erpc.MarshalSend(w, response)
}

// HCHeaderResponse defines the hash chain header's response
type HCHeaderResponse struct {
	Hash string
}

// hashChainHeaderHandler returns the header of the ipfs hash chain
// clients who want historicasl record of all activities can record the latest hash
// and then derive all the other files from it. This avoids a need for a direct endpoint
// that will serve data directly while leveraging ipfs.
func hashChainHeaderHandler() {
	http.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		var x HCHeaderResponse
		x.Hash = HashChainHeader
		xJson, err := json.Marshal(x)
		if err != nil {
			responseHandler(w, r, erpc.StatusInternalServerError)
			return
		}
		WriteToHandler(w, xJson)
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
		log.Fatal(err)
	}

	err = http.ListenAndServeTLS(":"+portString, "ssl/server.crt", "ssl/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
