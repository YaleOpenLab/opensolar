package main

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"
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
	var project core.Project
	var err error

	project.Index = 1
	project.SeedInvestmentCap = 4000
	project.Stage = 4
	project.MoneyRaised = 0
	project.TotalValue = 4000
	project.OwnershipShift = 0
	project.RecipientIndex = -1  // replace with real indices once created
	project.OriginatorIndex = -1 // replace with real indices once created
	project.GuarantorIndex = -1  // replace with real indices once created
	project.ContractorIndex = -1 // replace with real indices once created
	project.PaybackPeriod = time.Duration(time.Duration(consts.OneWeek) * time.Second)
	project.DeveloperFee = []float64{3000}
	project.Chain = "stellar"
	project.BrokerURL = "mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"
	project.DateInitiated = utils.Timestamp()
	project.DateFunded = utils.Timestamp()
	project.Metadata = "Aibonito Pilot Project"
	project.InvestmentType = "munibond"
	project.TellerURL = ""
	project.BrokerURL = "https://mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"

	// populate the CMS
	project.Content.Details = make(map[string]map[string]interface{})

	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	err = parseCMS("", 1)
	if err != nil {
		log.Println(err)
		return err
	}

	project, err = core.RetrieveProject(project.Index)
	if err != nil {
		log.Println(err)
		return err
	}

	stageString, err := utils.ToString(project.Stage)
	if err != nil {
		log.Println(err)
		return err
	}

	// add details that ashould be parsed from the yaml file here
	project.Name = project.Content.Details["Explore Tab"]["name"].(string)
	project.City = project.Content.Details["Explore Tab"]["city"].(string)
	project.State = project.Content.Details["Explore Tab"]["state"].(string)
	project.Country = project.Content.Details["Explore Tab"]["country"].(string)
	project.MainImage = project.Content.Details["Explore Tab"]["mainimage"].(string)
	project.Content.Details["Explore Tab"]["stage description"] = stageString + " | " + core.GetStageDescription(project.Stage)
	project.Content.Details["Explore Tab"]["location"] = project.Content.Details["Explore Tab"]["city"].(string) + ", " + project.Content.Details["Explore Tab"]["state"].(string) + ", " + project.Content.Details["Explore Tab"]["country"].(string)

	// end all the project fe related changes
	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	// start the main investment loop and the recipient acceptance loop
	password := "password"
	seedpwd := "x"
	invAmount := 4000.0
	run := utils.GetRandomString(5)
	exchangeRate := 10000000.0            // hardcode to 10**7
	stablecoinTrustLimit := 10000000000.0 // some very high limit, this isn't needed since we create the trust line on init, but still
	devFee := 3000.0

	var inv core.Investor
	var recp core.Recipient
	var guar core.Entity
	var dev core.Entity

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		inv, err = core.NewInvestor("mitdci"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created investor")
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		recp, err = core.NewRecipient("fabideas"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created recipient")
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		guar, err = core.NewGuarantor("guarantor"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created guarantor")
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		dev, err = core.NewDeveloper("inversol"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created developer")
	}(&wg)

	wg.Wait()

	invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	dev.PresentContractIndices = append(dev.PresentContractIndices, project.Index)

	devSeed, err := wallet.DecryptSeed(dev.U.StellarWallet.EncryptedSeed, "x")
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	err = dev.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	err = core.AddWaterfallAccount(1, dev.U.StellarWallet.PublicKey, devFee)
	if err != nil {
		log.Println(err)
		return err
	}

	project, err = core.RetrieveProject(1)
	if err != nil {
		log.Println(err)
		return err
	}

	project.OneTimeUnlock = "x" // needed for the developer to be able for the developer to request money
	project.EscrowLock = true
	project.MainDeveloperIndex = dev.U.Index
	project.DeveloperIndices = append(project.DeveloperIndices, dev.U.Index)
	project.RecipientIndex = recp.U.Index
	project.GuarantorIndex = guar.U.Index
	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	var wg1 sync.WaitGroup

	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("loading test investor with stablecoin, pubkey: ", inv.U.StellarWallet.PublicKey)
		stablecoin.GetTestStablecoin(inv.U.Username, inv.U.StellarWallet.PublicKey, seedpwd, exchangeRate)
	}(&wg1)

	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("loading receipient with stablecoin, pubkey: ", recp.U.StellarWallet.PublicKey)
		go stablecoin.GetTestStablecoin(recp.U.Username, recp.U.StellarWallet.PublicKey, seedpwd, exchangeRate)
	}(&wg1)

	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		txhash, err1 := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, stablecoinTrustLimit, consts.PlatformSeed)
		if err != nil {
			log.Fatal(err1)
		}
		log.Println("tx for platform trusting stablecoin:", txhash)
	}(&wg1)

	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("developer trusts stableUSD")
		_, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, devFee, devSeed)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg1)

	var token string

	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		token, err = getToken(recp.U.Username)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg1)

	wg1.Wait()
	time.Sleep(30 * time.Second) // wait for stablecoin to be issued

	if xlm.GetAssetBalance(inv.U.StellarWallet.PublicKey, consts.StablecoinCode) < 1 {
		return errors.New("stablecoin not present with the investor")
	}

	err = core.Invest(project.Index, inv.U.Index, invAmount, invSeed)
	if err != nil {
		log.Println("did not invest in order", err)
		return err
	}

	err = core.UnlockProject(recp.U.Username, token, 1, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	recp.NextPaymentInterval = utils.IntToHumanTime(utils.Unix() + int64(consts.OneWeek))

	err = recp.Save()
	if err != nil {
		log.Println(err)
	}

	log.Println("RECIPIENT CREDS: ", recp.U.Username)
	log.Println("INVESTOR CREDS: ", inv.U.Username)
	log.Println("DEVELOPER CREDS: ", dev.U.Username)
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
