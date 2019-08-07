package rpc

import (
	"errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	// opzones "github.com/YaleOpenLab/openx/platforms/ozones"
)

// setupInvestorRPCs sets up all RPCs related to the investor
func setupInvestorRPCs() {
	registerInvestor()
	validateInvestor()
	getAllInvestors()
	invest()
	voteTowardsProject()
	addLocalAssetInv()
	invAssetInv()
	sendEmail()
	// investInConstructionBond()
	// investInLivingUnitCoop()
}

// InvValidateHelper is a helper that is used to validate an ivnestor on the platform
func InvValidateHelper(w http.ResponseWriter, r *http.Request, options ...string) (opensolar.Investor, error) {
	var prepInvestor opensolar.Investor
	var err error
	// need to pass the pwhash param here
	if r.URL.Query() == nil {
		return prepInvestor, errors.New("url query can't be empty")
	}

	options = append(options, "username", "pwhash")

	for _, option := range options {
		if r.URL.Query()[option] == nil {
			return prepInvestor, errors.New("required param: " + option + "not specified, quitting")
		}
	}

	if len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepInvestor, errors.New("pwhash length not 128, quitting")
	}

	prepInvestor, err = opensolar.ValidateInvestor(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		log.Println("did not validate investor", err)
		return prepInvestor, err
	}

	return prepInvestor, nil
}

func registerInvestor() {
	http.HandleFunc("/investor/register", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		// to register, we need the name, username and pwhash
		if r.URL.Query()["name"] == nil || r.URL.Query()["username"] == nil ||
			r.URL.Query()["pwhash"] == nil || r.URL.Query()["seedpwd"] == nil {
			log.Println("missing basic set of params that can be used ot validate a user")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		name := r.URL.Query()["name"][0]
		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		seedpwd := r.URL.Query()["seedpwd"][0]

		// check for username collision here. If the username already exists, fetch details from that and register as investor
		_, err := database.CheckUsernameCollision(username)
		if err != nil {
			// user already exists on the platform, need to retrieve the user
			user, err := CheckReqdParams(w, r) // check whether this person is a user and has params
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
			}
			// username collision, check other fields by fetching user details for the collided user
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
			var a opensolar.Investor
			a.U = &user
			err = a.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, a)
			return
		}

		user, err := opensolar.NewInvestor(username, pwhash, seedpwd, name)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// validateInvestor retrieves the investor after valdiating if such an ivnestor exists
// by checking the pwhash of the given investor with the stored one
func validateInvestor() {
	http.HandleFunc("/investor/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepInvestor, err := InvValidateHelper(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, prepInvestor)
	})
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestors() {
	http.HandleFunc("/investor/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		investors, err := opensolar.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, investors)
	})
}

// Invest invests in a specific project of the user's choice
func invest() {
	http.HandleFunc("/investor/invest", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		// need the following params to invest in a project:
		// 1. Seed pwhash (for the investor)
		// 2. project index
		// 3. investment amount
		// 4. Login username (for the investor)
		// 5. Login pwhash (for the investor)
		investor, err := InvValidateHelper(w, r, "seedpwd", "projIndex", "amount")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd := r.URL.Query()["seedpwd"][0]
		investorSeed, err := wallet.DecryptSeed(investor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		investorPubkey, err := wallet.ReturnPubkey(investorSeed)
		if err != nil {
			log.Println("did not return pubkey", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if !xlm.AccountExists(investorPubkey) {
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}

		err = opensolar.Invest(projIndex, investor.U.Index, amount, investorSeed)
		if err != nil {
			log.Println("did not invest in order", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// voteTowardsProject votes towards a specific propsoed project of the user's choice.
func voteTowardsProject() {
	http.HandleFunc("/investor/vote", func(w http.ResponseWriter, r *http.Request) {
		investor, err := InvValidateHelper(w, r, "votes", "projIndex")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		votes, err := utils.ToFloat(r.URL.Query()["votes"][0])
		if err != nil {
			log.Println("votes not float, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = opensolar.VoteTowardsProposedProject(investor.U.Index, votes, projIndex)
		if err != nil {
			log.Println("did not vote towards proposed project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// addLocalAssetInv adds a local asset that can be traded in a p2p fashion wihtout direct involvement
// from the platform. The platform can have a UI that will deal with this or this can be
// made an emualtor only function so that only experienced users use this.
func addLocalAssetInv() {
	http.HandleFunc("/investor/localasset", func(w http.ResponseWriter, r *http.Request) {

		prepInvestor, err := InvValidateHelper(w, r, "assetName")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		prepInvestor.U.LocalAssets = append(prepInvestor.U.LocalAssets, assetName)
		err = prepInvestor.Save()
		if err != nil {
			log.Println("did not save investor", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// invAssetInv sends a local asset to a remote peer
func invAssetInv() {
	http.HandleFunc("/investor/sendlocalasset", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r, "assetName", "seedpwd", "destination", "amount")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.URL.Query()["assetName"][0]

		seedpwd := r.URL.Query()["seedpwd"][0]
		seed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		destination := r.URL.Query()["destination"][0]
		amount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		found := true
		for _, elem := range prepInvestor.U.LocalAssets {
			if elem == assetName {
				found = true
			}
		}

		if !found {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		_, txhash, err := assets.SendAssetFromIssuer(assetName, destination, amount, seed, prepInvestor.U.StellarWallet.PublicKey)
		if err != nil {
			log.Println("did not send asset from issuer", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, txhash)
	})
}

// sendEmail sends an email to a specific entity
func sendEmail() {
	http.HandleFunc("/investor/sendemail", func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r, "message", "to")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		message := r.URL.Query()["message"][0]
		to := r.URL.Query()["to"][0]
		err = notif.SendEmail(message, to, prepInvestor.U.Name)
		if err != nil {
			log.Println("did not send email", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

/*
// investInConstructionBond invests a specific amount in a bond of the user's choice
func investInConstructionBond() {
	http.HandleFunc("/constructionbond/invest", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		var err error

		prepInvestor, err := InvValidateHelper(w, r, "amount", "projIndex", "seedpwd")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		invAmount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		seedpwd := r.URL.Query()["seedpwd"][0]

		invSeed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not get investor seed from password", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = opzones.InvestInConstructionBond(projIndex, prepInvestor.U.Index, invAmount, invSeed)
		if err != nil {
			log.Println("did not invest in bond", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// InvestInCoop invests in a coop of the user's choice
func investInLivingUnitCoop() {
	http.HandleFunc("/livingunitcoop/invest", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		var err error

		prepInvestor, err := InvValidateHelper(w, r, "amount", "projIndex", "seedpwd")
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		invAmount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		seedpwd := r.URL.Query()["seedpwd"][0]

		invSeed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not get investor seed from password", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		recpSeed := "SA5LO2G3XR37YY7566K2NHWQCK6PFXMF7UE64WGFBCOAPFHEKNSWT6PE"
		err = opzones.InvestInLivingUnitCoop(projIndex, prepInvestor.U.Index, invAmount, invSeed, recpSeed)
		if err != nil {
			log.Println("did not invest in the coop", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
*/
