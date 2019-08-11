package rpc

import (
	"errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	openxrpc "github.com/YaleOpenLab/openx/rpc"

	core "github.com/YaleOpenLab/opensolar/core"
)

// setupRecipientRPCs sets up all RPCs related to the recipient
func setupRecipientRPCs() {
	registerRecipient()
	validateRecipient()
	getAllRecipients()
	payback()
	storeDeviceId()
	storeStartTime()
	storeDeviceLocation()
	chooseBlindAuction()
	chooseVickreyAuction()
	chooseTimeAuction()
	unlockOpenSolar()
	addEmail()
	finalizeProject()
	originateProject()
	calculateTrustLimit()
	// unlockCBond()
	storeStateHash()
}

// RecpValidateHelper is a helper that helps validates recipients in routes
func RecpValidateHelper(w http.ResponseWriter, r *http.Request, options ...string) (core.Recipient, error) {
	var prepRecipient core.Recipient
	var err error
	if r.URL.Query() == nil {
		return prepRecipient, errors.New("url query can't be empty")
	}

	options = append(options, "username", "pwhash")

	for _, option := range options {
		if r.URL.Query()[option] == nil {
			return prepRecipient, errors.New("required param: " + option + "not specified, quitting")
		}
	}

	if len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepRecipient, errors.New("pwhash length not 128, quitting")
	}

	prepRecipient, err = core.ValidateRecipient(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		log.Println("did not validate recipient", err)
		return prepRecipient, err
	}

	return prepRecipient, nil
}

// getAllRecipients gets a list of all the recipients who have registered on the platform
func getAllRecipients() {
	http.HandleFunc("/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipients, err := core.RetrieveAllRecipients()
		if err != nil {
			log.Println("did not retrieve all recipients", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, recipients)
	})
}

// registerRecipient creates and stores a new recipient on the platform
func registerRecipient() {
	http.HandleFunc("/recipient/register", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["name"] == nil || r.URL.Query()["username"] == nil || r.URL.Query()["pwd"] == nil || r.URL.Query()["seedpwd"] == nil {
			log.Println("missing basic set of params that can be used ot validate a user")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		name := r.URL.Query()["name"][0]
		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]

		// check for username collision here. If the username already exists, fetch details from that and register as investor
		if core.CheckUsernameCollision(username) {
			// user already exists on the platform, need to retrieve the user
			user, err := openxrpc.CheckReqdParams(w, r) // check whether this person is a user and has params
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
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
			var a core.Recipient
			a.U = &user
			err = a.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, a)
			return
		}

		user, err := core.NewRecipient(username, pwhash, seedpwd, name)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// validateRecipient validates a recipient on the platform
func validateRecipient() {
	http.HandleFunc("/recipient/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepRecipient, err := RecpValidateHelper(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, prepRecipient)
	})
}

// payback pays back towards an  invested order
func payback() {
	http.HandleFunc("/recipient/payback", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepRecipient, err := RecpValidateHelper(w, r, "assetName", "amount", "seedpwd", "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		recpIndex := prepRecipient.U.Index
		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		assetName := r.URL.Query()["assetName"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]
		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		recipientSeed, err := wallet.DecryptSeed(prepRecipient.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		log.Println(recpIndex, projIndex, assetName, amount, recipientSeed)
		err = core.Payback(recpIndex, projIndex, assetName, amount, recipientSeed)
		if err != nil {
			log.Println("did not payback", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeDeviceId stores the recipient's device id from the teller. Called by the teller
func storeDeviceId() {
	http.HandleFunc("/recipient/deviceId", func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		prepRecipient, err := RecpValidateHelper(w, r, "deviceId")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		// we have the recipient ready. Now set the device id
		prepRecipient.DeviceId = r.URL.Query()["deviceId"][0]
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeStartTime stores the start time of the remote device installed as part of an
// invested project. Called by the teller
func storeStartTime() {
	http.HandleFunc("/recipient/startdevice", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepRecipient, err := RecpValidateHelper(w, r, "start")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		prepRecipient.DeviceStarts = append(prepRecipient.DeviceStarts, r.URL.Query()["start"][0])
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeDeviceLocation stores the location of the remote device when it starts up. Called by the teller
func storeDeviceLocation() {
	http.HandleFunc("/recipient/storelocation", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepRecipient, err := RecpValidateHelper(w, r, "location")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		prepRecipient.DeviceLocation = r.URL.Query()["location"][0]
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseBlindAuction chooses a blind auction method to choose for the winner. Also commonly
// known as a 1st price auction.
func chooseBlindAuction() {
	http.HandleFunc("/recipient/auction/choose/blind", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not validate recipient projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bestContract, err := core.SelectContractBlind(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4)
		if err != nil {
			log.Println("did not set final project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseVickreyAuction chooses a vickrey auction method to choose the winning contractor.
// also known as a second price auction
func chooseVickreyAuction() {
	http.HandleFunc("/recipient/auction/choose/vickrey", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not retrieve recipient projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bestContract, err := core.SelectContractVickrey(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4)
		if err != nil {
			log.Println("did not set final project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseTimeAuction chooses the winning contractor based on least completion time
func chooseTimeAuction() {
	http.HandleFunc("/recipient/auction/choose/time", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r)
		if err != nil {
			log.Println("did not validate recipient", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if err != nil {
			log.Println("did not retrieve recipient projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bestContract, err := core.SelectContractTime(allContracts)
		if err != nil {
			log.Println("did not select contract", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = bestContract.SetStage(4)
		if err != nil {
			log.Println("did not set final project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// unlockOpenSolar unlocks a project which has just been invested in, signalling that the recipient
// has accepted the investment.
func unlockOpenSolar() {
	http.HandleFunc("/recipient/unlock/opensolar", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r, "seedpwd", "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("did not parse to integer", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = core.UnlockProject(recipient.U.Username, recipient.U.Pwhash, projIndex, seedpwd)
		if err != nil {
			log.Println("did not unlock project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// addEmail adds an email address to the recipient's profile
func addEmail() {
	http.HandleFunc("/recipient/addemail", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r, "email")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		email := r.URL.Query()["email"][0]
		err = recipient.U.AddEmail(email)
		if err != nil {
			log.Println("did not add email", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// finalizeProject finalizes (ie moves from stage 2 to 3) a specific project
func finalizeProject() {
	http.HandleFunc("/recipient/finalize", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		_, err := RecpValidateHelper(w, r, "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("did not parse to integer", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		project, err := core.RetrieveProject(projIndex)
		if err != nil {
			log.Println("did not retrieve project", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = project.SetStage(4)
		if err != nil {
			log.Println("did not set final project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// originateProject originates (ie moves from stage 0 to 1) a project
func originateProject() {
	http.HandleFunc("/recipient/originate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r, "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("did not parse to integer", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = core.RecipientAuthorize(projIndex, recipient.U.Index)
		if err != nil {
			log.Println("did not authorize project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// calculateTrustLimit calculates the trust limit associated with a specific asset.
func calculateTrustLimit() {
	http.HandleFunc("/recipient/trustlimit", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipient, err := RecpValidateHelper(w, r, "assetName")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		trustLimit, err := xlm.GetAssetTrustLimit(recipient.U.StellarWallet.PublicKey, assetName)
		if err != nil {
			log.Println("did not get trust limit", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, trustLimit)
	})
}

// storeStateHash stores the start time of the remote device installed as part of an invested project.
// Called by the teller
func storeStateHash() {
	http.HandleFunc("/recipient/ssh", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := RecpValidateHelper(w, r, "hash")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		prepRecipient.StateHashes = append(prepRecipient.StateHashes, r.URL.Query()["hash"][0])
		err = prepRecipient.Save()
		if err != nil {
			log.Println("did not save recipient", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
