package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"

	"github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	core "github.com/YaleOpenLab/opensolar/core"
)

// setupEntityRPCs sets up all entity related endpoints
func setupEntityRPCs() {
	validateEntity()
	getStage0Contracts()
	getStage1Contracts()
	getStage2Contracts()
	addCollateral()
	proposeOpensolarProject()
	registerEntity()
}

var EntityRPC = map[int][]string{
	1: []string{"/entity/validate", "GET"},                                                                  // GET
	2: []string{"/entity/stage0", "GET"},                                                                    // GET
	3: []string{"/entity/stage1", "GET"},                                                                    // GET
	4: []string{"/entity/stage2", "GET"},                                                                    // GET
	5: []string{"/entity/addcollateral", "POST", "amount", "collateral"},                                    // POST
	6: []string{"/entity/proposeproject/opensolar", "POST", "projIndex", "fee"},                             // POST
	7: []string{"/entity/register", "POST", "name", "username", "pwhash", "token", "seedpwd", "entityType"}, // POST
	8: []string{"/entity/contractor/dashboard", "GET"},                                                      // GET
}

// entityValidateHelper is a helper that helps validate an entity
func entityValidateHelper(w http.ResponseWriter, r *http.Request, options []string, method string) (core.Entity, error) {
	var prepEntity core.Entity

	err := checkReqdParams(w, r, options, method)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.NotEntityError)
		return prepEntity, errors.New("reqd params not present can't be empty")
	}

	var username, token string
	if method == "GET" {
		username, token = r.URL.Query()["username"][0], r.URL.Query()["token"][0]
	} else if method == "POST" {
		username, token = r.FormValue("username"), r.FormValue("token")
	} else {
		log.Println("method not recognized, quitting")
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return prepEntity, errors.New("invalid method, quitting")
	}

	prepEntity, err = core.ValidateEntity(username, token)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.NotEntityError)
		log.Println("did not validate investor", err)
		return prepEntity, err
	}

	return prepEntity, nil
}

// validateEntity is an endpoint that vlaidates is a specific entity is registered on the platform
func validateEntity() {
	http.HandleFunc(EntityRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, []string{}, EntityRPC[1][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}
		erpc.MarshalSend(w, prepEntity)
	})
}

// getStage0Contracts gets a list of all the pre origianted contracts on the platform
func getStage0Contracts() {
	http.HandleFunc(EntityRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, EntityRPC[2][2:], EntityRPC[2][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}

		x, err := core.RetrieveOriginatorProjects(core.Stage0.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving originator project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// getStage1Contracts gets a list of all the originated contracts on the platform
func getStage1Contracts() {
	http.HandleFunc(EntityRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, EntityRPC[3][2:], EntityRPC[3][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}

		x, err := core.RetrieveOriginatorProjects(core.Stage1.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving originator projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// getStage2Contracts gets a list of all the proposed contracts on the platform
func getStage2Contracts() {
	http.HandleFunc(EntityRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, EntityRPC[4][2:], EntityRPC[4][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}

		x, err := core.RetrieveContractorProjects(core.Stage2.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving contractor projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// addCollateral is a route that a contractor can use to add collateral
func addCollateral() {
	http.HandleFunc(EntityRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		prepEntity, err := entityValidateHelper(w, r, EntityRPC[5][2:], EntityRPC[5][1])
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, EntityRPC[5][2:], EntityRPC[5][1])
		if err != nil {
			log.Println(err)
			return
		}

		amountx := r.FormValue("amount")
		collateral := r.FormValue("collateral")

		amount, err := utils.ToFloat(amountx)
		if err != nil {
			log.Println("Error while converting string to float", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
			return
		}

		err = prepEntity.AddCollateral(amount, collateral)
		if err != nil {
			log.Println("Error while adding collateral", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// proposeOpensolarProject creates a contract which the contractor proposes towards a particular project
func proposeOpensolarProject() {
	http.HandleFunc(EntityRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepEntity, err := entityValidateHelper(w, r, EntityRPC[6][2:], EntityRPC[6][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}

		err = checkReqdParams(w, r, EntityRPC[6][2:], EntityRPC[6][1])
		if err != nil {
			log.Println(err)
			return
		}

		projIndexx := r.FormValue("projIndex")
		feex := r.FormValue("fee")

		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println("project idnex not int, quitting!")
			return
		}

		x, err := core.RetrieveProject(projIndex)
		if err != nil {
			log.Println("couldn't retrieve project with index")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		fee, err := utils.ToFloat(feex)
		if err != nil {
			log.Println("fee passed not integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
		}

		x.TotalValue += fee
		x.OriginatorFee = fee
		x.OriginatorIndex = prepEntity.U.Index
		x.Stage = 2

		err = x.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// registerEntity creates and stores a new entity on the platform
func registerEntity() {
	http.HandleFunc(EntityRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, EntityRPC[7][2:], EntityRPC[7][1])
		if err != nil {
			log.Println(err)
			return
		}

		name := r.FormValue("name")
		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")
		seedpwd := r.FormValue("seedpwd")
		entityType := r.FormValue("entityType")

		// check for username collision here. If the username already exists, fetch details from that and register as investor
		if core.CheckUsernameCollision(username) {
			// user already exists on the platform, need to retrieve the user
			user, err := userValidateHelper(w, r, nil, RecpRPC[2][1]) // check whether this person is a user and has params
			if err != nil {
				return
			}

			// this is the same user who wants to register as an investor now, check if encrypted seed decrypts
			seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			pubkey, err := wallet.ReturnPubkey(seed)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			if pubkey != user.StellarWallet.PublicKey {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
			}

			var a core.Entity
			switch entityType {
			case "developer":
				a.Developer = true
			case "contractor":
				a.Contractor = true
			case "guarantor":
				a.Guarantor = true
			case "originator":
				a.Originator = true
			}

			a.U = &user
			err = a.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, a)
			return
		}

		var user core.Entity
		switch entityType {
		case "developer":
			user, err = core.NewDeveloper(username, pwhash, seedpwd, name)
		case "contractor":
			user, err = core.NewContractor(username, pwhash, seedpwd, name)
		case "guarantor":
			user, err = core.NewGuarantor(username, pwhash, seedpwd, name)
		case "originator":
			user, err = core.NewOriginator(username, pwhash, seedpwd, name)

		}

		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}
