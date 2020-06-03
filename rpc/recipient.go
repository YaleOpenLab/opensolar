package rpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/YaleOpenLab/opensolar/messages"
	"github.com/YaleOpenLab/opensolar/oracle"

	tickers "github.com/Varunram/essentials/exchangetickers"
	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	wallet "github.com/Varunram/essentials/xlm/wallet"

	"github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
)

// setupRecipientRPCs sets up all RPCs related to the recipient
func setupRecipientRPCs() {
	registerRecipient()
	validateRecipient()
	getAllRecipients()
	payback()
	storeDeviceID()
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
	setOneTimeUnlock()
	storeTellerURL()
	storeTellerDetails()
	recpDashboard()
	storeTellerEnergy()
	setCompanyBoolRecp()
	setCompanyRecp()
}

// RecpRPC is a collection of all recipient RPC endpoints and their required params
var RecpRPC = map[int][]string{
	1:  {"/recipient/all", "GET"},                                                                                                   // GET
	2:  {"/recipient/register", "POST", "name", "username", "pwhash", "seedpwd"},                                                    // POST
	3:  {"/recipient/validate", "GET"},                                                                                              // GET
	4:  {"/recipient/payback", "POST", "assetName", "amount", "seedpwd", "projIndex"},                                               // POST
	5:  {"/recipient/deviceId", "POST", "deviceId"},                                                                                 // POST
	6:  {"/recipient/startdevice", "POST", "start"},                                                                                 // POST
	7:  {"/recipient/storelocation", "POST", "location"},                                                                            // POST
	8:  {"/recipient/auction/choose/blind", "GET"},                                                                                  // GET
	9:  {"/recipient/auction/choose/vickrey", "GET"},                                                                                // GET
	10: {"/recipient/auction/choose/time", "GET"},                                                                                   // GET
	11: {"/recipient/unlock/opensolar", "POST", "seedpwd", "projIndex"},                                                             // POST
	12: {"/recipient/addemail", "POST", "email"},                                                                                    // POST
	13: {"/recipient/finalize", "POST", "projIndex"},                                                                                // POST
	14: {"/recipient/originate", "POST", "projIndex"},                                                                               // POST
	15: {"/recipient/trustlimit", "GET", "assetName"},                                                                               // GET
	16: {"/recipient/ssh", "POST", "hash"},                                                                                          // POST
	17: {"/recipient/onetimeunlock", "POST", "projIndex", "seedpwd"},                                                                // POST
	18: {"/recipient/register/teller", "POST", "url", "projIndex"},                                                                  // POST
	19: {"/recipient/teller/details", "POST", "projIndex", "url", "brokerurl", "topic"},                                             // POST
	20: {"/recipient/dashboard", "GET"},                                                                                             // GET
	21: {"/recipient/company/set", "POST"},                                                                                          // POST
	22: {"/recipient/company/details", "POST", "companytype", "name", "legalname", "address", "country", "city", "zipcode", "role"}, // POST
	23: {"/recipient/teller/energy", "POST", "energy"},                                                                              // POST
}

// recpValidateHelper is a helper that helps validates recipients in routes
func recpValidateHelper(w http.ResponseWriter, r *http.Request, options []string, method string) (core.Recipient, error) {
	var prepRecipient core.Recipient
	var err error

	err = checkReqdParams(w, r, options, method)
	if err != nil {
		log.Println(err)
		return prepRecipient, errors.New("reqd params not present can't be empty")
	}

	var username, token string

	if r.Method == "GET" {
		username, token = r.URL.Query()["username"][0], r.URL.Query()["token"][0]
	} else if r.Method == "POST" {
		username, token = r.FormValue("username"), r.FormValue("token")
	}

	prepRecipient, err = core.ValidateRecipient(username, token)
	if erpc.Err(w, err, erpc.StatusUnauthorized, "did not validate recipient", messages.NotRecipientError) {
		return prepRecipient, err
	}

	return prepRecipient, nil
}

// getAllRecipients gets a list of all the recipients who have registered
func getAllRecipients() {
	http.HandleFunc(RecpRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := recpValidateHelper(w, r, RecpRPC[1][2:], RecpRPC[1][1])
		if err != nil {
			return
		}
		recipients, err := core.RetrieveAllRecipients()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve all recipients") {
			return
		}
		erpc.MarshalSend(w, recipients)
	})
}

// registerRecipient creates and stores a new recipient
func registerRecipient() {
	http.HandleFunc(RecpRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := checkReqdParams(w, r, RecpRPC[2][2:], RecpRPC[2][1])
		if err != nil {
			log.Println(err)
			return
		}

		name := r.FormValue("name")
		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")
		seedpwd := r.FormValue("seedpwd")

		// check for username collision here. If the username already exists, fetch details from that and register as investor
		if core.CheckUsernameCollision(username) {
			// user already exists, need to retrieve the user
			user, err := userValidateHelper(w, r, nil, RecpRPC[2][1]) // check whether this person is a user and has params
			if err != nil {
				return
			}

			// this is the same user who wants to register as an investor now, check if encrypted seed decrypts
			seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
			pubkey, err := wallet.ReturnPubkey(seed)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
			if pubkey != user.StellarWallet.PublicKey {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized)
				return
			}
			var a core.Recipient
			a.U = &user
			err = a.Save()
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}
			erpc.MarshalSend(w, a)
			return
		}

		user, err := core.NewRecipient(username, pwhash, seedpwd, name)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// validateRecipient validates a recipient
func validateRecipient() {
	http.HandleFunc(RecpRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[3][2:], RecpRPC[3][1])
		if err != nil {
			return
		}
		erpc.MarshalSend(w, prepRecipient)
	})
}

// payback pays back towards an  invested order
func payback() {
	http.HandleFunc(RecpRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[4][2:], RecpRPC[4][1])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")
		assetName := r.FormValue("assetName")
		seedpwd := r.FormValue("seedpwd")
		amountx := r.FormValue("amount")

		recpIndex := prepRecipient.U.Index
		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}
		amount, err := utils.ToFloat(amountx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		recipientSeed, err := wallet.DecryptSeed(prepRecipient.U.StellarWallet.EncryptedSeed, seedpwd)
		if erpc.Err(w, err, erpc.StatusBadRequest, "did not decrypt seed") {
			return
		}

		err = core.Payback(recpIndex, projIndex, assetName, amount, recipientSeed)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not payback") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeDeviceID stores the recipient's device id from the teller. Called by the teller
func storeDeviceID() {
	http.HandleFunc(RecpRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[5][2:], RecpRPC[5][1])
		if err != nil {
			return
		}

		deviceID := r.FormValue("deviceId")
		// we have the recipient ready. Now set the device id
		prepRecipient.DeviceID = deviceID
		err = prepRecipient.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not save recipient") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeStartTime stores the start time of the remote device installed as part of an
// invested project. Called by the teller
func storeStartTime() {
	http.HandleFunc(RecpRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[6][2:], RecpRPC[6][1])
		if err != nil {
			return
		}

		start := r.FormValue("start")

		prepRecipient.DeviceStarts = append(prepRecipient.DeviceStarts, start)
		err = prepRecipient.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not save recipient") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// storeDeviceLocation stores the location of the remote device when it starts up. Called by the teller
func storeDeviceLocation() {
	http.HandleFunc(RecpRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[7][2:], RecpRPC[7][1])
		if err != nil {
			log.Println(err)
			return
		}

		location := r.FormValue("location")

		prepRecipient.DeviceLocation = location
		err = prepRecipient.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not save recipient") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseBlindAuction chooses a blind auction method to choose for the winner. Also commonly
// known as a 1st price auction.
func chooseBlindAuction() {
	http.HandleFunc(RecpRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[8][2:], RecpRPC[8][1])
		if err != nil {
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve recipient projects") {
			return
		}

		bestContract, err := core.SelectContractBlind(allContracts)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not select contract") {
			return
		}

		err = bestContract.SetStage(4)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not set final project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseVickreyAuction chooses a vickrey auction method to choose the winning contractor.
// also known as a second price auction
func chooseVickreyAuction() {
	http.HandleFunc(RecpRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[9][2:], RecpRPC[9][1])
		if err != nil {
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve recipient projects") {
			return
		}

		bestContract, err := core.SelectContractBlind(allContracts)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not select contract") {
			return
		}

		err = bestContract.SetStage(4)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not set final project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// chooseTimeAuction chooses the winning contractor based on least completion time
func chooseTimeAuction() {
	http.HandleFunc(RecpRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[10][2:], RecpRPC[10][1])
		if err != nil {
			return
		}

		allContracts, err := core.RetrieveRecipientProjects(core.Stage2.Number, recipient.U.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve recipient projects") {
			return
		}

		bestContract, err := core.SelectContractBlind(allContracts)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not select contract") {
			return
		}

		err = bestContract.SetStage(4)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not set final project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// unlockOpenSolar unlocks a project which has just been invested in, signalling that the recipient
// has accepted the investment.
func unlockOpenSolar() {
	http.HandleFunc(RecpRPC[11][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[11][2:], RecpRPC[11][1])
		if err != nil {
			return
		}

		seedpwd := r.FormValue("seedpwd")
		projIndexx := r.FormValue("projIndex")
		token := r.FormValue("token")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		err = core.UnlockProject(recipient.U.Username, token, projIndex, seedpwd)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not unlock project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// addEmail adds an email address to the recipient's profile
func addEmail() {
	http.HandleFunc(RecpRPC[12][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[12][2:], RecpRPC[12][1])
		if err != nil {
			return
		}

		email := r.FormValue("email")

		err = recipient.U.AddEmail(email)
		if erpc.Err(w, err, erpc.StatusBadRequest, "did not add email") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// finalizeProject finalizes (ie moves from stage 2 to 3) a project
func finalizeProject() {
	http.HandleFunc(RecpRPC[13][0], func(w http.ResponseWriter, r *http.Request) {
		_, err := recpValidateHelper(w, r, RecpRPC[13][2:], RecpRPC[13][1])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(projIndex)
		if erpc.Err(w, err, erpc.StatusBadRequest, "did not retrieve project") {
			return
		}

		err = project.SetStage(4)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not set final project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// originateProject originates (ie moves from stage 0 to 1) a project
func originateProject() {
	http.HandleFunc(RecpRPC[14][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[14][2:], RecpRPC[14][1])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		err = core.RecipientAuthorize(projIndex, recipient.U.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not authorize project") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// calculateTrustLimit calculates the trust limit associated with a asset.
func calculateTrustLimit() {
	http.HandleFunc(RecpRPC[15][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[15][2:], RecpRPC[15][1])
		if err != nil {
			return
		}

		assetName := r.URL.Query()["assetName"][0]
		trustLimit := xlm.GetAssetTrustLimit(recipient.U.StellarWallet.PublicKey, assetName)
		erpc.MarshalSend(w, trustLimit)
	})
}

// storeStateHash stores the state hashes of the teller
func storeStateHash() {
	http.HandleFunc(RecpRPC[16][0], func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[16][2:], RecpRPC[16][1])
		if err != nil {
			return
		}

		hash := r.FormValue("hash")

		prepRecipient.StateHashes = append(prepRecipient.StateHashes, hash)
		err = prepRecipient.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not save recipient") {
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func setOneTimeUnlock() {
	http.HandleFunc(RecpRPC[17][0], func(w http.ResponseWriter, r *http.Request) {
		// first validate the recipient or anyone would be able to set device ids
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[17][2:], RecpRPC[17][1])
		if err != nil {
			return
		}

		projIndexx := r.FormValue("projIndex")
		seedpwd := r.FormValue("seedpwd")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		err = prepRecipient.SetOneTimeUnlock(projIndex, seedpwd)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not set one time unlock") {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func storeTellerURL() {
	http.HandleFunc(RecpRPC[18][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[18][2:], RecpRPC[18][1])
		if err != nil {
			return
		}

		err = r.ParseForm()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		projIndexx := r.FormValue("projIndex")
		url := r.FormValue("url")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(projIndex)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		if project.RecipientIndex != recipient.U.Index {
			log.Println("recipient indices don't match, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		project.TellerURL = url
		err = project.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		go core.MonitorTeller(projIndex, url)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func storeTellerDetails() {
	http.HandleFunc(RecpRPC[19][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[19][2:], RecpRPC[19][1])
		if err != nil {
			return
		}

		err = r.ParseForm()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		projIndexx := r.FormValue("projIndex")
		url := r.FormValue("url")
		brokerurl := r.FormValue("brokerurl")
		topic := r.FormValue("topic")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "could not convert to integer", messages.ConversionError) {
			return
		}

		if recipient.U.Index != projIndex {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		err = core.AddTellerDetails(projIndex, url, brokerurl, topic)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type recpDashboardHelper struct {
	YourProfile struct {
		Name           string `json:"Name"`
		ActiveProjects int    `json:"Active Projects"`
	} `json:"Your Profile"`

	YourEnergy struct {
		TiCP    string `json:"Total in Current Period"`
		AllTime string `json:"All Time"`
	} `json:"Your Energy"`

	YourWallet struct {
		ProjectWalletBalance float64 `json:"Project Wallet Balance"`
		AutoReload           string  `json:"Auto Reload"`
	} `json:"Your Wallet"`

	NActions struct {
		Notification    string `json:"Notification"`
		ActionsRequired string `json:"Actions Required"`
	} `json:"Notifications & Actions"`

	YourProjects []recpDashboardData `json:"Your Projects"`
}

type recpDashboardData struct {
	Index      int
	ExploreTab map[string]interface{} `json:"Explore Tab"`
	Role       string
	PSA        struct {
		Stage   string
		Actions []string
	} `json:"Project Stage & Actions"`
	ProjectWallets struct {
		Certificates [][]string `json:"Certificates"`
	}
	BillsRewards struct {
		PendingPayments []string `json:"Payments"`
		Link            string   `json:"PastPaymentLink"`
	}
	Documents map[string]interface{} `json:"Documentation and Smart Contracts"`
}

// recpDashboard returns the relevant data needed to populate the recipient dashboard
func recpDashboard() {
	http.HandleFunc(RecpRPC[20][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[20][2:], RecpRPC[20][1])
		if err != nil {
			return
		}

		var ret recpDashboardHelper

		if len(prepRecipient.ReceivedSolarProjectIndices) == 0 {
			erpc.MarshalSend(w, ret)
			return
		}

		project, err := core.RetrieveProject(prepRecipient.ReceivedSolarProjectIndices[0])
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		if len(prepRecipient.U.Name) == 0 {
			ret.YourProfile.Name = "No name set"
		} else {
			ret.YourProfile.Name = prepRecipient.U.Name
		}
		ret.YourProfile.ActiveProjects = len(prepRecipient.ReceivedSolarProjectIndices)

		data, err := erpc.GetRequest(consts.OpenxURL + "/user/tellerfile")
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		type energyStruct struct {
			EnergyTimestamp string `json:"energy_timestamp"`
			Unit            string `json:"unit"`
			Value           uint32 `json:"value"`
			OwnerID         string `json:"owner_id"`
			AssetID         string `json:"asset_id"`
		}

		var EnergyValue uint32
		EnergyValue = 0

		reader := bufio.NewReader(bytes.NewReader(data))

		for {
			var data1 []byte

			for i := 0; i < 7; i++ { // formatted according to the responses received from the lumen unit
				// which is further read by the subscriber
				line, _, err := reader.ReadLine()
				if err != nil {
					break
				}
				data1 = append(data1, line...)
			}

			var x energyStruct
			err = json.Unmarshal(data1, &x)
			if err != nil {
				break
			}

			EnergyValue += x.Value
		}

		ret.YourEnergy.AllTime, err = utils.ToString(EnergyValue)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		ret.YourEnergy.TiCP = ret.YourEnergy.AllTime

		ret.YourWallet.AutoReload = "On"
		ret.NActions.Notification = "None"
		ret.NActions.ActionsRequired = "None"

		if consts.Mainnet {
			ret.YourWallet.ProjectWalletBalance += xlm.GetAssetBalance(project.EscrowPubkey, consts.AnchorUSDCode)
		} else {
			ret.YourWallet.ProjectWalletBalance += xlm.GetAssetBalance(project.EscrowPubkey, consts.StablecoinCode)
		}

		if ret.YourWallet.ProjectWalletBalance < 0 {
			ret.YourWallet.ProjectWalletBalance = 0
		}

		ret.YourProjects = make([]recpDashboardData, len(prepRecipient.ReceivedSolarProjectIndices))
		for i, elem := range prepRecipient.ReceivedSolarProjectIndices {
			var x recpDashboardData
			project, err := core.RetrieveProject(elem)
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.Index = elem
			x.ExploreTab = make(map[string]interface{})
			x.ExploreTab = project.Content.Details["Explore Tab"]
			x.ExploreTab["location"] = project.Content.Details["Explore Tab"]["city"].(string) + ", " + project.Content.Details["Explore Tab"]["state"].(string) + ", " + project.Content.Details["Explore Tab"]["country"].(string)
			x.ExploreTab["money raised"] = project.MoneyRaised
			x.ExploreTab["total value"] = project.TotalValue
			x.Role = "You are an Offtaker"
			sStage, err := utils.ToString(project.Stage)
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.PSA.Stage = "Project is in Stage: " + sStage
			x.PSA.Actions = []string{"Contractor Actions", "No pending action"}
			x.ProjectWallets.Certificates = make([][]string, 2)
			x.ProjectWallets.Certificates[0] = []string{"Carbon & Climate Certificates", "0"}

			pp, err := utils.ToString(float64(EnergyValue) * oracle.MonthlyBill() / 1000000) // /1000 is for kWh
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}

			var temp int64
			if project.DateLastPaid == 0 {
				layout := "Monday, 02-Jan-06 15:04:05 MST"
				t, err := time.Parse(layout, project.DateFunded)
				if err != nil {
					log.Println(err)
					erpc.MarshalSend(w, erpc.StatusInternalServerError)
					return
				}
				temp = t.Unix()
			} else {
				temp = project.DateLastPaid
			}

			dlp, err := utils.ToString(temp + 2419200) // consts.FourWeeksInSecond
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}

			dlpI, err := strconv.ParseInt(dlp, 10, 64)
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}

			dlp = time.Unix(dlpI, 0).String()[0:10]

			xlmUSD, err := tickers.BinanceTicker()
			if erpc.Err(w, err, erpc.StatusInternalServerError, "", messages.TickerError) {
				return
			}

			var wg sync.WaitGroup

			var primNativeBalance, secNativeBalance, primUsdBalance float64

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				primNativeBalance = xlm.GetNativeBalance(prepRecipient.U.StellarWallet.PublicKey) * xlmUSD
				if primNativeBalance < 0 {
					primNativeBalance = 0
				}
			}(&wg)

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				secNativeBalance = xlm.GetNativeBalance(prepRecipient.U.SecondaryWallet.PublicKey) * xlmUSD
				if secNativeBalance < 0 {
					secNativeBalance = 0
				}
			}(&wg)

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				primUsdBalance = xlm.GetAssetBalance(prepRecipient.U.StellarWallet.PublicKey, consts.StablecoinCode)
				if primUsdBalance < 0 {
					primUsdBalance = 0
				}
			}(&wg)

			wg.Wait()

			accBal, err := utils.ToString(primUsdBalance + primNativeBalance)
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}

			x.BillsRewards.PendingPayments = []string{"Your Pending Payment", pp + " due on " + dlp, "Your Account Balance", accBal, "Energy Tariff", "20 ct/kWh"}
			x.BillsRewards.Link = "https://testnet.steexp.com/account/" + prepRecipient.U.StellarWallet.PublicKey + "#transactions"
			x.Documents = make(map[string]interface{})
			x.Documents = project.Content.Details["Documents"]
			ret.YourProjects[i] = x
		}

		erpc.MarshalSend(w, ret)
	})
}

func setCompanyBoolRecp() {
	http.HandleFunc(RecpRPC[21][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[21][2:], RecpRPC[21][1])
		if err != nil {
			return
		}

		err = prepRecipient.SetCompany()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func setCompanyRecp() {
	http.HandleFunc(RecpRPC[22][0], func(w http.ResponseWriter, r *http.Request) {
		prepRecipient, err := recpValidateHelper(w, r, RecpRPC[22][2:], RecpRPC[22][1])
		if err != nil {
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

		err = prepRecipient.SetCompanyDetails(companyType, name, legalName, adminEmail, phoneNumber, address,
			country, city, zipCode, taxIDNumber, role)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func storeTellerEnergy() {
	http.HandleFunc(RecpRPC[23][0], func(w http.ResponseWriter, r *http.Request) {
		recipient, err := recpValidateHelper(w, r, RecpRPC[23][2:], RecpRPC[23][1])
		if err != nil {
			return
		}

		err = r.ParseForm()
		if erpc.Err(w, err, erpc.StatusBadRequest) {
			return
		}

		energy := r.FormValue("energy")

		energyInt, err := utils.ToInt(energy)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		recipient.TellerEnergy = uint32(energyInt)
		recipient.PastTellerEnergy = append(recipient.PastTellerEnergy, uint32(energyInt))

		err = recipient.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
