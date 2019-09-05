package rpc

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
)

func checkReqdParams(w http.ResponseWriter, r *http.Request, options []string) error {

	if r.Method == "GET" {
		if r.URL.Query() == nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("url query can't be empty")
		}

		options = append(options, "username", "token") // default for all endpoints

		for _, option := range options {
			if r.URL.Query()[option] == nil {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return errors.New("required param: " + option + " not specified, quitting")
			}
		}

		if len(r.URL.Query()["token"][0]) != 32 {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("token length not 32, quitting")
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return err
		}

		if r.FormValue("username") == "" || r.FormValue("token") == "" {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("required params username or token missing")
		}

		if len(r.FormValue("token")) != 32 {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return errors.New("token length not 32, quitting")
		}

		for _, option := range options {
			if r.FormValue(option) == "" {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return errors.New("required param: " + option + " not specified, quitting")
			}
		}
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
				log.Println(err)
				return
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
		log.Println("couldn't parse passed port, setting it to default 80")
		port = "80"
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}
