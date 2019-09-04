package rpc

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

func checkReqdParams(r *http.Request, options ...string) error {

	if r.URL.Query() == nil {
		return prepUser, errors.New("url query can't be empty")
	}

	options = append(options, "username", "token") // default for all endpoints

	for _, option := range options {
		if r.URL.Query()[option] == nil {
			return errors.New("required param: " + option + " not specified, quitting")
		}
	}

	if len(r.URL.Query()["token"][0]) != consts.AccessTokenLength {
		return prepUser, errors.New("pwhash length not 128, quitting")
	}

	return nil
}

// relayGetRequest relays get requests to openx
func relayRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := erpc.CheckGet(w, r)
			if err != nil {
				log.Println(err)
				return
			}
			body := consts.OpenxURL + r.URL.String()
			log.Println(body)
			data, err := erpc.GetRequest(body)
			if err != nil {
				log.Println("could not relay get request", err)
				erpc.ResponseHandler(w, http.StatusInternalServerError)
				return
			}

			var x interface{}
			_ = json.Unmarshal(data, &x)
			erpc.MarshalSend(w, x)
		} else if r.Method == "POST" {
			err := erpc.CheckPost(w, r)
			if err != nil {
				log.Println(err)
				return
			}
			body := consts.OpenxURL + r.URL.String()
			log.Println(body)

			err = r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			data, err := erpc.PostForm(body, r.Form)
			var x interface{}
			_ = json.Unmarshal(data, &x)
			erpc.MarshalSend(w, x)
		}
	})
}

// StartServer starts the opensolar backend server
func StartServer(portx int, insecure bool) {
	erpc.SetupPingHandler()
	relayRequest()
	setupProjectRPCs()
	setupUserRpcs()
	setupInvestorRPCs()
	setupRecipientRPCs()
	setupPublicRoutes()
	setupEntityRPCs()
	setupParticleHandlers()
	setupSwytchApis()
	setupStagesHandlers()
	setupAdminHandlers()

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
