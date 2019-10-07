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
	userInfo()
}

// UserRPC is a collection of all user RPC endpoints and their required params
var UserRPC = map[int][]string{
	1: []string{"/user/update", "POST"},              // POST
	2: []string{"/user/report", "POST", "projIndex"}, // POST
	3: []string{"/user/info", "GET"},                // GET
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

		user, err := userValidateHelper(w, r, UserRPC[1][2:])
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
		investor, err := InvValidateHelper(w, r, UserRPC[1][2:])
		if err == nil {
			investor.U = &user
			err = investor.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		}
		recipient, err := recpValidateHelper(w, r, UserRPC[1][2:])
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

		user, err := userValidateHelper(w, r, UserRPC[1][2:])
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

// validateParams is a struct used fro validating user params
type validateParams struct {
	// Role is a string identifying the user on the pilot opensolar platform
	Role string
	// Entity is an interface containing the user struct
	Entity interface{}
}

// userInfo validates a user and returns whether the user is an investor or recipient on the opensolar platform
func userInfo() {
	http.HandleFunc(UserRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		// need to pass the pwhash param here
		prepUser, err := userValidateHelper(w, r, UserRPC[3][2:])
		if err != nil {
			return
		}
		// no we need to see whether this guy is an investor or a recipient.
		var prepInvestor core.Investor
		var prepRecipient core.Recipient
		var prepEntity core.Entity

		var x validateParams

		prepInvestor, err = core.RetrieveInvestor(prepUser.Index)
		if err == nil && prepInvestor.U.Index != 0 {
			x.Role = "Investor"
			x.Entity = prepInvestor
			erpc.MarshalSend(w, x)
			return
		}

		prepRecipient, err = core.RetrieveRecipient(prepUser.Index)
		if err == nil && prepRecipient.U.Index != 0 {
			x.Role = "Recipient"
			x.Entity = prepRecipient
			erpc.MarshalSend(w, x)
			return
		}

		prepEntity, err = core.RetrieveEntity(prepUser.Index)
		if err == nil && prepEntity.U.Index != 0 {
			x.Role = "Entity"
			x.Entity = prepEntity
			erpc.MarshalSend(w, x)
			return
		}
	})
}
