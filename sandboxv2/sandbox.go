package sandbox

import (
	"encoding/json"
	"io/ioutil"
	"log"

	core "github.com/YaleOpenLab/opensolar/core"
	utils "github.com/Varunram/essentials/utils"
)

type TestStruct struct {
	Field1 string `json:Field1`
	Field2 struct {
		Field21 string `json:Field21`
		Field22 string `json:Field21`
	} `json:Field2`
	Field3 struct {
		Field31 struct {
			Field311 string `json:Field311`
			Field312 string `json:Field312`
		} `json:Field31`
	} `json:Field3`
	Field4 []string   `json:Field4`
	Field5 [][]string `json:Field5`
}

func ReadValues() (TestStruct, error) {
	var err error
	var x TestStruct

	data, err := ioutil.ReadFile("sandboxv2/values.json")
	if err != nil {
		return x, err
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, err
	}

	return x, nil
}

func populateUsers() error {
	// func NewInvestor(uname string, pwd string, seedpwd string, Name string) (Investor, error) {
	// func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	// func NewEntity(uname string, pwhash string, seedpwd string, name string role string) (Entity, error) {
	pwhash := utils.SHA3hash("password")
	orig, err := core.NewOriginator("orig", pwhash, "x", "orig")
	if err != nil {
		return err
	}

	inv, err := core.NewInvestor("inv", pwhash, "x", "inv")
	if err != nil {
		return err
	}

	recp, err := core.NewRecipient("recp", pwhash, "x", "recp")
	if err != nil {
		return err
	}

	dev, err := core.NewDeveloper("dev", pwhash, "x", "dev")
	if err != nil {
		return err
	}

	guar, err := core.NewGuarantor("guar", pwhash, "x", "guar")
	if err != nil {
		return err
	}

	log.Printf("Pubkeys list:\nInvestor:%s\nRecipient:%s\nDeveloper:%s\n",
		orig.U.StellarWallet.PublicKey, inv.U.StellarWallet.PublicKey, recp.U.StellarWallet.PublicKey, dev.U.StellarWallet.PublicKey, guar.U.StellarWallet.PublicKey)

	return nil
}

func Test() error {
	log.Println("populating test sandbox endpoints")
	var err error

	values, err := ReadValues()
	if err != nil {
		return err
	}

	log.Println(values)
	return nil
}
