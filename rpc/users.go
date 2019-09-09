package rpc

import (
	"github.com/pkg/errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	openx "github.com/YaleOpenLab/openx/database"
	// openxrpc "github.com/YaleOpenLab/openx/rpc"
)

// setupUserRpcs sets up user related RPCs
func setupUserRpcs() {
	updateUser()
	reportProject()
}

// UserRPC is a collection of all user RPC endpoints and their required params
var UserRPC = map[int][]string{
	1: []string{"/user/update"},
	2: []string{"/user/report", "projIndex"},
}

func userValidateHelper(w http.ResponseWriter, r *http.Request, options []string) (openx.User, error) {
	var user openx.User

	err := checkReqdParams(w, r, options)
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return user, err
	}

	username := r.URL.Query()["username"][0]
	token := r.URL.Query()["token"][0]

	user, err = core.ValidateUser(username, token)
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return user, err
	}

	if !user.Admin {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return user, errors.New("unauthorized")
	}

	return user, nil
}

// updateUser updates credentials of the user
func updateUser() {
	http.HandleFunc(UserRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := userValidateHelper(w, r, UserRPC[1][1:])
		if err != nil {
			return
		}

		if r.FormValue("name") != "" {
			user.Name = r.FormValue("name")
		}
		if r.FormValue("city") != "" {
			user.City = r.FormValue("city")
		}
		if r.FormValue("zipcode") != "" {
			user.ZipCode = r.FormValue("zipcode")
		}
		if r.FormValue("country") != "" {
			user.Country = r.FormValue("country")
		}
		if r.FormValue("recoveryphone") != "" {
			user.RecoveryPhone = r.FormValue("recoveryphone")
		}
		if r.FormValue("address") != "" {
			user.Address = r.FormValue("address")
		}
		if r.FormValue("description") != "" {
			user.Description = r.FormValue("description")
		}
		if r.FormValue("email") != "" {
			user.Email = r.FormValue("email")
		}
		if r.FormValue("notification") != "" {
			if r.FormValue("notification") != "true" {
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
		recipient, err := recpValidateHelper(w, r, UserRPC[1][1:])
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
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := userValidateHelper(w, r, UserRPC[1][1:])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")

		projIndex, err := utils.ToInt(projIndexx)
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
