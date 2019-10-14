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

	core "github.com/YaleOpenLab/opensolar/core"
	notif "github.com/YaleOpenLab/opensolar/notif"
)

// setupInvestorRPCs sets up all investor related RPCs
func setupInvestorRPCs() {
	registerInvestor()
	validateInvestor()
	getAllInvestors()
	invest()
	voteTowardsProject()
	addLocalAssetInv()
	invAssetInv()
	sendEmail()
}

// InvRPC contains a list of all investor related endpoints
var InvRPC = map[int][]string{
	1: []string{"/investor/register", "POST", "name", "username", "pwhash", "token", "seedpwd"},      // POST
	2: []string{"/investor/validate", "GET"},                                                         // GET
	3: []string{"/investor/all", "GET"},                                                              // GET
	4: []string{"/investor/invest", "POST", "seedpwd", "projIndex", "amount"},                        // POST
	5: []string{"/investor/vote", "POST", "votes", "projIndex"},                                      // POST
	6: []string{"/investor/localasset", "POST", "assetName"},                                         // POST
	7: []string{"/investor/sendlocalasset", "POST", "assetName", "seedpwd", "destination", "amount"}, // POST
	8: []string{"/investor/sendemail", "POST", "message", "to"},                                      // POST
}

// InvValidateHelper is a helper used to validate an investor on the platform
func InvValidateHelper(w http.ResponseWriter, r *http.Request, options []string, method string) (core.Investor, error) {
	var prepInvestor core.Investor
	var err error

	err = checkReqdParams(w, r, options, method)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return prepInvestor, errors.New("reqd params not present can't be empty")
	}

	var username, token string
	if method == "GET" {
		username, token = r.URL.Query()["username"][0], r.URL.Query()["token"][0]
	} else {
		username, token = r.FormValue("username"), r.FormValue("token")
	}

	prepInvestor, err = core.ValidateInvestor(username, token)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		log.Println("did not validate investor", err)
		return prepInvestor, err
	}

	return prepInvestor, nil
}

func registerInvestor() {
	http.HandleFunc(InvRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, InvRPC[1][2:], InvRPC[1][1])
		if err != nil {
			log.Println(err)
			return
		}

		name := r.FormValue("name")
		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")
		token := r.FormValue("token")
		seedpwd := r.FormValue("seedpwd")

		// check for username collision here. If the username already exists, fetch details from that and register as investor
		if core.CheckUsernameCollision(username) {
			// user already exists on the platform, need to retrieve the user
			user, err := core.ValidateUser(username, token) // check whether this person is a user and has params
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
			var a core.Investor
			a.U = &user
			err = a.Save()
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, a)
			return
		}

		user, err := core.NewInvestor(username, pwhash, seedpwd, name)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// validateInvestor validates the username and pwhash of a given investor
func validateInvestor() {
	http.HandleFunc(InvRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		prepInvestor, err := InvValidateHelper(w, r, InvRPC[2][2:], InvRPC[2][1])
		if err != nil {
			return
		}
		erpc.MarshalSend(w, prepInvestor)
	})
}

// getAllInvestors gets a list of all the investors in the database
func getAllInvestors() {
	http.HandleFunc(InvRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = InvValidateHelper(w, r, InvRPC[3][2:], InvRPC[3][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		investors, err := core.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, investors)
	})
}

// Invest invests in a project of the investor's choice
func invest() {
	http.HandleFunc(InvRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		investor, err := InvValidateHelper(w, r, InvRPC[4][2:], InvRPC[4][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		seedpwd := r.FormValue("seedpwd")
		projIndexx := r.FormValue("projIndex")
		amountx := r.FormValue("amount")

		investorSeed, err := wallet.DecryptSeed(investor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		amount, err := utils.ToFloat(amountx)
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

		log.Println("reaches here", investorPubkey)
		if !xlm.AccountExists(investorPubkey) {
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}

		err = core.Invest(projIndex, investor.U.Index, amount, investorSeed)
		if err != nil {
			log.Println("did not invest in order", err)
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// voteTowardsProject votes towards a proposed project of the user's choice.
func voteTowardsProject() {
	http.HandleFunc(InvRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		investor, err := InvValidateHelper(w, r, InvRPC[5][2:], InvRPC[5][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		votesx := r.FormValue("votes")
		projIndexx := r.FormValue("projIndex")

		votes, err := utils.ToFloat(votesx)
		if err != nil {
			log.Println("votes not float, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println("error while converting project index to int: ", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = core.VoteTowardsProposedProject(investor.U.Index, votes, projIndex)
		if err != nil {
			log.Println("did not vote towards proposed project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// addLocalAssetInv adds a local asset that can be traded in a p2p fashion wihtout direct involvement
// from the platform
func addLocalAssetInv() {
	http.HandleFunc(InvRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepInvestor, err := InvValidateHelper(w, r, InvRPC[6][2:], InvRPC[6][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.FormValue("assetName")

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
	http.HandleFunc(InvRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepInvestor, err := InvValidateHelper(w, r, InvRPC[7][2:], InvRPC[7][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		assetName := r.FormValue("assetName")
		seedpwd := r.FormValue("seedpwd")
		destination := r.FormValue("desination")
		amountx := r.FormValue("amount")

		seed, err := wallet.DecryptSeed(prepInvestor.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println("did not decrypt seed", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		amount, err := utils.ToFloat(amountx)
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
	http.HandleFunc(InvRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepInvestor, err := InvValidateHelper(w, r, InvRPC[8][2:], InvRPC[8][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		message := r.FormValue("message")
		to := r.FormValue("to")

		err = notif.SendEmail(message, to, prepInvestor.U.Name)
		if err != nil {
			log.Println("did not send email", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
