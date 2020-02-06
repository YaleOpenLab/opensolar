package main

import (
	"log"

	core "github.com/YaleOpenLab/opensolar/core"
)

func demoData() error {
	var project core.Project

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
	project.Content.DetailPageStub.Box.ProjectType = "Grant"
	project.Content.DetailPageStub.Box.OriginatorName = "Martin Wainstein"
	project.Content.DetailPageStub.Box.Description = "5kW solar to be owned by the FabIDEAS Community Cooperative in Aibonito. The Coop is part of the Instituto Nueva Escuela (INE), and will host a Fab Lab manufacturing montessori school supplies."
	project.Content.DetailPageStub.Box.Bullet1 = "Powering a unique social entrepreneurship model through a local cooperative"
	project.Content.DetailPageStub.Box.Bullet2 = "Power purchase agreement tied to the standard local Aibonito tariff"
	project.Content.DetailPageStub.Box.Bullet3 = "Pay-to-own solar model for the coop with all proceeds reinvested in the Fab Lab"
	project.Content.DetailPageStub.Box.Solar = "5kW"
	project.Content.DetailPageStub.Box.Battery = "5kWh"
	project.Content.DetailPageStub.Box.Return = "%0"
	project.Content.DetailPageStub.Box.Rating = "N/A"
	project.Content.DetailPageStub.Box.Maturity = project.Acquisition
	project.Content.DetailPageStub.Box.MoneyRaised = project.MoneyRaised
	project.Content.DetailPageStub.Box.TotalValue = project.TotalValue

	project.Content.OtherDetails.Tax = "N/A"
	project.Content.OtherDetails.Storage = "250 Wh"
	project.Content.OtherDetails.Tariff = "0.20$"

	// project.Content.DetailPageStub.Tabs.Terms
	project.Content.DetailPageStub.Tabs.Terms.Purpose = "Proceeds from this project's raise are granted for the development of a pilot solar installation in the FabIDEAS cooperative in Aibonito. The pilot will be used to test the Open Solar platform’s smart contract and financial technology capabilities and is part of a research initiative of the Digital Currency Initiative of the MIT Media Lab and the Yale Open Innovation Lab."
	project.Content.DetailPageStub.Tabs.Terms.Table.Columns = []string{"Variable, Value, Relevant Party", "Note", "Status", "Support Doc"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows = make([][]string, 6)
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[0] = []string{"Investment Type", "Donation", "InverSOL", "Solar Equipment"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[1] = []string{"PPA Tariff", "0.24 Ct/KWh", "Oracle / PREPA", "Variable Anchored To Local Tariff"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[2] = []string{"Return (TEY)", "3.1%", "Broker Dealer", "Tax-Adjusted Yield"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[3] = []string{"Maturity", "+/- 2025", "Broker Dealer", "Variable Tied To Tariff"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[4] = []string{"Guarantee, 50%", "Foundation X", "First-Loss Upon Breach"}
	project.Content.DetailPageStub.Tabs.Terms.Table.Rows[5] = []string{"Insurance", "Premium", "Allianz CS", "Hurricane Coverage"}
	project.Content.DetailPageStub.Tabs.Terms.SecurityNote = "This project does not entail the issuance of a financial security, and is structured exclusively as a restricted research grant. The project does not entail an investment since there are no financial returns offered to donors. All funds accrued through power purchasing are donated back to the Cooperative. Learn more at <link>"

	// project.Content.DetailPageStub.Tabs.Overview
	columns := []string{"Investment", "Financials", "Project Size", "Sustainability Metrics"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.Columns = columns
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData = make(map[string][]string)
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[0]] = []string{"Capex", "$5000"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[0]] = []string{"Hardware", "60%"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[0]] = []string{"First-Loss Escrow", "30%"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[0]] = []string{"Certification Costs", "N/A"}

	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[1]] = []string{"Return (TEY)", "3.1%"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[1]] = []string{"Insurance", "Premium"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[1]] = []string{"Tariff (Variable)", "0.24 ct/kWh"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[1]] = []string{"Maturity (Variable)", "2028"}

	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[2]] = []string{"PV Solar", "1 kW"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[2]] = []string{"Storage", "200 Wh"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[2]] = []string{"% Critical", "2%"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[2]] = []string{"Inverter Capacity", "2.5 kW"}

	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[3]] = []string{"Carbon Drawdown", "0.1t/kWh"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[3]] = []string{"Community Value", "5/7"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[3]] = []string{"LCA", "N/A"}
	project.Content.DetailPageStub.Tabs.Overview.ExecutiveSummary.ColumnData[columns[3]] = []string{"Resilience Rating", "80%"}

	project.Content.DetailPageStub.Tabs.Overview.ImageLink = ""
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Description = `Cooperativa Fábrica de Ideas de Aibonito (FabIDEAS Coop), is a project that has received the support of Instituto Nueva Escuela (INE), an independent 501c3 nonprofit organization dedicated to transforming the public education system in Puerto Rico through the Montessori philosophy and methodology. FabIDEAS Coop, is an initiative of the community linked to the INE public school S.U. Pasto in the rural town of Aibonito. FabIDEAS Coop aims to create an economic model in which a cooperative of Montessori materials with five initial members serves as an economic engine for the production and distribution of educational products and children furniture, where each additional member of the community that joins the production guild can learn product design and gradually increase his/her income source. It will act as a hub for education in distributed manufacturing to the students of SU Pasto and as an emergency shelter for community members. 

The current project entails the installation of a 5kW system with InverSol’s Lumen battery and inverter unit. Solar will power critical loads in the building, including emergency lights, a telecommunication system, and main manufacturing equipment. The installed system is priced at $12’000, with $9000 being donated by Council Rock / InverSol to support the pilot and coop. The Digital Currency Initiative at the MIT Media Lab will provide a grant of $4000 to cover $3000 of labor cost by the inverSol team and $1000 of other installation services.  The FabIDEAS Coop has agreed to match this $4000 by paying for the solar electricity at the standard utility price (which the building is subject to when it purchases power from the grid) until reaches this amount in cumulative solar power payments. Once these funds accrue they will be reinvested in manufacturing units for the Fab Lab. 

The project is the first full pilot of the Open Solar platform and will test the smart contracts and digital currency enabled by the platform to automate all the dynamics behind the $4000 grant, the solar power payments, and the Renewable Energy Certificates generated by the system. Financial transactions will be automated based on the data read by inverSol’s Lumen unit. 
		`
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals = make([]string, 5)
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[0] = "Demonstrate contractual automation and disintermediation of renewable energy project finance using blockchain-based smart contracts, as featured in the OpenSolar platform"
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[1] = "Demonstrate alternative finance schemes with pay-to-own models for community ownership of solar assets."
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[2] = "Demonstrate the integration between data from internet-of-things (IoT) devices into payment schemes and climate asset tokenization (Renewable Energy Certificates)."
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[3] = "Stress test all features in the OpenSolar platforms, including user experiences, fiat on and offramps and smart contracts."
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.PilotGoals[4] = "Provide a blueprint for a finance plan to transform all of Puerto Rico’s public schools into solar powered emergency shelters."

	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images = make([]string, 3)
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images[0] = ""
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images[1] = ""
	project.Content.DetailPageStub.Tabs.Overview.Opportunity.Images[2] = ""

	project.Content.DetailPageStub.Tabs.Overview.Context = `Two years after hurricane Maria hit the island, schools and local communities are still exposed to a centralized and high-carbon energy system vulnerable to climate impacts. The 2020 Earthquake left ⅓ of the island without power. Cooperatives and schools like FabIDEAS and SU Pasto are ideal places for community owned microgrids to be deployed, in order to provide greater power resilience and usher in a new energy economy to Puerto Rico. Since Hurricane Maria, community cooperatives have become nodal points facilitating discussions of concerned parents on how to increase climate & social resilience in the whole community.

The Puerto Rican (PR) government and the department of education are working to appoint schools as emergency shelters —nodes with robust energy and communication systems— for the community to reach out in the event of unavoidable climate shocks. Financing is a key gap. This project acts as a pilot finance mechanism that can help bridge the finance gap to make solar powered schools and community centers more affordable.`

	// project.Content.DetailPageStub.Tabs.Project
	project.Content.DetailPageStub.Tabs.Project.Architecture.MapLayoutImage = ""
	project.Content.DetailPageStub.Tabs.Project.Architecture.SolarOutputImage = ""
	project.Content.DetailPageStub.Tabs.Project.Architecture.DesignDescription = "The solar installation will be a behind-the-meter backup setup, to avoid net metering with PREPA’s grid. Future expansion deployments could consider a grid-tied two-way system. The 5kW solar photovoltaics will be installed on the FabIDEAS main building’s roof and connected to the inverSOL’s Lumen unit equipped with a 5kWh battery, a 5kW inverter, a charge regulator and internet-of-things (IoT) functionality."
	project.Content.DetailPageStub.Tabs.Project.Architecture.Description = "This 10 kW system was installed in a Grid-Tied design by CT solar Developers. It has a 15 kW smart solarEdge inverter, and an Itron revenue grade REC meter. The system is roof-mounted on a 55 degree angle on a SE facing view. However, its high efficient Jingko panels provide a 78% efficiency rating."

	project.Content.DetailPageStub.Tabs.Project.Layout.InstallationArchetype = "This will be a model installation in that the solar and battery support a subpanel of the building circuitry, where only critical loads have been connected. Large manufacturing machinery will not be connected to the subpanel. The system will be configured as a grid-tied installation, in that the main grid can also support other loads in the panel as well as be used to power the battery bank. The installation allows for the interconnection of an emergency generator if needed."
	project.Content.DetailPageStub.Tabs.Project.Layout.ITInfrastructure = "Main power data readings will come directly from the Lumen all-in-one powermeter unit, transmitting secure data via MQTT protocol. A second revenue-grade meter with IoT pre-pay functionality will be added for further testing integrations. IoT readings from the Lumen system will be used in a smart contract oracle to verify & validate readings for payment and REC generation. A whole building non-invasive powermeter is also contemplated to critical vs. general loads."
	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Description = `Lumen by inverSOL is a smart renewable energy system for the home providing greater energy independence and backup power. Lithium NMC (LiNMC ) batteries used in Lumen are validated and produced with uncompromised safety and quality control. Wireless connectivity and computing platform allow for remote control through an app, software upgrades and smart energy management features.

The Lumen smart features minimize wasted solar power and reduce energy bills, eliminating the need for net metering. The proprietary algorithm built in the Lumen brain ensures solar energy is used even when there is no Sun. Enhanced user experience through an interactive touchscreen and remote control through a mobile app allow to track energy usage and savings. New features available with software updates. Robust and sleek design make Lumen a seamless fit for any interior. Touchscreen and Interactive Design ensure enhanced user experience.`
	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Images = make([]string, 2)
	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Images[0] = ""
	project.Content.DetailPageStub.Tabs.Project.Layout.HighlightedProduct.Images[1] = ""

	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.Columns = []string{"Consultation", "Participation", "Outreach", "Governance"}
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.ColumnData = make([][]string, len(project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.Columns))
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.ColumnData[0] = []string{"", "The MIT and Yale team will convene meetings with the FabIDEAS cooperative board to discuss project details and outreach opportunities. The team has already convened a meeting with the Parent-Teacher Organisation of the SU Pasto school, thanks to the coordination of the school’s principal Janice Alejandro, to discuss the role of new finance mechanisms for solar in the local community. Over 50 members of the community gathered to discuss the project, with unanimous approval and significant interest for its replication."}
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.ColumnData[1] = []string{"The FabIDEAS cooperative community will source volunteers and champions to act as caretakers of the system to monitor its status, report any qualitative information and coordinate with the operation & maintenance required."}
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.ColumnData[2] = []string{"The system will be installed with instructions and visual explanations so that it can act as a pedagogical site for students and community members to learn about the merits of solar energy, electricity and basic electronics. Talks about solar energy will be convened every semester in the context of climate change communication to the community."}
	project.Content.DetailPageStub.Tabs.Project.CommunityEngagement.ColumnData[3] = []string{"The board of the Cooperative and its acting President Maria Pastor will convene bi yearly meeting with the Yale-MIT team (i.e. the originators) to review processes and performance of the solar system and the smart contract."}

	project.Content.DetailPageStub.Tabs.Project.BizNumbers.Description = "The system will be funded by an in-kind donation of inverSOL, providing the solar hardware, and a grant from the Digital Currency Initiative at MIT to cover labor and other service costs. inverSOL’s donation involves the $9000 for 5kW system with issued by the PR Department of Education covered the principle cost, used for labor and materials. The PPA revenue accrues to pay coupons and mature the bond."
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.GeneralPaymentLogic = ""
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.CapitalExpenditure = "The expected capital cost of the project is $13000, including the U$S 9000 product value of a 5kW solar array with a Lumen unity (donated by inverSOL), $3000 of labor costs and $1000 for contingency and other services (covered by the DCI grant)."
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.CapitalExpenditureImage = ""
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.ProjectRevenue = "The FabIDEAS cooperative will pay for the solar electricity generated at a standard $/kWh local tariff using an Open Solar platform wallet. Once accumulated payments reach $4000 (stored in the project’s smart contract escrow), these will be released back to the FabIDEAS coop wallet to be used for reinvesting in the fab lab."
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.ProjectExpenses = "The project has an O&M (Operation & Management) contingency fund of $1000, but will otherwise will be covered by inverSOL’s guarantee for 5 years. After this period, the cooperative will be responsible for O&M."
	project.Content.DetailPageStub.Tabs.Project.BizNumbers.NonProfit = "No net-income or profits will be generated by this project."

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

	err := project.Save()
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
