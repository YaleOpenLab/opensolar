package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	core "github.com/YaleOpenLab/opensolar/core"
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
	getInvTopReputationPublic()
	getRecpTopReputationPublic()
}

// sanitizeInvestor removes sensitive fields from the investor struct
func sanitizeInvestor(investor core.Investor) SnInvestor {
	var sanitize SnInvestor
	sanitize.Name = investor.U.Name
	sanitize.InvestedSolarProjects = investor.InvestedSolarProjects
	sanitize.AmountInvested = investor.AmountInvested
	sanitize.PublicKey = investor.U.StellarWallet.PublicKey
	sanitize.Reputation = investor.U.Reputation
	return sanitize
}

// sanitizeRecipient removes sensitive fields from the recipient struct
func sanitizeRecipient(recipient core.Recipient) SnRecipient {
	var sanitize SnRecipient
	sanitize.Name = recipient.U.Name
	sanitize.PublicKey = recipient.U.StellarWallet.PublicKey
	sanitize.Reputation = recipient.U.Reputation
	sanitize.ReceivedSolarProjects = recipient.ReceivedSolarProjects
	return sanitize
}

// sanitizeAllInvestors sanitizes an array of investors
func sanitizeAllInvestors(investors []core.Investor) []SnInvestor {
	var arr []SnInvestor
	for _, elem := range investors {
		arr = append(arr, sanitizeInvestor(elem))
	}
	return arr
}

// sanitizeAllRecipients sanitizes an array of recipients
func sanitizeAllRecipients(recipients []core.Recipient) []SnRecipient {
	var arr []SnRecipient
	for _, elem := range recipients {
		arr = append(arr, sanitizeRecipient(elem))
	}
	return arr
}

// getAllInvestors gets a list of all the investors in the database
func getAllInvestorsPublic() {
	http.HandleFunc("/public/investor/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		investors, err := core.RetrieveAllInvestors()
		if err != nil {
			log.Println("did not retrieve all investors", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(investors)
		erpc.MarshalSend(w, sInvestors)
	})
}

// getAllRecipients gets a list of all the investors in the database
func getAllRecipientsPublic() {
	http.HandleFunc("/public/recipient/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		recipients, err := core.RetrieveAllRecipients()
		if err != nil {
			log.Println("did not retrieve all recipients", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(recipients)
		erpc.MarshalSend(w, sRecipients)
	})
}

// getRecpTopReputationPublic gets a list of the recipients sorted by descending order of reputation
func getRecpTopReputationPublic() {
	http.HandleFunc("/public/recipient/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		allRecps, err := core.TopReputationRecipients()
		if err != nil {
			log.Println("did not retrieve all top reputaiton recipients", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sRecipients := sanitizeAllRecipients(allRecps)
		erpc.MarshalSend(w, sRecipients)
	})
}

// getInvTopReputationPublic gets a list of the investors sorted by descending order of reputation
func getInvTopReputationPublic() {
	http.HandleFunc("/public/investor/reputation/top", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		allInvs, err := core.TopReputationInvestors()
		if err != nil {
			log.Println("did not retrieve all top reputation investors", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		sInvestors := sanitizeAllInvestors(allInvs)
		erpc.MarshalSend(w, sInvestors)
	})
}
