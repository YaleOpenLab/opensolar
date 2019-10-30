// +build all travis

package core

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	opensolarconsts "github.com/YaleOpenLab/opensolar/consts"
	openxconsts "github.com/YaleOpenLab/openx/consts"
	openxdb "github.com/YaleOpenLab/openx/database"
	build "github.com/stellar/go/txnbuild"
)

// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	openxconsts.SetConsts(false)
	opensolarconsts.SetTnConsts()
	os.Remove(openxconsts.DbDir + "/opensolar.db")
	CreateHomeDir() // create home directory if it doesn't exist yet
	openxdb.CreateHomeDir()
	db, err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	inv, err := NewInvestor("investor1", "blah", "blah", "Investor1")
	if err != nil {
		t.Fatal(err)
	}
	inv.U.AccessToken = "ACCESSTOKEN"
	inv.U.AccessTokenTimeout = utils.Unix() + 1000000
	err = inv.U.Save()
	if err != nil {
		t.Fatal(err)
	}
	xc, err := RetrieveAllInvestors()
	if len(xc) != 2 { // one is created by default
		log.Println("len of all inv: ", len(xc))
		t.Fatal(err)
	}
	log.Println(xc[0].U, len(xc))
	// try retrieving existing stuff
	inv1, err := RetrieveInvestor(1)
	if err != nil {
		t.Fatal(err)
	}
	if inv1.U.Name != "Investor1" {
		t.Fatalf("Investor names don't match, quitting!")
	}
	// func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	recp, err := NewRecipient("recipient1", "blah", "blah", "Recipient1")
	if err != nil {
		t.Fatal(err)
	}
	recp.U.AccessToken = "ACCESSTOKEN"
	recp.U.AccessTokenTimeout = utils.Unix() + 1000000

	err = recp.U.Save()
	if err != nil {
		t.Fatal(err)
	}
	rec1, err := RetrieveRecipient(recp.U.Index)
	if err != nil {
		t.Fatal(err)
	}

	if rec1.U.Name != "Recipient1" {
		t.Fatalf("Recipient usernames don't match. quitting!")
	}

	allRec, err := RetrieveAllRecipients()
	if err != nil || len(allRec) != 2 { // one is created by default
		log.Println("length of all recipients not 2")
		t.Fatal(err)
	}

	user, err := NewUser("user1", "blah", "blah", "User1")
	if err != nil {
		t.Fatal(err)
	}
	user.AccessToken = "ACCESSTOKEN"
	user.AccessTokenTimeout = utils.Unix() + 1000000
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	_, err = RetrieveInvestor(1000)
	if err == nil {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	_, err = RetrieveRecipient(1000)
	if err == nil {
		t.Fatalf("Recipient shouldn't exist, but does, quitting!")
	}

	user1, err := RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	if user1.Name != "User1" {
		t.Fatalf("Usernames don't match. quitting!")
	}

	tmpuser, _ := RetrieveUser(1000)
	if tmpuser.Index != 0 {
		t.Fatalf("Investor shouldn't exist, but does, quitting!")
	}

	allUsers, err := openxdb.RetrieveAllUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(allUsers) != 3 {
		t.Fatalf("Unknown users existing, quitting!")
	}

	// check if each of the validate functions work
	_, err = ValidateInvestor("investor1", "ACCESSTOKEN")
	if err != nil {
		log.Println(err)
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("recipient1", "ACCESSTOKEN")
	if err != nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateUser("user1", "ACCESSTOKEN")
	if err != nil {
		log.Println(err)
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateInvestor("blah", "ACCESSTOKEN")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}

	_, err = ValidateRecipient("blah", "ACCESSTOKEN")
	if err == nil {
		t.Fatalf("Data in bucket morphed, quitting!")
	}
	// check voting balance routes
	voteBalance := inv.VotingBalance
	err = inv.ChangeVotingBalance(10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 10000 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	err = inv.ChangeVotingBalance(-10000)
	if err != nil {
		t.Fatal(err)
	}
	if inv.VotingBalance-voteBalance != 0 {
		t.Fatalf("Voting Balance not added, quitting!")
	}
	// check CanInvest Route
	if inv.CanInvest(1000) {
		t.Fatalf("CanInvest Returns true!")
	}
	err = user.GenKeys("blah")
	if err != nil {
		t.Fatalf("Not able to generate keys, quitting!")
	}

	// check the asset functions below. For some weird reason, placing these tests
	// above confuses the other routes, so placing everything here so that we can
	// isolate them from the other routes.
	err = xlm.GetXLM(recp.U.StellarWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	err = xlm.GetXLM(inv.U.StellarWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	err = inv.U.IncreaseTrustLimit("blah", 10)
	if err != nil {
		t.Fatal(err)
	}
	_ = build.CreditAsset{"blah", recp.U.StellarWallet.PublicKey}
	invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, "blah")
	if err != nil {
		t.Fatal(err)
	}
	hash, err := assets.TrustAsset("blah", recp.U.StellarWallet.PublicKey, 100, invSeed)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("HASH IS: ", hash)
	_, err = assets.TrustAsset("blah", recp.U.StellarWallet.PublicKey, -1, "blah")
	if err == nil {
		t.Fatalf("can trust asset with invalid s eed!")
	}
	pkSeed, _, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_ = build.CreditAsset{"blah2", pkSeed} // this account doesn't exist yet, so this should fail
	_, err = assets.TrustAsset("blah2", "", -1, "blah")
	if err == nil {
		t.Fatalf("can trust invalid asset")
	}
	_, err = openxdb.RetrieveAllUsersWithoutKyc()
	if err != nil {
		t.Fatal(err)
	}
	_, err = openxdb.RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}
	err = user.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = user.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = openxdb.TopReputationUsers()
	if err != nil {
		t.Fatal(err)
	}
	err = user.Authorize(user.Index)
	if err == nil {
		t.Fatalf("Not able to catch inspector permission error")
	}
	err = user.SetBan(100)
	if err == nil {
		t.Fatalf("able to ban a user even though person is not an inspector, quitting")
	}
	err = openxdb.AddInspector(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	user, err = openxdb.RetrieveUser(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	err = user.Authorize(user.Index)
	if err != nil {
		t.Fatal(err)
	}
	err = user.SetBan(user.Index)
	if err == nil {
		t.Fatalf("able to  set a ban on self, quitting")
	}
	err = user.SetBan(-1)
	if err == nil {
		t.Fatalf("able to set ban on user who doesn't exist, quitting")
	}
	var banTest openxdb.User
	banTest.Index = 1000
	err = banTest.Save()
	if err != nil {
		t.Fatalf("not able to save user for banning, quitting")
	}
	user.Admin = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = user.SetBan(1000)
	if err != nil {
		log.Println(err)
		t.Fatalf("Not able to set ban on legitimate user, quitting")
	}
	err = user.SetBan(1000)
	if err != nil {
		t.Fatalf("Able to set ban on user even if ban is already set, quitting")
	}
	user.Kyc = true
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = user.Authorize(user.Index)
	if err == nil {
		t.Fatalf("Able to authorize KYC'd user, exiting!")
	}
	_, err = openxdb.RetrieveAllUsersWithKyc()
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationRecipients()
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.ChangeReputation(-1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = recp.U.AddEmail("blah@blah.com")
	if err != nil {
		t.Fatal(err)
	}
	_, err = openxdb.SearchWithEmailId("blah@blah.com")
	if err != nil {
		t.Fatal(err)
	}
	testuser, _ := openxdb.SearchWithEmailId("blahx@blah.com")
	if testuser.StellarWallet.PublicKey != "" {
		t.Fatalf("user with invalid email exists")
	}
	err = xlm.GetXLM(inv.U.SecondaryWallet.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.MoveFundsFromSecondaryWallet(10, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.MoveFundsFromSecondaryWallet(-1, "blah")
	if err == nil {
		t.Fatalf("not able to catch invalid amount error")
	}
	err = inv.U.SweepSecondaryWallet("blah")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.SweepSecondaryWallet("invalidseedpwd")
	if err == nil {
		t.Fatalf("no able to catch invalid seedpwd")
	}
	_, err = TopReputationInvestors()
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.ChangeReputation(1.0)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.U.ChangeReputation(-1.0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationRecipients()
	if err != nil {
		t.Fatal(err)
	}
	_, err = TopReputationInvestors()
	if err != nil {
		t.Fatal(err)
	}
	var blah Investor
	blah.U = new(openxdb.User)
	blah.U.Name = "Cool"
	blahBytes, err := json.Marshal(blah)
	if err != nil {
		t.Fatal(err)
	}

	var uBlah Investor
	err = json.Unmarshal(blahBytes, &uBlah)
	if err != nil {
		t.Fatal(err)
	}

	err = DeleteKeyFromBucket(inv.U.Index, InvestorBucket)
	if err != nil {
		t.Fatal(err)
	}

	os.Remove(openxconsts.DbDir + "/openx.db")
}
