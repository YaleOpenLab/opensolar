package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	openxrpc "github.com/YaleOpenLab/openx/rpc"
)

// UserRPC is a collection of all user RPC endpoints and their required params
var UserRPC = map[int][]string{
	1: []string{"/user/update"},
	2: []string{"/user/report", "projIndex"},
}

// setupUserRpcs sets up user related RPCs
func setupUserRpcs() {
	updateUser()
	reportProject()
}

// updateUser updates credentials of the user
func updateUser() {
	http.HandleFunc(UserRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		user, err := openxrpc.CheckReqdParams(w, r, UserRPC[1][1:])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		if r.URL.Query()["name"] != nil {
			user.Name = r.URL.Query()["name"][0]
		}
		if r.URL.Query()["city"] != nil {
			user.City = r.URL.Query()["city"][0]
		}
		if r.URL.Query()["zipcode"] != nil {
			user.ZipCode = r.URL.Query()["zipcode"][0]
		}
		if r.URL.Query()["country"] != nil {
			user.Country = r.URL.Query()["country"][0]
		}
		if r.URL.Query()["recoveryphone"] != nil {
			user.RecoveryPhone = r.URL.Query()["recoveryphone"][0]
		}
		if r.URL.Query()["address"] != nil {
			user.Address = r.URL.Query()["address"][0]
		}
		if r.URL.Query()["description"] != nil {
			user.Description = r.URL.Query()["description"][0]
		}
		if r.URL.Query()["email"] != nil {
			user.Email = r.URL.Query()["email"][0]
		}
		if r.URL.Query()["notification"] != nil {
			if r.URL.Query()["notification"][0] != "true" {
				user.Notification = false
			} else {
				user.Notification = true
			}
		}

		err = user.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// check whether given user is an investor or recipient
		investor, err := InvValidateHelper(w, r, UserRPC[1][1:])
		if err == nil {
			investor.U = &user
			err = investor.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		recipient, err := RecpValidateHelper(w, r, UserRPC[1][1:])
		if err == nil {
			recipient.U = &user
			err = recipient.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// reportProject updates credentials of the user
func reportProject() {
	http.HandleFunc(UserRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		user, err := openxrpc.CheckReqdParams(w, r, UserRPC[2][1:])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = core.UserMarkFlagged(projIndex, user.Index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
