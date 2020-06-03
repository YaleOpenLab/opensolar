package rpc

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	openx "github.com/YaleOpenLab/openx/database"
	// openxrpc "github.com/YaleOpenLab/openx/rpc"
)

// setupUserRpcs sets up user related RPCs
func setupUserRpcs() {
	updateUser()
	reportProject()
	userInfo()
	registerUser()
	getUserRoles()
}

// UserRPC is a collection of all user RPC endpoints and their required params
var UserRPC = map[int][]string{
	1: {"/update", "POST"},                                                  // POST
	2: {"/user/report", "POST", "projIndex"},                                // POST
	3: {"/user/info", "GET"},                                                // GET
	4: {"/user/register", "POST", "email", "username", "pwhash", "seedpwd"}, // POST
	5: {"/user/roles", "GET"},                                               // GET
}

func userValidateHelper(w http.ResponseWriter, r *http.Request, options []string, method string) (openx.User, error) {
	var user openx.User

	err := checkReqdParams(w, r, options, method)
	if erpc.Err(w, err, erpc.StatusUnauthorized, "", messages.NotUserError) {
		return user, err
	}

	var username, token string
	if method == "GET" {
		username, token = r.URL.Query()["username"][0], r.URL.Query()["token"][0]
	} else {
		username, token = r.FormValue("username"), r.FormValue("token")
	}

	user, err = core.ValidateUser(username, token)
	if erpc.Err(w, err, erpc.StatusUnauthorized, "", messages.NotUserError) {
		return user, err
	}

	return user, nil
}

// updateUser updates credentials of the user
func updateUser() {
	http.HandleFunc(UserRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		// updateUser must first call the openx rpc to update the user struct

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
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var user openx.User
		err = json.Unmarshal(data, &user)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		token := r.FormValue("token")

		if user.Index != 0 {
			// check whether given user is an investor or recipient
			investor, err := core.ValidateInvestor(user.Username, token)
			if err == nil {
				investor.U = &user
				err = investor.Save()
				if erpc.Err(w, err, erpc.StatusInternalServerError, "unable to save investor") {
					return
				}
			}
			recipient, err := core.ValidateRecipient(user.Username, token)
			if err == nil {
				recipient.U = &user
				err = recipient.Save()
				if erpc.Err(w, err, erpc.StatusInternalServerError, "unable to save recipient") {
					return
				}
			}
			entity, err := core.ValidateEntity(user.Username, token)
			if err == nil {
				entity.U = &user
				err = entity.Save()
				if erpc.Err(w, err, erpc.StatusInternalServerError, "unable to save entity") {
					return
				}
			}
			erpc.MarshalSend(w, user)
		} else {
			log.Println("user not updated")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
	})
}

// reportProject updates credentials of the user
func reportProject() {
	http.HandleFunc(UserRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[1][2:], UserRPC[1][1])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		err = core.UserMarkFlagged(projIndex, user.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
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
		// need to pass the pwhash param here
		prepUser, err := userValidateHelper(w, r, UserRPC[3][2:], UserRPC[3][1])
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

// registerUser creates a new user
func registerUser() {
	http.HandleFunc(UserRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		// parse form the check whether required params are present
		err = r.ParseForm()
		if erpc.Err(w, err, erpc.StatusUnauthorized) {
			return
		}

		for _, option := range UserRPC[4][2:] {
			if r.FormValue(option) == "" {
				log.Println("required param: ", option, " not found")
				erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.ParamError(option))
				return
			}
		}

		email := r.FormValue("email")
		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")
		seedpwd := r.FormValue("seedpwd")

		user, err := core.NewUser(username, pwhash, seedpwd, email)
		if erpc.Err(w, err, erpc.StatusNotFound, "unable to save user") {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// UserRoleStruct is a handler used to return /user/validate
type UserRoleStruct struct {
	User      openx.User
	Investor  core.Investor
	Recipient core.Recipient
	Entity    core.Entity
}

// getUserRoles gets a list of the roles that an investor partakes
func getUserRoles() {
	http.HandleFunc(UserRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		user, err := userValidateHelper(w, r, UserRPC[5][2:], UserRPC[5][1])
		if err != nil {
			return
		}

		var ret UserRoleStruct
		// we now have the user struct, search invetors, recipients, entities for what role the
		// user is

		ret.User = user
		inv, err := core.SearchForInvestor(user.Username)
		if err == nil {
			ret.Investor = inv
		}

		recp, err := core.SearchForRecipient(user.Username)
		if err == nil {
			ret.Recipient = recp
		}

		entity, err := core.SearchForEntity(user.Username)
		if err == nil {
			ret.Entity = entity
		}

		erpc.MarshalSend(w, ret)
	})
}
