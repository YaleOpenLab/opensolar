package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/Varunram/essentials/xlm"
	"github.com/Varunram/essentials/xlm/assets"
	"github.com/Varunram/essentials/xlm/wallet"

	"github.com/Varunram/essentials/utils"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	"github.com/YaleOpenLab/opensolar/stablecoin"
)

func sandbox() error {
	project, err := core.RetrieveProject(1)
	if err != nil {
		log.Println(err)
		return err
	}

	password := "password"
	//pwhash := utils.SHA3hash(password)
	seedpwd := "x"
	// invAmount := 4000.0
	run := utils.GetRandomString(5)

	txhash, err := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 100000000000, consts.PlatformSeed)
	if err != nil {
		return err
	}

	log.Println("tx for platform trusting stablecoin:", txhash)

	//exchangeAmount := 1.0

	inv, err := core.NewInvestor("inv"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}

	invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	recp, err := core.NewRecipient("recp"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}

	project.RecipientIndex = recp.U.Index
	project.GuarantorIndex = 1
	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	// start all the teim consuming calls
	// inv.U.Legal = true
	err = inv.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	go xlm.GetXLM(inv.U.StellarWallet.PublicKey)
	go xlm.GetXLM(recp.U.StellarWallet.PublicKey)

	time.Sleep(10 * time.Second) // wait for the accounts to be setup

	if xlm.GetNativeBalance(inv.U.StellarWallet.PublicKey) < 1 {
		return errors.New("inv account not setup")
	}

	if xlm.GetNativeBalance(recp.U.StellarWallet.PublicKey) < 1 {
		return errors.New("recp account not setup")
	}

	log.Println("loading test investor with stablecoin")
	log.Println("INVPUBKEY: ", inv.U.StellarWallet.PublicKey)
	go stablecoin.GetTestStablecoin(inv.U.Username, inv.U.StellarWallet.PublicKey, invSeed, 10000000)

	time.Sleep(35 * time.Second)

	if xlm.GetAssetBalance(inv.U.StellarWallet.PublicKey, consts.StablecoinCode) < 1 {
		return errors.New("stablecoin not present with the investor")
	}

	err = core.Invest(project.Index, inv.U.Index, 4000, invSeed)
	if err != nil {
		log.Println("did not invest in order", err)
		return err
	}

	// now we have to wait for a while and then unlock the project

	// get a token for the recipient

	token, err := getToken(recp.U.Username)
	if err != nil {
		log.Println(err)
		return err
	}

	err = core.UnlockProject(recp.U.Username, token, 1, "x")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("RECIPIENT CREDS: ", recp.U.Username)
	return nil
}

type tokenResponse struct {
	Token string `json:"Token"`
}

func getToken(username string) (string, error) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("pwhash", utils.SHA3hash("password"))

	retdata, err := erpc.PostForm(consts.OpenxURL+"/token", form)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var x tokenResponse

	log.Println("TOKEN : ", string(retdata))

	err = json.Unmarshal(retdata, &x)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return x.Token, nil
}
