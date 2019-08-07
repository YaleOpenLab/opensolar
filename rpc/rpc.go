package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// also evaluate on how easy it would be to rewrite this in nodeJS since the
// frontend is in react. Not many advantages per se and this works fine, so I guess
// we'll stay with this one for a while
import (
	"encoding/json"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

func relayGetRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// validate if the person requesting this is a vlaid user on the platform
		err := erpc.CheckGet(w, r) // check origin of request as well if needed
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

// StartServer runs on the server side ie the server with the frontend.
// having to define specific endpoints for this because this
// is the system that would be used by the backend, so has to be built secure.
func StartServer(portx int, insecure bool) {
	// we have a sub handlers for each major entity. These handlers
	// call the relevant internal endpoints and return a erpc.StatusResponse message.
	// we also have to process data from the pi itself, and that should have its own
	// functions somewhere else that can be accessed by the rpc.

	// also, this is assumed to run on localhost and hence has no authentication mehcanism.
	// in the case we want to expose the API, we must add some stuff that secures this.
	// right now, its just the CORS header, since we want to allow all localhost processes
	// to access the API
	// a potential improvement will be to add something like macaroons
	// so that we can serve over an authenticated channel
	// setup all related handlers
	erpc.SetupPingHandler()
	relayGetRequest()
	setupProjectRPCs()
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
