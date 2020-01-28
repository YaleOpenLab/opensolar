package main

import (
	"log"
	"time"

	"github.com/Varunram/essentials/xlm/assets"

	"github.com/Varunram/essentials/utils"

	"github.com/Varunram/essentials/xlm/wallet"

	"github.com/Varunram/essentials/xlm"

	"github.com/Varunram/essentials/xlm/stablecoin"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
)

func demoData() error {
	var project core.Project
	project.Name = "5kW Solar at FabIDEAS Coop - Pilot 1"
	project.City = "Aibonito"
	project.State = "Puerto Rico"
	project.Country = "USA"
	project.Location = "Puerto Rico"
	project.DonationType = "Grant"
	project.Originator = "Martin Wainstein"
	project.PanelSize = "50 x 100W"
	project.DailyAvgGeneration = "20 kWh"
	project.InverterSize = "4800W"
	project.Description = "5kW solar to be owned by the FabIDEAS Community Cooperative in Aibonito. The Coop is part of the Instituto Nueva Escuela (INE), and will host a Fab Lab manufacturing montessori school supplies."
	project.Bullet1 = "Powering a unique social entrepreneurship model through a local cooperative"
	project.Bullet2 = "Power purchase agreement tied to the standard local Aibonito tariff"
	project.Bullet3 = "Pay-to-own solar model for the coop with all proceeds reinvested in the Fab Lab"
	project.PilotGoals = append(project.PilotGoals, "Demonstrate contractual automation and disintermediation of renewable energy project finance using blockchain-based smart contracts, as featured in the OpenSolar platform")
	project.PilotGoals = append(project.PilotGoals, "Demonstrate alternative finance schemes with pay-to-own models for community ownership of solar assets.")
	project.PilotGoals = append(project.PilotGoals, "Demonstrate the integration between data from internet-of-things (IoT) devices into payment schemes and climate asset tokenization (Renewable Energy Certificates).")
	project.PilotGoals = append(project.PilotGoals, "Stress test all features in the OpenSolar platforms, including user experiences, fiat on and offramps and smart contracts.")
	project.PilotGoals = append(project.PilotGoals, "Provide a blueprint for a finance plan to transform all of Puerto Rico’s public schools into solar powered emergency shelters.")
	project.Context = `Two years after hurricane Maria hit the island, schools and local communities are still exposed to a centralized and high-carbon energy system vulnerable to climate impacts. The 2020 Earthquake left ⅓ of the island without power. Cooperatives and schools like FabIDEAS and SU Pasto are ideal places for community owned microgrids to be deployed, in order to provide greater power resilience and usher in a new energy economy to Puerto Rico. Since Hurricane Maria, community cooperatives have become nodal points facilitating discussions of concerned parents on how to increase climate & social resilience in the whole community.

The Puerto Rican (PR) government and the department of education are working to appoint schools as emergency shelters —nodes with robust energy and communication systems— for the community to reach out in the event of unavoidable climate shocks. Financing is a key gap. This project acts as a pilot finance mechanism that can help bridge the finance gap to make solar powered schools and community centers more affordable.`
	project.OpportunityDescription = `Cooperativa Fábrica de Ideas de Aibonito (FabIDEAS Coop), is a project that has received the support of Instituto Nueva Escuela (INE), an independent 501c3 nonprofit organization dedicated to transforming the public education system in Puerto Rico through the Montessori philosophy and methodology. FabIDEAS Coop, is an initiative of the community linked to the INE public school S.U. Pasto in the rural town of Aibonito. FabIDEAS Coop aims to create an economic model in which a cooperative of Montessori materials with five initial members serves as an economic engine for the production and distribution of educational products and children furniture, where each additional member of the community that joins the production guild can learn product design and gradually increase his/her income source. It will act as a hub for education in distributed manufacturing to the students of SU Pasto and as an emergency shelter for community members. 

The current project entails the installation of a 5kW system with InverSol’s Lumen battery and inverter unit. Solar will power critical loads in the building, including emergency lights, a telecommunication system, and main manufacturing equipment. The installed system is priced at $12’000, with $9000 being donated by Council Rock / InverSol to support the pilot and coop. The Digital Currency Initiative at the MIT Media Lab will provide a grant of $4000 to cover $3000 of labor cost by the inverSol team and $1000 of other installation services.  The FabIDEAS Coop has agreed to match this $4000 by paying for the solar electricity at the standard utility price (which the building is subject to when it purchases power from the grid) until reaches this amount in cumulative solar power payments. Once these funds accrue they will be reinvested in manufacturing units for the Fab Lab. 

The project is the first full pilot of the Open Solar platform and will test the smart contracts and digital currency enabled by the platform to automate all the dynamics behind the $4000 grant, the solar power payments, and the Renewable Energy Certificates generated by the system. Financial transactions will be automated based on the data read by inverSol’s Lumen unit. 
		`
	project.ArchitectureDesignDescription = "The solar installation will be a behind-the-meter backup setup, to avoid net metering with PREPA’s grid. Future expansion deployments could consider a grid-tied two-way system. The 5kW solar photovoltaics will be installed on the FabIDEAS main building’s roof and connected to the inverSOL’s Lumen unit equipped with a 5kWh battery, a 5kW inverter, a charge regulator and internet-of-things (IoT) functionality."
	project.InstallationArchetype = "This will be a model installation in that the solar and battery support a subpanel of the building circuitry, where only critical loads have been connected. Large manufacturing machinery will not be connected to the subpanel. The system will be configured as a grid-tied installation, in that the main grid can also support other loads in the panel as well as be used to power the battery bank. The installation allows for the interconnection of an emergency generator if needed."
	project.ITInfrastructure = "Main power data readings will come directly from the Lumen all-in-one powermeter unit, transmitting secure data via MQTT protocol. A second revenue-grade meter with IoT pre-pay functionality will be added for further testing integrations. IoT readings from the Lumen system will be used in a smart contract oracle to verify & validate readings for payment and REC generation. A whole building non-invasive powermeter is also contemplated to critical vs. general loads."
	project.HighlightedProduct = `
		inverSOL Lumen:
Lumen by inverSOL is a smart renewable energy system for the home providing greater energy independence and backup power. Lithium NMC (LiNMC ) batteries used in Lumen are validated and produced with uncompromised safety and quality control. Wireless connectivity and computing platform allow for remote control through an app, software upgrades and smart energy management features. 

The Lumen smart features minimize wasted solar power and reduce energy bills, eliminating the need for net metering. The proprietary algorithm built in the Lumen brain ensures solar energy is used even when there is no Sun. Enhanced user experience through an interactive touchscreen and remote control through a mobile app allow to track energy usage and savings. New features available with software updates. Robust and sleek design make Lumen a seamless fit for any interior. Touchscreen and Interactive Design ensure enhanced user experience.
		`
	project.CommunityEngagement.Consultation = "The MIT and Yale team will convene meetings with the FabIDEAS cooperative board to discuss project details and outreach opportunities. The team has already convened a meeting with the Parent-Teacher Organisation of the SU Pasto school, thanks to the coordination of the school’s principal Janice Alejandro, to discuss the role of new finance mechanisms for solar in the local community. Over 50 members of the community gathered to discuss the project, with unanimous approval and significant interest for its replication."
	project.CommunityEngagement.Participation = "The FabIDEAS cooperative community will source volunteers and champions to act as caretakers of the system to monitor its status, report any qualitative information and coordinate with the operation & maintenance required. "
	project.CommunityEngagement.Outreach = "The system will be installed with instructions and visual explanations so that it can act as a pedagogical site for students and community members to learn about the merits of solar energy, electricity and basic electronics. Talks about solar energy will be convened every semester in the context of climate change communication to the community. "
	project.CommunityEngagement.Governance = "The board of the Cooperative and its acting President Maria Pastor will convene bi yearly meeting with the Yale-MIT team (i.e. the originators) to review processes and performance of the solar system and the smart contract."
	project.BusinessNumbers.Description = "The system will be funded by an in-kind donation of inverSOL, providing the solar hardware, and a grant from the Digital Currency Initiative at MIT to cover labor and other service costs. inverSOL’s donation involves the $9000 for 5kW system with  issued by the PR Department of Education covered the principle cost, used for labor and materials. The PPA revenue accrues to pay coupons and mature the bond. The MIT is registered as a first-loss guarantor."
	project.BusinessNumbers.CapitalExpenditure = "The expected capital cost of the project is $13000, including the U$S 9000 product value of a 5kW solar array with a Lumen unity (donated by inverSOL), $3000 of labor costs and $1000 for contingency and other services (covered by the DCI grant)."
	project.BusinessNumbers.ProjectRevenue = "The FabIDEAS cooperative will pay for the solar electricity generated at a standard $/kWh local tariff using an Open Solar platform wallet. Once accumulated payments reach $4000 (stored in the project’s smart contract escrow), these will be released back to the FabIDEAS coop wallet to be used for reinvesting in the fab lab. "
	project.BusinessNumbers.ProjectExpenses = "The project has an O&M (Operation & Management) contingency fund of $1000, but will otherwise will be covered by inverSOL’s guarantee for 5 years. After this period, the cooperative will be responsible for O&M. "
	project.BusinessNumbers.NonProfit = "No net-income or profits will be generated by this project. "
	project.Solar = "5kW"
	project.Battery = "5kWh"
	project.Return = "0%"
	project.Rating = "N/A"
	project.Maturity = "2021"
	project.Storage = "This is a sample storage ipsum"
	project.Tariff = "Unlimited tariff"
	project.Return = "Unlimited Return"
	project.Rating = "AAA"
	project.Tax = "1000"
	project.Acquisition = "Sample Acquisition"

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

	project.RecipientIndex = -1                      // replace with real indices once created
	project.OriginatorIndex = -1                     // replace with real indices once created
	project.GuarantorIndex = -1                      // replace with real indices once created
	project.ContractorIndex = -1                     // replace with real indices once created
	project.PaybackPeriod = consts.FourWeeksInSecond // four weeks payback time
	project.Stage = 4
	project.Chain = "stellar"
	project.OwnershipShift = 0
	project.BrokerUrl = "mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"
	project.Index = 1
	project.TotalValue = 4000
	project.MoneyRaised = 0
	project.Metadata = "Aibonito Pilot Project"
	project.InvestmentType = "munibond"
	project.TellerUrl = ""
	project.BrokerUrl = "https://mqtt.openx.solar"
	project.TellerPublishTopic = "opensolartest"

	err := project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	txhash, err := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, consts.PlatformSeed)
	if err != nil {
		return err
	}

	log.Println("tx for platform trusting stablecoin:", txhash)

	password := "password"
	pwhash := utils.SHA3hash(password)
	seedpwd := "x"
	exchangeAmount := 1.0
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

	recp, err := core.NewRecipient("recp"+run, password, seedpwd, "varunramganesh@gmail.com")
	if err != nil {
		log.Println(err)
		return err
	}

	err = xlm.GetXLM(inv.U.StellarWallet.PublicKey)
	if err != nil {
		log.Println("could not get XLM: ", err)
		return err
	}

	project.RecipientIndex = recp.U.Index
	err = project.Save()
	if err != nil {
		log.Println(err)
		return err
	}

	invSeed, err := wallet.DecryptSeed(inv.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return err
	}

	err = stablecoin.Exchange(inv.U.StellarWallet.PublicKey, invSeed, exchangeAmount)
	if err != nil {
		log.Println("did not exchange xlm", err)
		return err
	}

	time.Sleep(5 * time.Second)

	err = core.Invest(project.Index, inv.U.Index, invAmount, invSeed)
	if err != nil {
		log.Println("did not invest in order", err)
		return err
	}

	time.Sleep(10 * time.Second)

	err = core.UnlockProject(recp.U.Username, pwhash, project.Index, seedpwd)
	if err != nil {
		return err
	}

	return nil
}
