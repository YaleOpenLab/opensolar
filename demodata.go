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
	// project.Content.DetailPageStub.Box
	project.Content.DetailPageStub.Box.Name = project.Name
	project.Content.DetailPageStub.Box.Location = project.City + ", " + project.State + ", " + project.Country
	project.Content.DetailPageStub.Box.Maturity = project.Acquisition
	project.Content.DetailPageStub.Box.MoneyRaised = project.MoneyRaised
	project.Content.DetailPageStub.Box.TotalValue = project.TotalValue

	// project.Content.DetailPageStub.Tabs.Overview
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.Columns = make(map[string]map[string]string)

	// project.Content.DetailPageStub.Tabs.StageForecast
	project.Content.DetailPageStub.Tabs.StageForecast.DevelopmentStage.Image = ""
	project.Content.DetailPageStub.Tabs.StageForecast.DevelopmentStage.StageTitle = "Construction"
	project.Content.DetailPageStub.Tabs.StageForecast.DevelopmentStage.StageDescription = "The project is in the contract development and signing stage. In this stage, the power purchase agreement and general financial variables behind the Open Solar platform’s smart contract are carefully negotiated, drafted and signed by all relevant parties. Full funding of the project is not available."
	project.Content.DetailPageStub.Tabs.StageForecast.DevelopmentStage.OtherLink = ""

	// project.Content.DetailPageStub.Tabs.StageForecast.SolarStage
	// project.Content.DetailPageStub.Tabs.Documents
	project.Content.DetailPageStub.Tabs.Documents.Description = ""
	project.Content.DetailPageStub.Tabs.Documents.LegalContracts.Image = ""
	project.Content.DetailPageStub.Tabs.Documents.LegalContracts.Title = ""
	project.Content.DetailPageStub.Tabs.Documents.LegalContracts.Description = ""
	project.Content.DetailPageStub.Tabs.Documents.SmartContractsImage = ""
	project.Content.DetailPageStub.Tabs.Documents.SCReviewDescription = ""

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

// parseCMS parses the yaml file and converts it into the CMS format we have
func parseCMS(fileName string, projIndex int) error {
	viper.SetConfigType("yaml")
	// viper.SetConfigName(fileName)
	viper.SetConfigName("cms")
	// viper.AddConfigPath("./data")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	project, err := core.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	etm := viper.Get("Explore Tab Modal").(map[string]interface{})
	project.Content.DetailPageStub.Box.ProjectType = etm["projecttype"].(string)
	project.Content.DetailPageStub.Box.OriginatorName = etm["originatorname"].(string)
	project.Content.DetailPageStub.Box.Description = etm["description"].(string)
	project.Content.DetailPageStub.Box.Bullet1 = etm["bullet1"].(string)
	project.Content.DetailPageStub.Box.Bullet2 = etm["bullet2"].(string)
	project.Content.DetailPageStub.Box.Bullet3 = etm["bullet3"].(string)
	project.Content.DetailPageStub.Box.Solar = etm["solar"].(string)
	project.Content.DetailPageStub.Box.Battery = etm["battery"].(string)
	project.Content.DetailPageStub.Box.Return = etm["return"].(string)
	project.Content.DetailPageStub.Box.Rating = etm["rating"].(string)

	od := viper.Get("Other Details").(map[string]interface{})
	project.Content.OtherDetails.Tax = od["tax"].(string)
	project.Content.OtherDetails.Storage = od["storage"].(string)
	project.Content.OtherDetails.Tariff = od["tariff"].(string)

	terms := viper.Get("Terms").(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Terms.Purpose = terms["purpose"].(string)

	table := terms["table"].(map[string]interface{})

	tableColumns := table["columns"].([]interface{})
	for _, column := range tableColumns {
		project.Content.DetailPageStub.Tabs.Terms.Table.Columns = append(project.Content.DetailPageStub.Tabs.Terms.Table.Columns, column.(string))
	}

	tableRows := table["rows"].([]interface{})
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows = make([][]string, len(tableRows))
	for i, xrow := range tableRows {
		row := xrow.([]interface{})
		for _, subrow := range row {
			project.Content.DetailPageStub.Tabs.Terms.Table.Rows[i] = append(project.Content.DetailPageStub.Tabs.Terms.Table.Rows[i], subrow.(string))
		}
	}

	project.Content.DetailPageStub.Tabs.Terms.SecurityNote = terms["securitynote"].(string)

	overview := viper.Get("overview").(map[string]interface{})
	log.Println(overview)
	execSummary := overview["executive summary"].(map[string]interface{})

	for execKeys, execVals := range execSummary {
		exec := execVals.(map[string]interface{})
		project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.Columns[execKeys] = make(map[string]string)
		for key, value := range exec {
			project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.Columns[execKeys][key] = value.(string)
		}
	}

	project.Content.DetailPageStub.Tabs.Overview.ImageLink = overview["imagelink"].(string)

	opportunity := overview["opportunity"].(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Description = opportunity["description"].(string)

	pilotgoals := opportunity["pilotgoals"].([]interface{})
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals = make([]string, len(pilotgoals))
	for i, goal := range pilotgoals {
		project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[i] = goal.(string)
	}

	oimages := opportunity["images"].([]interface{})
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images = make([]string, len(oimages))
	for i, image := range oimages {
		project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images[i] = image.(string)
	}

	project.Content.DetailPageStub.Tabs.Overview.Context = opportunity["context"].(string)

	projDetails := viper.Get("project details").(map[string]interface{})

	projArch := projDetails["architecture"].(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Project.Architecture.MapLayoutImage = projArch["maplayoutimage"].(string)
	project.Content.DetailPageStub.Tabs.Project.Architecture.SolarOutputImage = projArch["solaroutputimage"].(string)
	project.Content.DetailPageStub.Tabs.Project.Architecture.DesignDescription = projArch["designdescription"].(string)
	project.Content.DetailPageStub.Tabs.Project.Architecture.Description = projArch["description"].(string)

	projLayout := projDetails["layout"].(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Project.Layout.InstallationArchetype = projLayout["installationarchetype"].(string)
	project.Content.DetailPageStub.Tabs.Project.Layout.ITInfrastructure = projLayout["itinfrastructure"].(string)

	hProduct := projLayout["highlightedproduct"].(map[string]interface{})

	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Description = hProduct["description"].(string)

	hpImages := hProduct["images"].([]interface{})
	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Images = make([]string, len(hpImages))

	for i, image := range hpImages {
		project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Images[i] = image.(string)
	}

	project.Content.DetailPageStub.Tabs.Project.Layout.Description = projLayout["description"].(string)

	comEng := projDetails["community engagement"].(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.Columns = make(map[string][]string, len(comEng))

	for cKeys, cVals := range comEng {
		log.Println(cVals)
		arr := cVals.([]interface{})
		var columns []string
		for _, strings := range arr {
			columns = append(columns, strings.(string))
		}
		project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.Columns[cKeys] = columns
	}

	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.Description = comEng["description"].([]interface{})[0].(string)

	bizNumbers := projDetails["biznumbers"].(map[string]interface{})
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.Description = bizNumbers["description"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.GeneralPaymentLogic = bizNumbers["generalpaymentlogic"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.CapitalExpenditure = bizNumbers["capitalexpenditure"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.CapitalExpenditureImage = bizNumbers["capitalexpenditureimage"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.ProjectRevenue = bizNumbers["projectrevenue"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.ProjectExpenses = bizNumbers["projectexpenses"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.NonProfit = bizNumbers["nonprofit"].(string)
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.OtherLinks = bizNumbers["otherlinks"].(string)

	err = project.Save()
	if err != nil {
		return err
	}

	return nil
}
