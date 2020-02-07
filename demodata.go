package main

import (
	"log"
	"strings"

	core "github.com/YaleOpenLab/opensolar/core"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func demoData() error {
	var project core.Project
	var err error

	project.Name = "5kW Solar at FabIDEAS Coop - Pilot 1"
	project.City = "Aibonito"
	project.State = "Puerto Rico"
	project.Country = "USA"
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
	project.PaybackPeriod = 4    // four weeks payback time
	project.Acquisition = "2025"
	project.Chain = "stellar"
	project.BrokerUrl = "mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"
	project.Metadata = "Aibonito Pilot Project"
	project.InvestmentType = "munibond"
	project.TellerUrl = ""
	project.BrokerUrl = "https://mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"

	// populate the CMS
	// project.Content.DetailPageStub.Box
	project.Content.DetailPageStub.Box = make(map[string]interface{})
	project.Content.OtherDetails = make(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Terms = make(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Overview = make(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Project = make(map[string]interface{})
	project.Content.DetailPageStub.Tabs.StageForecast = make(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Documents = make(map[string]interface{})

	project.Content.DetailPageStub.Box["Name"] = project.Name
	project.Content.DetailPageStub.Box["Location"] = project.City + ", " + project.State + ", " + project.Country
	project.Content.DetailPageStub.Box["Maturity"] = project.Acquisition
	project.Content.DetailPageStub.Box["Money Raised"] = project.MoneyRaised
	project.Content.DetailPageStub.Box["Total Value"] = project.TotalValue

	// project.Content.DetailPageStub.Tabs.Overview
	/*
		recp, err := core.NewRecipient("aibonito", utils.SHA3hash("password"), "password", "Maria Pastor")
		if err != nil {
			return err
		}
		project.RecipientIndex = recp.U.Index

		orig, err := core.NewOriginator("mwainstein", "password", "password", "Martin Wainstein")
		if err != nil {
			return err
		}
		project.OriginatorIndex = orig.U.Index

		cont, err := NewContractor("contractor", "password", "password", "Contractor Name")
		if err != nil {
			return err
		}
		project.ContractorIndex = cont.U.Index

		dev, err := core.NewDeveloper("developer", "password", "password", "Developer Name")
		if err != nil {
			return err
		}
		project.MainDeveloperIndex = dev.U.Index

		guar, err := core.NewGuarantor("guarantor", "password", "password", "Guarantor Name")
		if err != nil {
			return err
		}
		project.GuarantorIndex = guar.U.Index
	*/

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

	/*
		txhash, err := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, consts.PlatformSeed)
		if err != nil {
			return err
		}

		log.Println("tx for platform trusting stablecoin:", txhash)

		password := "password"
		//pwhash := utils.SHA3hash(password)
		seedpwd := "x"
		//exchangeAmount := 1.0
		invAmount := 4000.0
		run := utils.GetRandomString(5)

		inv, err := core.NewInvestor("inv"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Println(err)
			return err
		}

		// inv.U.Legal = true
		err = inv.Save()
		if err != nil {
			log.Println(err)
			return err
		}

		err = xlm.GetXLM(inv.U.StellarWallet.PublicKey)
		if err != nil {
			log.Println("could not get XLM: ", err)
			return err
		}

		invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, seedpwd)
		if err != nil {
			log.Println(err)
			return err
		}

		err = stablecoin.GetTestStablecoin(inv.U.Username, inv.U.StellarWallet.PublicKey, invSeed, 1000000)
		if err != nil {
			log.Println(err)
			return err
		}

		recp, err := core.NewRecipient("recp"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Println(err)
			return err
		}

		err = xlm.GetXLM(recp.U.StellarWallet.PublicKey)
		if err != nil {
			log.Println("could not get XLM: ", err)
			return err
		}

		project.RecipientIndex = recp.U.Index
		project.GuarantorIndex = 1
		err = project.Save()
		if err != nil {
			log.Println(err)
			return err
		}

		err = core.Invest(project.Index, inv.U.Index, invAmount, invSeed)
		if err != nil {
			log.Println("did not invest in order", err)
			return err
		}

		log.Println("RECIPIENT CREDS: ", recp.U.Username, recp.U.AccessToken, recp.U.Pwhash, project.Index)
	*/
	return nil
}

func ifString(x interface{}) bool {
	switch x.(type) {
	case string:
		return true
	default:
		return false
	}
}

func ifMapStringInterface(x interface{}) bool {
	switch x.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

// parseCMS parses the yaml file and converts it into the CMS format we have
func parseCMS(fileName string, projIndex int) error {
	viper.SetConfigType("yaml")
	// viper.SetConfigName(fileName)
	viper.SetConfigName("cms")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	project, err := core.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	etm := viper.Get("Explore Tab Modal").(map[string]interface{})
	for key, value := range etm {
		titleKey := strings.Title(key)
		project.Content.DetailPageStub.Box[titleKey] = value
	}

	od := viper.Get("Other Details").(map[string]interface{})
	for key, value := range od {
		titleKey := strings.Title(key)
		project.Content.OtherDetails[titleKey] = value
	}

	terms := viper.Get("Terms").(map[string]interface{})
	for key, value := range terms {
		titleKey := strings.Title(key)
		if ifString(value) {
			project.Content.DetailPageStub.Tabs.Terms[titleKey] = value
		}
		if ifMapStringInterface(value) {
			msi := make(map[string]interface{})
			for tKey, tValue := range value.(map[string]interface{}) {
				msi[strings.Title(tKey)] = tValue
			}
			project.Content.DetailPageStub.Tabs.Terms[strings.Title(key)] = make(map[string]interface{})
			project.Content.DetailPageStub.Tabs.Terms[strings.Title(key)] = msi
		}
	}

	overview := viper.Get("overview").(map[string]interface{})
	for key, value := range overview {
		titleKey := strings.Title(key)
		if ifString(value) {
			project.Content.DetailPageStub.Tabs.Overview[titleKey] = value
		}
		if ifMapStringInterface(value) {
			msi := make(map[string]interface{})
			for tKey, tValue := range value.(map[string]interface{}) {
				msi[strings.Title(tKey)] = tValue
			}
			project.Content.DetailPageStub.Tabs.Overview[strings.Title(key)] = make(map[string]interface{})
			project.Content.DetailPageStub.Tabs.Overview[strings.Title(key)] = msi
		}
	}

	projDetails := viper.Get("project details").(map[string]interface{})
	for key, value := range projDetails {
		titleKey := strings.Title(key)
		if ifString(value) {
			project.Content.DetailPageStub.Tabs.Project[titleKey] = value
		}
		if ifMapStringInterface(value) {
			msi := make(map[string]interface{})
			for tKey, tValue := range value.(map[string]interface{}) {
				msi[strings.Title(tKey)] = tValue
				log.Println("TKEY: ", strings.Title(tKey))
			}
			project.Content.DetailPageStub.Tabs.Project[strings.Title(key)] = make(map[string]interface{})
			project.Content.DetailPageStub.Tabs.Project[strings.Title(key)] = msi
		}
	}

	sForecast := viper.Get("stage").(map[string]interface{})
	for key, value := range sForecast {
		titleKey := strings.Title(key)
		if ifString(value) {
			project.Content.DetailPageStub.Tabs.StageForecast[titleKey] = value
		}
		if ifMapStringInterface(value) {
			msi := make(map[string]interface{})
			for tKey, tValue := range value.(map[string]interface{}) {
				msi[strings.Title(tKey)] = tValue
				log.Println("TKEY: ", strings.Title(tKey))
			}
			project.Content.DetailPageStub.Tabs.StageForecast[strings.Title(key)] = make(map[string]interface{})
			project.Content.DetailPageStub.Tabs.StageForecast[strings.Title(key)] = msi
		}
	}

	documents := viper.Get("documents").(map[string]interface{})
	for key, value := range documents {
		titleKey := strings.Title(key)
		if ifString(value) {
			project.Content.DetailPageStub.Tabs.Documents[titleKey] = value
		}
		if ifMapStringInterface(value) {
			msi := make(map[string]interface{})
			for tKey, tValue := range value.(map[string]interface{}) {
				msi[strings.Title(tKey)] = tValue
				log.Println("TKEY: ", strings.Title(tKey))
			}
			project.Content.DetailPageStub.Tabs.Documents[strings.Title(key)] = make(map[string]interface{})
			project.Content.DetailPageStub.Tabs.Documents[strings.Title(key)] = msi
		}
	}

	err = project.Save()
	if err != nil {
		return err
	}

	return nil
}
