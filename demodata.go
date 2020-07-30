package main

import (
	"log"
	"strconv"

	"github.com/Varunram/essentials/utils"

	core "github.com/YaleOpenLab/opensolar/core"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func demoData() error {
	demoProjects := createDemoProjects()
	var err error

	for _, project := range demoProjects{
		err = project.Save()

		if err != nil {
			log.Println(err)
			return err
		}

		filename := "./demodata/cms" + strconv.Itoa(project.Index)

		err = parseCMS(filename, project.Index)
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

		// add details that should be parsed from the yaml file here
		project.Name = project.Content.Details["Explore Tab"]["name"].(string)
		project.City = project.Content.Details["Explore Tab"]["city"].(string)
		project.State = project.Content.Details["Explore Tab"]["state"].(string)
		project.Country = project.Content.Details["Explore Tab"]["country"].(string)
		project.MainImage = project.Content.Details["Explore Tab"]["mainimage"].(string)
		project.Content.Details["Explore Tab"]["stage description"] = stageString + " | " + core.GetStageDescription(project.Stage)
		project.Content.Details["Explore Tab"]["location"] = project.Content.Details["Explore Tab"]["city"].(string) + ", " + project.Content.Details["Explore Tab"]["state"].(string) + ", " + project.Content.Details["Explore Tab"]["country"].(string)

		password := "password"
		seedpwd := "x"
		run := utils.GetRandomString(5)

		inv, err := core.NewInvestor("mitdci"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Println(err)
			return err
		}

		recp, err := core.NewRecipient("fabideas"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Println(err)
			return err
		}

		dev, err := core.NewDeveloper("inversol"+run, password, seedpwd, "varunramganesh@gmail.com")
		if err != nil {
			log.Println(err)
			return err
		}

		inv.InvestedSolarProjectsIndices = append(inv.InvestedSolarProjectsIndices, project.Index)
		recp.ReceivedSolarProjectIndices = append(recp.ReceivedSolarProjectIndices, project.Index)
		dev.PresentContractIndices = append(dev.PresentContractIndices, project.Index)

		project.MainDeveloperIndex = dev.U.Index
		project.DeveloperFee = []float64{3000}
		project.RecipientIndex = recp.U.Index
		project.GuarantorIndex = 1

		err = inv.Save()
		if err != nil {
			log.Println(err)
			return err
		}

		err = recp.Save()
		if err != nil {
			log.Println(err)
			return err
		}

		err = dev.Save()
		if err != nil {
			log.Println(err)
			return err
		}

		err = project.Save()

		if err != nil {
			log.Println(err)
			return err
		}
	}
	return err
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

func convert1(x interface{}) interface{} {
	msi := x.(map[string]interface{})
	for key, value := range msi {
		switch value.(type) {
		case []interface{}:
			msi[key] = value.([]interface{})
		case interface{}:
			msi[key] = value.(interface{})
		}
	}
	return msi
}

func convert2(x map[string]interface{}) map[string]interface{} {
	temp := make(map[string]interface{})
	for key1, value1 := range x {
		depth := findInterfaceDepth(value1)
		switch depth {
		case 0:
			temp[key1] = value1.(interface{})
		case 1:
			temp[key1] = make(map[string]interface{})
			temp[key1] = convert1(value1)
		}
	}
	return temp
}

// parseCMS parses the yaml file and converts it into the CMS format we have
func parseCMS(fileName string, projIndex int) error {
	viper.SetConfigType("yaml")
	// viper.SetConfigName(fileName)
	viper.SetConfigName(fileName)
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
			// we can't use convert1 here because it returns interface{}
			msi := viper.Get(key).(map[string]interface{})
			temp := make(map[string]interface{})
			for key, value := range msi {
				switch value.(type) {
				case []interface{}:
					temp[key] = value.([]interface{})
				case interface{}:
					temp[key] = value.(interface{})
				}
			}

			project.Content.Details[key] = make(map[string]interface{})
			project.Content.Details[key] = temp
		case 2:
			msmsi := viper.Get(key).(map[string]interface{})
			project.Content.Details[key] = make(map[string]interface{})
			project.Content.Details[key] = convert2(msmsi)
		case 3:
			msmsmsi := viper.Get(key).(map[string]interface{})
			project.Content.Details[key] = make(map[string]interface{})
			for key1, value1 := range msmsmsi {
				depth := findInterfaceDepth(value1)
				switch depth {
				case 0:
					project.Content.Details[key][key1] = value1.(interface{})
				case 1:
					project.Content.Details[key][key1] = make(map[string]interface{})
					project.Content.Details[key][key1] = convert1(value1)
				case 2:
					msmsi := value1.(map[string]interface{})
					project.Content.Details[key][key1] = make(map[string]interface{})
					project.Content.Details[key][key1] = convert2(msmsi)
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

func createDemoProjects() []core.Project{
	var demoProjects []core.Project
	var project1, project2, project3, project4 core.Project

	project1.Index = 1
	project1.SeedInvestmentCap = 4000
	project1.Stage = 4
	project1.MoneyRaised = 0
	project1.TotalValue = 4000
	project1.OwnershipShift = 0
	project1.RecipientIndex = -1  // replace with real indices once created
	project1.OriginatorIndex = -1 // replace with real indices once created
	project1.GuarantorIndex = -1  // replace with real indices once created
	project1.ContractorIndex = -1 // replace with real indices once created
	project1.PaybackPeriod = 4    // four weeks payback time
	project1.Chain = "stellar"
	project1.BrokerURL = "mqtt.openx.solar"
	project1.TellerPublishTopic = "opensolartest"
	project1.Metadata = "Aibonito Pilot Project"
	project1.InvestmentType = "munibond"
	project1.TellerURL = ""
	project1.BrokerURL = "https://mqtt.openx.solar"
	project1.TellerPublishTopic = "opensolartest"
	project1.Content.Details = make(map[string]map[string]interface{})

	demoProjects = append(demoProjects, project1)

	project2.Index = 2
	project2.SeedInvestmentCap = 4000
	project2.Stage = 4
	project2.MoneyRaised = 0
	project2.TotalValue = 10000.0
	project2.OwnershipShift = 0
	project2.RecipientIndex = -1  // replace with real indices once created
	project2.OriginatorIndex = -1 // replace with real indices once created
	project2.GuarantorIndex = -1  // replace with real indices once created
	project2.ContractorIndex = -1 // replace with real indices once created
	project2.PaybackPeriod = 4    // four weeks payback time
	project2.Chain = "stellar"
	project2.BrokerURL = "mqtt.openx.solar"
	project2.TellerPublishTopic = "opensolartest"
	project2.Metadata = "Pasto Public School - POC 1 kW"
	project2.InvestmentType = "munibond"
	project2.TellerURL = ""
	project2.BrokerURL = "https://mqtt.openx.solar"
	project2.TellerPublishTopic = "opensolartest"
	project2.Content.Details = make(map[string]map[string]interface{})

	demoProjects = append(demoProjects, project2)

	project3.Index = 3
	project3.SeedInvestmentCap = 15000
	project3.Stage = 4
	project3.MoneyRaised = 0
	project3.TotalValue = 30000
	project3.OwnershipShift = 0
	project3.RecipientIndex = -1  // replace with real indices once created
	project3.OriginatorIndex = -1 // replace with real indices once created
	project3.GuarantorIndex = -1  // replace with real indices once created
	project3.ContractorIndex = -1 // replace with real indices once created
	project3.PaybackPeriod = 4    // four weeks payback time
	project3.Chain = "stellar"
	project3.BrokerURL = "mqtt.openx.solar"
	project3.TellerPublishTopic = "opensolartest"
	project3.Metadata = "New Haven Shelter Solar 2"
	project3.InvestmentType = "munibond"
	project3.TellerURL = ""
	project3.BrokerURL = "https://mqtt.openx.solar"
	project3.TellerPublishTopic = "opensolartest"
	project3.Content.Details = make(map[string]map[string]interface{})

	demoProjects = append(demoProjects, project3)

	project4.Index = 4
	project4.SeedInvestmentCap = 230000
	project4.Stage = 4
	project4.MoneyRaised = 0
	project4.TotalValue = 230000
	project4.OwnershipShift = 0
	project4.RecipientIndex = -1  // replace with real indices once created
	project4.OriginatorIndex = -1 // replace with real indices once created
	project4.GuarantorIndex = -1  // replace with real indices once created
	project4.ContractorIndex = -1 // replace with real indices once created
	project4.PaybackPeriod = 4    // four weeks payback time
	project4.Chain = "stellar"
	project4.BrokerURL = "mqtt.openx.solar"
	project4.TellerPublishTopic = "opensolartest"
	project4.Metadata = "Village Energy Collective, Rural Microgrid"
	project4.InvestmentType = "munibond"
	project4.TellerURL = ""
	project4.BrokerURL = "https://mqtt.openx.solar"
	project4.TellerPublishTopic = "opensolartest"
	project4.Content.Details = make(map[string]map[string]interface{})

	demoProjects = append(demoProjects, project4)

	return demoProjects
}
