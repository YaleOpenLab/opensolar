package rpc

import (
	"encoding/json"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

// relayGetRequest relays get requests to openx
func relayGetRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		body := consts.OpenxURL + r.URL.String()
		log.Println(body)
		data, err := erpc.GetRequest(body)
		if err != nil {
			log.Println("could not submit transacation to testnet, quitting")
			erpc.ResponseHandler(w, http.StatusInternalServerError)
			return
		}

		var x interface{}
		_ = json.Unmarshal(data, &x)
		erpc.MarshalSend(w, x)
	})
}

// StartServer starts the opensolar backend server
func StartServer(portx int, insecure bool) {
	erpc.SetupPingHandler()
	relayGetRequest()
	setupProjectRPCs()
	setupUserRpcs()
	setupInvestorRPCs()
	setupRecipientRPCs()
	setupPublicRoutes()
	setupEntityRPCs()
	setupParticleHandlers()
	setupSwytchApis()
	setupStagesHandlers()
	// adminHandlers()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}
