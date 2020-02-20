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
	seedpwd := "x"
	invAmount := 4000.0
	run := utils.GetRandomString(5)

	txhash, err := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 100000000000, consts.PlatformSeed)
	if err != nil {
		return err
	}
	log.Println("tx for platform trusting stablecoin:", txhash)

	inv, err := core.NewInvestor("inv"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}
	// inv.U.Legal = true

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

	guar, err := core.NewGuarantor("guar"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}

	dev, err := core.NewDeveloper("inversol"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}

	dev.PresentContractIndices = append(dev.PresentContractIndices, project.Index)

	err = dev.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	err = core.AddWaterfallAccount(1, dev.U.StellarWallet.PublicKey, 3000)
	if err != nil {
		log.Println(err)
		return err
	}

	project.OneTimeUnlock = "x" // needed for the developer to be able for the developer to request money
	project.MainDeveloperIndex = dev.U.Index
	project.DeveloperIndices = append(project.DeveloperIndices, dev.U.Index)
	project.RecipientIndex = recp.U.Index
	project.GuarantorIndex = guar.U.Index
	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	// start all the time consuming calls

	go xlm.GetXLM(inv.U.StellarWallet.PublicKey)
	go xlm.GetXLM(recp.U.StellarWallet.PublicKey)
	go xlm.GetXLM(dev.U.StellarWallet.PublicKey)

	time.Sleep(10 * time.Second) // wait for the accounts to be setup

	if xlm.GetNativeBalance(inv.U.StellarWallet.PublicKey) < 1 {
		return errors.New("inv account not setup")
	}

	if xlm.GetNativeBalance(recp.U.StellarWallet.PublicKey) < 1 {
		return errors.New("recp account not setup")
	}

	log.Println("loading test investor with stablecoin, pubkey: ", inv.U.StellarWallet.PublicKey)

	go stablecoin.GetTestStablecoin(inv.U.Username, inv.U.StellarWallet.PublicKey, seedpwd, 10000000)
	time.Sleep(35 * time.Second)

	if xlm.GetAssetBalance(inv.U.StellarWallet.PublicKey, consts.StablecoinCode) < 1 {
		return errors.New("stablecoin not present with the investor")
	}

	err = core.Invest(project.Index, inv.U.Index, invAmount, invSeed)
	if err != nil {
		log.Println("did not invest in order", err)
		return err
	}

	// get a token for the recipient
	token, err := getToken(recp.U.Username)
	if err != nil {
		log.Println(err)
		return err
	}

	err = core.UnlockProject(recp.U.Username, token, 1, seedpwd)
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
