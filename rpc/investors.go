package rpc

import (
	"errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"

	tickers "github.com/Varunram/essentials/exchangetickers"
	consts "github.com/YaleOpenLab/opensolar/consts"
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
	invDashboard()
	setCompanyBool()
	setCompany()
}

// InvRPC contains a list of all investor related endpoints
var InvRPC = map[int][]string{
	1:  []string{"/investor/register", "POST", "name", "username", "pwhash", "token", "seedpwd"},      // POST
	2:  []string{"/investor/validate", "GET"},                                                         // GET
	3:  []string{"/investor/all", "GET"},                                                              // GET
	4:  []string{"/investor/invest", "POST", "seedpwd", "projIndex", "amount"},                        // POST
	5:  []string{"/investor/vote", "POST", "votes", "projIndex"},                                      // POST
	6:  []string{"/investor/localasset", "POST", "assetName"},                                         // POST
	7:  []string{"/investor/sendlocalasset", "POST", "assetName", "seedpwd", "destination", "amount"}, // POST
	8:  []string{"/investor/sendemail", "POST", "message", "to"},                                      // POST
	9:  []string{"/investor/dashboard", "GET"},                                                        // GET
	10: []string{"/investor/company/set", "POST"},                                                     // POST
	11: []string{"/investor/company/details", "POST", "companytype",
		"name", "legalname", "address", "country", "city", "zipcode", "role"}, // POST
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
		err := checkReqdParams(w, r, InvRPC[1][2:], InvRPC[1][1])
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
		_, err := InvValidateHelper(w, r, InvRPC[3][2:], InvRPC[3][1])
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

type invDHelper struct {
	Index            int     `json:"Index"`
	Stage            int     `json:"Stage"`
	Name             string  `json:"Project Name"`
	Location         string  `json:"Location"`
	Capacity         string  `json:"Capacity"`
	YourInvestment   float64 `json:"Your Investment"`
	YourReturn       string  `json:"Your Return"`
	InvestmentRating string  `json:"Investment Rating"`
	ImpactRating     string  `json:"Impact Rating"`
	ProjectActions   string  `json:"Project Actions"`
}

type invDashboardStruct struct {
	YourProfile struct {
		Name  string `json:"Name"`
		Roles string `json:"Roles"`
	} `json:"Your Profile"`
	YourInvestments struct {
		TotalInvestments float64 `json:"Total Investments"`
		ProjectsInvested int     `json:"Projects Invested"`
	} `json:"Your Investments"`
	YourReturns struct {
		NetReturns   string `json:"Net Returns"`
		RecsReceived string `json:"RECs Received"`
	} `json:"Your Returns"`
	EFacilitate struct {
		DirectContributions string `json:"My Direct Contributions"`
		TotalContributions  string `json:"Total Contributions"`
	} `json:"Energy You Facilitate"`
	PrimaryAddress   string       `json:"Main Wallet"`
	SecondaryAddress string       `json:"Secondary Wallet"`
	AccountBalance1  float64      `json:"Account Balance 1"`
	AccountBalance2  float64      `json:"Account Balance 2"`
	NetBalance       float64      `json:"Balance"`
	InvestedProjects []invDHelper `json:"Your Invested Projects"`
}

// invDashboard returns the parameters needed for displaying details on the frontend
func invDashboard() {
	http.HandleFunc(InvRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r, InvRPC[9][2:], InvRPC[9][1])
		if err != nil {
			return
		}

		var ret invDashboardStruct
		for _, index := range prepInvestor.InvestedSolarProjectsIndices {
			project, err := core.RetrieveProject(index)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			var temp invDHelper
			temp.Stage = project.Stage
			temp.Name = project.Name
			temp.Location = project.City + ", " + project.State + ", " + project.Country
			temp.Capacity = project.Content.Details["Other Details"]["panel size"].(string)
			temp.YourInvestment = project.InvestorMap[prepInvestor.U.StellarWallet.PublicKey] * project.TotalValue
			temp.YourReturn = "Donation"
			temp.InvestmentRating = "N/A"
			temp.ImpactRating = "4/4"
			temp.ProjectActions = "No immediate action"
			temp.Index = project.Index

			ret.InvestedProjects = append(ret.InvestedProjects, temp)
		}

		inv, err := core.SearchForInvestor(prepInvestor.U.Name)
		if err == nil && inv.U.Name != "" {
			ret.YourProfile.Roles += " Investor"
		}

		recp, err := core.SearchForRecipient(prepInvestor.U.Name)
		if err == nil && recp.U.Name != "" {
			ret.YourProfile.Roles += " Recipient"
		}

		entity, err := core.SearchForEntity(prepInvestor.U.Name)
		if err == nil && entity.U.Name != "" {
			ret.YourProfile.Roles += " Entity"
		}

		ret.YourProfile.Name = inv.U.Name
		ret.YourInvestments.TotalInvestments = prepInvestor.AmountInvested
		ret.YourInvestments.ProjectsInvested = len(prepInvestor.InvestedSolarProjects)
		ret.YourReturns.NetReturns = "$0"
		ret.YourReturns.RecsReceived = "10 MWh"
		ret.EFacilitate.DirectContributions = "1000 KWh"
		ret.EFacilitate.TotalContributions = "1000 KWh"

		ret.PrimaryAddress = prepInvestor.U.StellarWallet.PublicKey
		ret.SecondaryAddress = prepInvestor.U.SecondaryWallet.PublicKey

		if !consts.Mainnet {
			primNativeBalance := xlm.GetNativeBalance(prepInvestor.U.StellarWallet.PublicKey) * 10000000.0
			if primNativeBalance < 0 {
				primNativeBalance = 0
			}

			secNativeBalance := xlm.GetNativeBalance(prepInvestor.U.SecondaryWallet.PublicKey) * 10000000.0
			if secNativeBalance < 0 {
				secNativeBalance = 0
			}

			primUsdBalance := xlm.GetAssetBalance(prepInvestor.U.StellarWallet.PublicKey, consts.StablecoinCode)
			if primUsdBalance < 0 {
				primUsdBalance = 0
			}

			secUsdBalance := xlm.GetAssetBalance(prepInvestor.U.SecondaryWallet.PublicKey, consts.StablecoinCode)
			if secUsdBalance < 0 {
				secUsdBalance = 0
			}

			ret.AccountBalance1 = primNativeBalance + primUsdBalance
			ret.AccountBalance2 = secNativeBalance + secUsdBalance
		} else {
			xlmUSD, err := tickers.BinanceTicker()
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			}

			primNativeBalance := xlm.GetNativeBalance(prepInvestor.U.StellarWallet.PublicKey) * xlmUSD
			if primNativeBalance < 0 {
				primNativeBalance = 0
			}

			secNativeBalance := xlm.GetNativeBalance(prepInvestor.U.SecondaryWallet.PublicKey) * xlmUSD
			if secNativeBalance < 0 {
				secNativeBalance = 0
			}

			primUsdBalance := xlm.GetAssetBalance(prepInvestor.U.StellarWallet.PublicKey, consts.AnchorUSDCode)
			if primUsdBalance < 0 {
				primUsdBalance = 0
			}

			secUsdBalance := xlm.GetAssetBalance(prepInvestor.U.SecondaryWallet.PublicKey, consts.AnchorUSDCode)
			if secUsdBalance < 0 {
				secUsdBalance = 0
			}

			ret.AccountBalance1 = primNativeBalance + primUsdBalance
			ret.AccountBalance2 = secNativeBalance + secUsdBalance
		}

		if ret.AccountBalance2 < 0 {
			ret.AccountBalance2 = 0
		}

		if ret.AccountBalance1 < 0 {
			ret.AccountBalance1 = 0
		}

		ret.NetBalance = ret.AccountBalance1 + ret.AccountBalance2
		erpc.MarshalSend(w, ret)
	})
}

func setCompanyBool() {
	http.HandleFunc(InvRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r, InvRPC[10][2:], InvRPC[10][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		err = prepInvestor.SetCompany()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func setCompany() {
	http.HandleFunc(InvRPC[11][0], func(w http.ResponseWriter, r *http.Request) {
		prepInvestor, err := InvValidateHelper(w, r, InvRPC[11][2:], InvRPC[11][1])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		companyType := r.FormValue("companytype")
		switch companyType {
		case "For-Profit":
			log.Println("company type: For-Profit")
		case "Social Enterprise":
			log.Println("company type: Social Enterprise")
		case "Non Governmental":
			log.Println("company type: Non Governmental")
		case "Cooperative":
			log.Println("company type: Cooperative")
		case "Other":
			log.Println("company type: Other")
		default:
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		legalName := r.FormValue("legalname")
		address := r.FormValue("address")
		country := r.FormValue("country")
		city := r.FormValue("city")
		zipCode := r.FormValue("zipcode")
		role := r.FormValue("role")
		switch role {
		case "ceo":
			log.Println("role: ceo")
		case "employee":
			log.Println("role: employee")
		case "other":
			log.Println("role: other")
		default:
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// these are params which are not necessary
		var adminEmail, phoneNumber, taxIDNumber string

		if lenParseCheck(adminEmail) != nil || lenParseCheck(phoneNumber) != nil ||
			lenParseCheck(taxIDNumber) != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if r.FormValue("adminemail") != "" {
			adminEmail = r.FormValue("adminemail")
		}
		if r.FormValue("phonenumber") != "" {
			phoneNumber = r.FormValue("phoneNumber")
		}
		if r.FormValue("taxidnumber") != "" {
			taxIDNumber = r.FormValue("taxidnumber")
		}

		err = prepInvestor.SetCompanyDetails(companyType, name, legalName, adminEmail, phoneNumber, address,
			country, city, zipCode, taxIDNumber, role)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
