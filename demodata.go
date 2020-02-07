package main

import (
	"log"

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
	// project.Content.Details.ExploreTab
	project.Content.Details = make(map[string]map[string]interface{})

	// project.Content.Details.Tabs.Overview
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

func max(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	max := arr[0]
	for _, elem := range arr {
		if elem > max {
			max = elem
		}
	}
	return max
}

func findInterfaceDepth(x interface{}) int {
	switch x.(type) {
	case []interface{}:
		var depths []int
		y := x.([]interface{})
		for _, val := range y {
			depths = append(depths, findInterfaceDepth(val))
		}
		return max(depths)
	case map[string]interface{}:
		var depths []int
		y := x.(map[string]interface{})
		for _, val := range y {
			depths = append(depths, findInterfaceDepth(val))
		}
		return 1 + max(depths)
	default:
		return 0
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

	project.Content.Keys = []string{"Explore Tab", "Other Details", "Terms", "Overview", "project details", "Stage", "Documents"}

	for _, key := range project.Content.Keys {
		if !viper.IsSet(key) {
			log.Println("required content" + key + " not found for project, quitting")
			return errors.New("required content not found for project, quitting")
		}
	}

	for _, key := range project.Content.Keys {
		depth := findInterfaceDepth(viper.Get(key))
		switch depth {
		case 1:
			msi := viper.Get(key).(map[string]interface{})
			project.Content.Details[key] = make(map[string]interface{})
			for key1, value1 := range msi {
				project.Content.Details[key][key1] = value1
			}
		case 2:
			msmsi := viper.Get(key).(map[string]interface{})
			project.Content.Details[key] = make(map[string]interface{})
			for key1, value1 := range msmsi {
				depth := findInterfaceDepth(value1)
				switch depth {
				case 0:
					project.Content.Details[key][key1] = value1.(interface{})
				case 1:
					msi := value1.(map[string]interface{})
					project.Content.Details[key][key1] = make(map[string]interface{})
					for key2, value2 := range msi {
						switch value2.(type) {
						case []interface{}:
							msi[key2] = value2.([]interface{})
						case interface{}:
							msi[key2] = value2.(interface{})
						}
					}
					project.Content.Details[key][key1] = msi
				}
			}
		case 3:
			msmsmsi := viper.Get(key).(map[string]interface{})
			project.Content.Details[key] = make(map[string]interface{})
			for key1, value1 := range msmsmsi {
				depth := findInterfaceDepth(value1)
				switch depth {
				case 0:
					project.Content.Details[key][key1] = value1.(interface{})
				case 1:
					msi := value1.(map[string]interface{})
					project.Content.Details[key][key1] = make(map[string]interface{})
					for key2, value2 := range msi {
						switch value2.(type) {
						case []interface{}:
							msi[key2] = value2.([]interface{})
						case interface{}:
							msi[key2] = value2.(interface{})
						}
					}
					project.Content.Details[key][key1] = msi
				case 2:
					msmsi := viper.Get(key).(map[string]interface{})
					for key2, value2 := range msmsi {
						depth := findInterfaceDepth(value2)
						switch depth {
						case 0:
							msmsi[key2] = value2.(interface{})
						case 1:
							msi := value2.(map[string]interface{})
							msmsi[key2] = make(map[string]interface{})
							for key3, value3 := range msi {
								switch value3.(type) {
								case []interface{}:
									msi[key3] = value3.([]interface{})
								case interface{}:
									msi[key3] = value3.(interface{})
								}
							}
							msmsi[key2] = msi
							project.Content.Details[key][key1] = make(map[string]map[string]interface{})
							project.Content.Details[key][key1] = msmsi
						}
					}
				}
			}
		default:
			log.Println("cool")
		}
	}

	err = project.Save()
	if err != nil {
		return err
	}

	return nil
}
