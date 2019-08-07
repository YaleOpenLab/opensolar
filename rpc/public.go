package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
)

// SnInvestor defines a sanitized investor
type SnInvestor struct {
	Name                  string
	InvestedSolarProjects []string
	AmountInvested        float64
	InvestedBonds         []string
	InvestedCoops         []string
	PublicKey             string
	Reputation            float64
}

// SnRecipient defines a sanitized recipient
type SnRecipient struct {
	Name                  string
	PublicKey             string
	ReceivedSolarProjects []string
	Reputation            float64
}

// SnUser defines a sanitized user
type SnUser struct {
	Name       string
	PublicKey  string
	Reputation float64
}

func setupPublicRoutes() {
	getAllInvestorsPublic()
	getAllRecipientsPublic()
	getTopReputationPublic()
	getInvTopReputationPublic()
	getRecpTopReputationPublic()
	getUserInfo()
}

// public contains all the RPC routes that we explicitly intend to make public. Other
// routes such as the invest route are things we could make private as well, but that
// doesn't change the security model since we ask for username+pwauth

// sanitizeInvestor removes sensitive fields frm the investor struct in order to be able
// to return the investor field in a public route
func sanitizeInvestor(investor opensolar.Investor) SnInvestor {
	// this is a public route, so we shouldn't ideally return all parameters that are present
	// in the investor struct
	var sanitize SnInvestor
	sanitize.Name = investor.U.Name
	sanitize.InvestedSolarProjects = investor.InvestedSolarProjects
	sanitize.AmountInvested = investor.AmountInvested
	sanitize.PublicKey = investor.U.StellarWallet.PublicKey
	sanitize.Reputation = investor.U.Reputation
	return sanitize
}

// sanitizeRecipient removes sensitive fields from the recipient struct in order to be
// able to return the recipient fields in a public route
func sanitizeRecipient(recipient opensolar.Recipient) SnRecipient {
	// this is a public route, so we shouldn't ideally return all parameters that are present
	// in the investor struct
	var sanitize SnRecipient
	sanitize.Name = recipient.U.Name
	sanitize.PublicKey = recipient.U.StellarWallet.PublicKey
	sanitize.Reputation = recipient.U.Reputation
	sanitize.ReceivedSolarProjects = recipient.ReceivedSolarProjects
	return sanitize
}

// sanitizeAllInvestors sanitizes an array of investors
func sanitizeAllInvestors(investors []opensolar.Investor) []SnInvestor {
	var arr []SnInvestor
	for _, elem := range investors {
		arr = append(arr, sanitizeInvestor(elem))
	}
	return arr
}

// sanitizeUser sanitizes a particular user
func sanitizeUser(user database.User) SnUser {
	var sanitize SnUser
	sanitize.Name = user.Name
	sanitize.PublicKey = user.StellarWallet.PublicKey
	sanitize.Reputation = user.Reputation
	return sanitize
}

// sanitizeAllRecipients sanitizes an array of recipients
func sanitizeAllRecipients(recipients []opensolar.Recipient) []SnRecipient {
	var arr []SnRecipient
	for _, elem := range recipients {
		arr = append(arr, sanitizeRecipient(elem))
	}
	return arr
}

// sanitizeAllUsers sanitizes an arryay of users
func sanitizeAllUsers(users []database.User) []SnUser {
	var arr []SnUser
	for _, elem := range users {
		arr = append(arr, sanitizeUser(elem))
	}
	return arr
}

// getAllInvestors gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllInvestorsPublic() {
	http.HandleFunc("/public/investor/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		investors, err := opensolar.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(investors)
		erpc.MarshalSend(w, sInvestors)
	})
}

// getAllRecipients gets a list of all the investors in the system so that we can
// display it to some entity that is interested to view such stats
func getAllRecipientsPublic() {
	http.HandleFunc("/public/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		recipients, err := opensolar.RetrieveAllRecipients()
		if err != nil {
			log.Println("did not retrieve all recipients", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(recipients)
		erpc.MarshalSend(w, sRecipients)
	})
}

// this is to publish a list of the users with the best feedback in the system in order
// to award them badges or something similar
func getTopReputationPublic() {
	http.HandleFunc("/public/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allUsers, err := database.TopReputationUsers()
		if err != nil {
			log.Println("did not retrive all top reputation users", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sUsers := sanitizeAllUsers(allUsers)
		erpc.MarshalSend(w, sUsers)
	})
}

// getRecpTopReputationPublic gets a list of the recipients who have the best reputation on the platform
func getRecpTopReputationPublic() {
	http.HandleFunc("/public/recipient/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allRecps, err := opensolar.TopReputationRecipients()
		if err != nil {
			log.Println("did not retrieve all top reputaiton recipients", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(allRecps)
		erpc.MarshalSend(w, sRecipients)
	})
}

// getInvTopReputationPublic gets a lsit of the investors who have the best reputation on the platform
func getInvTopReputationPublic() {
	http.HandleFunc("/public/investor/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allInvs, err := opensolar.TopReputationInvestors()
		if err != nil {
			log.Println("did not retrieve all top reputation investors", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(allInvs)
		erpc.MarshalSend(w, sInvestors)
	})
}

func getUserInfo() {
	http.HandleFunc("/public/user", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["index"] == nil {
			log.Println("no index retrieved, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		user, err := database.RetrieveUser(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		sUser := sanitizeUser(user)
		erpc.MarshalSend(w, sUser)
	})
}
