package opensolar

import (
	"github.com/pkg/errors"
	"log"
	"time"

	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	escrow "github.com/YaleOpenLab/openx/chains/xlm/escrow"
	issuer "github.com/YaleOpenLab/openx/chains/xlm/issuer"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	notif "github.com/YaleOpenLab/openx/notif"
	oracle "github.com/YaleOpenLab/openx/oracle"
)

// platform is designed to be monolithic so we can have everything in one place

const (
	InvestorWeight         = 0.1 // the percentage weight of the project's total reputation assigned to the investor
	OriginatorWeight       = 0.1 // the percentage weight of the project's total reputation assigned to the originator
	ContractorWeight       = 0.3 // the percentage weight of the project's total reputation assigned to the contractor
	DeveloperWeight        = 0.2 // the percentage weight of the project's total reputation assigned to the developer
	RecipientWeight        = 0.3 // the percentage weight of the project's total reputation assigned to the recipient
	NormalThreshold        = 1   // NormalThreshold is the normal payback interval of 1 payback period. Regular notifications are sent regardless of whether the user has paid back towards the project.
	AlertThreshold         = 2   // AlertThreshold is the threshold above which the user gets a nice email requesting a quick payback whenever possible
	SternAlertThreshold    = 4   // SternAlertThreshold is the threshold above when the user gets a warning that services will be disconnected if the user doesn't payback soon.
	DisconnectionThreshold = 6   // DisconnectionThreshold is the threshold above which the user gets a notification telling that services have been disconnected.
)

// SetStage sets the stage of a project
func (a *Project) SetStage(number int) error {
	switch number {
	case 3:
		a.Reputation = a.TotalValue // upgrade reputation since totalValue might have changed from the originated contract
		err := a.Save()
		if err != nil {
			log.Println("Error while saving project", err)
			return err
		}
		err = RepOriginatedProject(a.OriginatorIndex, a.Index) // modify originator reputation now that the final price is fixed
		if err != nil {
			log.Println("Error while increasing reputation", err)
			return err
		}
	case 5:
		contractor, err := RetrieveEntity(a.ContractorIndex)
		if err != nil {
			log.Println("error while retrieving entity from db, quitting")
			return err
		}
		err = contractor.U.ChangeReputation(a.TotalValue * ContractorWeight) // modify contractor Reputation now that a project has been installed
		if err != nil {
			log.Println("Couldn't increase contractor reputation", err)
			return err
		}

		for _, i := range a.InvestorIndices {
			elem, err := RetrieveInvestor(i)
			if err != nil {
				log.Println("Error while retrieving investor", err)
				return err
			}
			err = elem.U.ChangeReputation(a.TotalValue * InvestorWeight)
			if err != nil {
				log.Println("Couldn't change investor reputation", err)
				return err
			}
		}
	case 6:
		recp, err := RetrieveRecipient(a.RecipientIndex)
		if err != nil {
			return err
		}
		err = recp.U.ChangeReputation(a.TotalValue * RecipientWeight) // modify recipient reputation now that the system had begun power generation
		if err != nil {
			log.Println("Error while changing recipient reputation", err)
			return err
		}
	default:
		log.Println("default")
	}
	a.Stage = number
	return a.Save()
}

var Stage0 = Stage{
	Number:       0,
	FriendlyName: "Handshake",
	Name:         "Idea Consolidation",
	Activities: []string{
		"[Originator] proposes project and either secures or agrees to serve as [Solar Developer]. NOTE: Originator is the community leader or catalyst for the project, they may opt to serve as the solar developer themselves, or pass that responsibility off, going forward we will use solar developer to represent the interest of both.",
		"[Solar Developer] creates general estimation of project (eg. with an automatic calculation through Google Project Sunroof, PV) ",
		"If [Originator]/[Solar Developer] is not landowner [Host] states legal ownership of site (hard proof is optional at this stage)",
	},
	StateTrigger: []string{
		"Matching of originator with receiver, and mutual approval/intention of interest.",
	},
}

// StageXtoY promtoes a contract from  stage X.Number to stage Y.Number
func StageXtoY(index int) error {
	// check for out of bound errors
	// retrieve the project
	project, err := RetrieveProject(index)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	if project.Stage < 0 || project.Stage > 8 {
		log.Println("project stage number out of bounds, quitting")
		return errors.New("stage number out of bounds or not eligible for stage updation")
	}

	if project.StageChecklist == nil || project.StageData == nil {
		log.Println("stage checklist or stage data is nil, quitting")
		return errors.New("stage checklist or stage data is nil, quitting")
	}

	var baseStage Stage
	var finalStage Stage

	switch project.Stage {
	case 0:
		baseStage = Stage0
		finalStage = Stage1
	case 1:
		baseStage = Stage1
		finalStage = Stage2
	case 2:
		baseStage = Stage2
		finalStage = Stage3
	case 3:
		baseStage = Stage3
		finalStage = Stage4
	case 4:
		baseStage = Stage4
		finalStage = Stage5
	case 5:
		baseStage = Stage5
		finalStage = Stage6
	case 6:
		baseStage = Stage6
		finalStage = Stage7
	case 7:
		baseStage = Stage7
		finalStage = Stage8
	case 8:
		baseStage = Stage8
		finalStage = Stage9
	default:
		// shouldn't come here? in case it does, error out.
		return errors.New("base stage doesn't match with predefined stages, quitting")
	}

	if len(project.StageChecklist[baseStage.Number]) != len(baseStage.Activities) {
		log.Println("length of checklists don't match, quitting")
		return errors.New("length of checklists don't match, quitting")
	}

	if len(project.StageData[baseStage.Number]) == 0 {
		log.Println("baseStage data is empty, can't upgrade stages!")
		return errors.New("baseStage data is empty, can't upgrade stages")
	}

	// go through the checklist and see if something's wrong
	for _, check := range project.StageChecklist[baseStage.Number] {
		if !check {
			log.Println("checklist not satisfied, quitting")
			return errors.New("checklist not satisfied, quitting")
		}
	}

	// everything in the checklist is set to true, so we can upgrade from stage 0 to 1 safely
	log.Println("Upgrading: ", project.Index, " from stage: ", baseStage.Number, " to stage: ", finalStage.Number)
	return project.SetStage(finalStage.Number)
}

var Stage1 = Stage{
	Number:       1,
	FriendlyName: "Engagement",
	Name:         "RFP Development",
	Activities: []string{
		"[Solar Developer] Analyse parameters, create financial model (proforma)",
		"[Host] & [Solar Developer] engage [Legal] & begin scoping site for planning constraints and opportunities (viability analysis)",
		"[Solar Developer] Create RFP (‘Request For Proposal’)",
		"Simple: Automatic calculation (eg. Sunroof style)",
		"Complex: Public project with 3rd party RFP consultant (independent engineer)",
		"[Originator][Solar Developer][Offtaker] Post project for RFP",
		"[Beneficiary/Host] Define and select RFP developer.",
		"[Investor] First angel investment option (high risk)",
		"Allow ‘time banking’ as sweat equity, monetized as tokenized capital or shadow stock",
	},
	StateTrigger: []string{
		"Issue an RFP",
		"Letter of Intent or MOU between originator and developer",
	},
}

var Stage2 = Stage{
	Number:       2,
	FriendlyName: "Quotes",
	Name:         "Actions",
	Activities: []string{
		"[Solar Developer][Beneficiary/Offtaker][Legal] PPA model negotiation.",
		"[Originator][Beneficiary]  Compare quotes from bidders: ",
		"[Engineering Procurement and Construction] (labor)",
		"[Vendors] (Hardware)",
		"[Insurers]",
		"[Issuer]",
		"[Intermediary Portal]",
		"[Originator/Receiver] Begin negotiation with [Utility]",
		"[Solar Developer] checks whether site upgrades are necessary.",
		"[Solar Developer][Host] Prepare submission for permitting and planning",
		"[Investor] Angel incorporation (less risk)",
	},
	StateTrigger: []string{
		"Selection of quotes and vendors",
		"Necessary identification of entities: Installers and offtaker",
	},
}

var Stage3 = Stage{
	Number:       3,
	FriendlyName: "Signing",
	Name:         "Contract Execution",
	Activities: []string{
		"[Solar Developer] pays [Legal] for PPA finalization.",
		"[Solar Developer][Host] Signs site Lease with landowner.",
		"[Solar Developer] OR [Issuer] signs Offering Agreement with [Intermediary Portal].",
		"[Solar Developer][Beneficiary] selects and signs contracts with: ",
		"[Engineering Procurement and Construction] (labor)",
		"[Vendors] (Hardware)",
		"[Insurers]",
		"[Issuer] OR [Intermediary Portal]",
		"[Offtaker] OR [Solar Developer][Engineering, Procurement and Construction] sign vendor/developer EPC Contracts",
		"[Solar Developer][Offtaker] signs PPA/Offtake Agreement",
		"[Investor] 2nd stage of eligible funding",
		"[Solar Developer][Beneficiary] makes downpayment to [Engineering Procurement and Construction] (labor)",
		"[Investor] Profile with risk ",
	},
	StateTrigger: []string{
		"Execution of contracts - Sign!",
	},
}

var Stage4 = Stage{
	Number:       4,
	FriendlyName: "The Raise",
	Name:         "Finance and Capitalization",
	Activities: []string{
		"[Issuer] engages [Intermediary Portal] to develop Form C or prospectus",
		"[Intermediary Portal] lists [Issuer] project",
		"[Originator][Solar Developer][Offtaker] market the crowdfunded offering",
		"[Investors] Commit capital to the project",
		"[Intermediary Portal] closes offering and disburses capital from Escrow account to [Issuers]",
		"If [Issuer] is not also [Solar Developer] then [Issuer] passes funds to [Solar Developer] ",
	},
	StateTrigger: []string{
		"Project account receives funds that cover the raise amount. Raise amount: normally includes both project capital expenditure (i.e. hardware and labor) and ongoing Operation & Management costs",
	},
}

var Stage5 = Stage{
	Number:       5,
	FriendlyName: "Construction",
	Name:         "Payments and Construction",
	Activities: []string{
		"[Solar Developer] coordinates installation dates and arrangements with [Host][Off-takers]",
		"[Solar Developer] OR [Engineering, Procurement and Construction] take delivery of equipment from [Vendor]",
		"[Utility] issues conditional interconnection",
		"[Solar Developer] schedules installation with [Engineering, Procurement and Construction]",
		"[Engineering, Procurement and Construction] completes installation.",
		"[Solar Developer] pays [Engineering, Procurement and Construction] for substantial completion of the project.",
		"[Insurers] verifies policy, [Solar Developer] pays [Insurers]",
		"[Investor] role?",
	},
	StateTrigger: []string{
		"Installation reaches substantial completion",
		"IoT devices detect energy generation",
	},
}

var Stage6 = Stage{
	Number:       6,
	FriendlyName: "Interconnection",
	Name:         "Contract Execution",
	Activities: []string{
		"[Solar Developer] coordinates with [Engineering Procurement and Construction] to schedule interconnection dates with [Utility] ",
		"[Engineering, Procurement and Construction] submits ‘as-built’ drawings to City/County Inspectors and schedules interconnection with [Utility]",
		"[Solar Developer] schedules City/County Building Inspector visit",
		"[Utility] visits site for witness test",
		"[Utility] places project in service ",
	},
	StateTrigger: []string{
		"[Utility] places project in service",
	},
}

var Stage7 = Stage{
	Number:       7,
	FriendlyName: "Legacy",
	Name:         "Operation and Management",
	Activities: []string{
		"[Solar Developer] hires OR becomes [Manager]",
		"[Manager] hires [Operations & Maintenance] provider",
		"[Manager] sets up billing system and issues monthly bills to [Offtaker] and collects payment on bills",
		"[Manager] monitors for breaches of payment or contract, other indentures, force majeure or adverse conditions [see below for Breach Conditions]",
		"[Manager] files annual taxes",
		"[Manager] handles annual true-up on net-metering payments",
		"[Manager] makes annual cash distributions and issues 1099-DIV to [Investors] or coordinates share repurchase from [Investors]",
		"If applicable, [Manager] executes flip between [Solar Developer] ownership interest and [Tax equity investor]",
		"[Manager] OR [Operations & Maintenance] monitors system performance and coordinates with [Off-takers] to schedule routine maintenance",
		"[Manager] OR [Operations & Maintenance] coordinates with [Engineering, Procurement and Construction] to change inverters or purchase replacements from [Vendors] as needed.",
		"[Investors] can engage in secondary market (i.e. re-selling its securities). ",
	},
	StateTrigger: []string{
		"[Investors] reach preferred return rate, or Power Purchase Agreement stipulates ownership flip date or conditions ",
	},
	BreachCondition: []string{
		"[Offtaker] fails to make $/kWh payments after X period of time due. ",
	},
}

var Stage8 = Stage{
	Number:       8,
	FriendlyName: "Handoff",
	Name:         "Ownership Flip",
	Activities: []string{
		"[Beneficiary/Offtakers] Payments accrue to cover the [Investor] principle (i.e. total raised amount)",
		"Escrow account (eg. capital account) pays off principle to [Investor]",
	},
	StateTrigger: []string{
		"[Beneficiary] (eg. Host, Holding)  becomes full legal owner of physical assets",
		"[Investors] exit the project",
	},
}

var Stage9 = Stage{
	Number:       9,
	FriendlyName: "End of Life",
	Name:         "Disposal",
	Activities: []string{
		"[IoT] Solar equipment is generating below a productivity threshold, or shows general malfunction",
		"[Beneficiaries][Developers] dispose of the equipment to a recycling program",
		"[Developer/Recycler] Certifies equipment is received",
	},
	StateTrigger: []string{
		"Project termination",
		"Wallet terminations",
	},
}

// VerifyBeforeAuthorizing verifies some information on the originator before upgrading the project stage
func VerifyBeforeAuthorizing(projIndex int) bool {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return false
	}
	originator, err := RetrieveEntity(project.OriginatorIndex)
	if err != nil {
		return false
	}
	log.Println("ORIGINATOR'S NAME IS:", originator.U.Name, " and PROJECT's METADATA IS: ", project.Metadata)
	if originator.U.Kyc && !originator.U.Banned {
		return true
	}
	return false
}

// RecipientAuthorize allows a recipient to authorize a specific project
func RecipientAuthorize(projIndex int, recpIndex int) error {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	if project.Stage != 0 {
		return errors.New("Project stage not zero")
	}
	if !VerifyBeforeAuthorizing(projIndex) {
		return errors.New("Originator not verified")
	}
	recipient, err := RetrieveRecipient(recpIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient")
	}
	if project.RecipientIndex != recipient.U.Index {
		return errors.New("You can't authorize a project which is not assigned to you!")
	}

	err = project.SetStage(1) // set the project as originated
	if err != nil {
		return errors.Wrap(err, "Error while setting origin project")
	}

	err = RepOriginatedProject(project.OriginatorIndex, project.Index)
	if err != nil {
		return errors.Wrap(err, "error while increasing reputation of originator")
	}

	return nil
}

// —VOTING SCHEMES—
// MW: Lets design this together. Very cool to have votes (which are 'Likes'), but why only investors can vote? Why not projects at stage 1?
// What does it mean if a project has high votes?

// VoteTowardsProposedProject is a handler that an investor would use to vote towards a
// specific proposed project on the platform.
func VoteTowardsProposedProject(invIndex int, votes float64, projectIndex int) error {
	inv, err := RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve investor")
	}
	if votes > inv.VotingBalance {
		return errors.New("Can't vote with an amount greater than available balance")
	}

	project, err := RetrieveProject(projectIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	if project.Stage != 2 {
		return errors.New("You can't vote for a project with stage not equal to 2")
	}

	project.Votes += votes
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}

	err = inv.ChangeVotingBalance(votes)
	if err != nil {
		return errors.Wrap(err, "error while deducitng voting balance of investor")
	}

	log.Println("CAST VOTE TOWARDS PROJECT SUCCESSFULLY")
	return nil
}

// preInvestmentChecks associated with the opensolar platform when an Investor bids an investment amount of a specific project
func preInvestmentCheck(projIndex int, invIndex int, invAmount float64, seed string) (Project, error) {
	var project Project
	var investor Investor
	var err error

	project, err = RetrieveProject(projIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve project")
	}

	investor, err = RetrieveInvestor(invIndex)
	if err != nil {
		return project, errors.Wrap(err, "couldn't retrieve investor")
	}

	if !investor.CanInvest(invAmount) {
		return project, errors.New("Investor has less balance than what is required to invest in this project")
	}

	pubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return project, errors.Wrap(err, "could not get pubkey from seed")
	}

	if !xlm.AccountExists(pubkey) {
		return project, errors.New("accoutn doesn't exist yet, quitting")
	}
	// check if investment amount is greater than or equal to the project requirements
	if invAmount > project.TotalValue-project.MoneyRaised {
		return project, errors.New("Investment amount greater than what is required! Adjust your investment")
	}

	// the checks till here are common for all chains. The stuff following this is exclusive to stellar.
	if project.Chain == "stellar" || project.Chain == "" {
		if project.SeedAssetCode == "" && project.InvestorAssetCode == "" {
			// this project does not have an asset issuer associated with it yet since there has been
			// no seed round nor investment round
			project.InvestorAssetCode = assets.AssetID(consts.InvestorAssetPrefix + project.Metadata) // you can retrieve asetCodes anywhere since metadata is assumed to be unique
			err = project.Save()
			if err != nil {
				return project, errors.Wrap(err, "couldn't save project")
			}
			err = issuer.InitIssuer(consts.OpenSolarIssuerDir, projIndex, consts.IssuerSeedPwd)
			if err != nil {
				return project, errors.Wrap(err, "error while initializing issuer")
			}
			err = issuer.FundIssuer(consts.OpenSolarIssuerDir, projIndex, consts.IssuerSeedPwd, consts.PlatformSeed)
			if err != nil {
				return project, errors.Wrap(err, "error while funding issuer")
			}
		}

		return project, nil
	} else if project.Chain == "algorand" {
		return project, errors.Wrap(err, "algorand investments not supported yet, quitting")
	} else {
		return project, errors.Wrap(err, "chain not supported, quitting")
	}
}

// SeedInvest is the seed investment function of the opensolar platform
func SeedInvest(projIndex int, invIndex int, invAmount float64, invSeed string) error {

	project, err := preInvestmentCheck(projIndex, invIndex, invAmount, invSeed)
	if err != nil {
		return errors.Wrap(err, "error while performing pre investment check")
	}

	if project.Stage != 1 && project.Stage != 2 {
		return errors.New("project stage not at 1, you either have passed the seed stage or project is not at seed stage yet")
	}

	if project.InvestmentType != "munibond" {
		return errors.New("investment models other than munibonds are not supported right now, quitting")
	}

	if project.SeedInvestmentCap < invAmount {
		return errors.New("you can't invest more than what the seed investment cap permits you to, quitting")
	}

	if project.Chain == "stellar" || project.Chain == "" {
		if project.SeedAssetCode == "" {
			log.Println("assigning a seed asset code")
			project.SeedAssetCode = "SEEDASSET"
		}
		err = MunibondInvest(consts.OpenSolarIssuerDir, invIndex, invSeed, invAmount, projIndex,
			project.SeedAssetCode, project.TotalValue, project.SeedInvestmentFactor)
		if err != nil {
			return errors.Wrap(err, "error while investing")
		}

		err = project.updateAfterInvestment(invAmount, invIndex)
		if err != nil {
			return errors.Wrap(err, "couldn't update project after investment")
		}

		return err
	} else {
		return errors.New("other chain investments not supported  yet")
	}
}

// Invest is the main invest function of the opensolar platform
func Invest(projIndex int, invIndex int, invAmount float64, invSeed string) error {
	var err error

	// run preinvestment checks to make sure everything is okay
	project, err := preInvestmentCheck(projIndex, invIndex, invAmount, invSeed)
	if err != nil {
		return errors.Wrap(err, "pre investment check failed")
	}

	if project.InvestmentType != "munibond" {
		return errors.New("other investment models are not supported right now, quitting")
	}

	if project.Chain == "stellar" || project.Chain == "" {
		if project.Stage != 4 {
			// if the project is not at stage 4 due to some reason, catch it here
			if project.Stage == 1 || project.Stage == 2 {
				// need to redirect it to the seedinvest function
				return SeedInvest(projIndex, invIndex, invAmount, invSeed)
			}
			return errors.New("project not at stage where it can solicit investment, quitting")
		}
		// call the model and invest in the particular project
		err = MunibondInvest(consts.OpenSolarIssuerDir, invIndex, invSeed, invAmount, projIndex,
			project.InvestorAssetCode, project.TotalValue, 1)
		if err != nil {
			log.Println("Error while seed investing", err)
			return errors.Wrap(err, "error while investing")
		}

		// once the investment is complete, update the project and store in the database
		err = project.updateAfterInvestment(invAmount, invIndex)
		if err != nil {
			return errors.Wrap(err, "failed to update project after investment")
		}
		return err
	} else {
		return errors.New("other chain investment not supported right now")
	}
}

// the updateAfterInvestment of the opensolar platform
func (project *Project) updateAfterInvestment(invAmount float64, invIndex int) error {
	// MW: It seems that all your messages strings relate to errors, but not to confirmed transactions. It would be useful to add those
	var err error

	project.MoneyRaised += invAmount
	project.InvestorIndices = append(project.InvestorIndices, invIndex)
	log.Println("INV INDEX: ", invIndex)
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}

	if project.MoneyRaised == project.TotalValue {
		project.Lock = true
		err = project.Save()
		if err != nil {
			return errors.Wrap(err, "couldn't save project")
		}

		err = project.sendRecipientNotification()
		if err != nil {
			return errors.Wrap(err, "error while sending notifications to recipient")
		}

		go sendRecipientAssets(project.Index)
	}

	// we need to udpate the project investment map here
	project.InvestorMap = make(map[string]float64) // make the map
	log.Println("INVESTOR INDICES: ", project.InvestorIndices)
	for i := range project.InvestorIndices {
		investor, err := RetrieveInvestor(project.InvestorIndices[i])
		if err != nil {
			return errors.Wrap(err, "error while retrieving investors, quitting")
		}

		log.Println(investor.U.StellarWallet.PublicKey, project.InvestorAssetCode)

		var balance1 float64
		var balance2 float64

		balance1, err = xlm.GetAssetBalance(investor.U.StellarWallet.PublicKey, project.InvestorAssetCode)
		if err != nil {
			balance1 = 0
			// return errors.Wrap(err, "error while retrieving asset balance, quitting")
		}

		balance2, err = xlm.GetAssetBalance(investor.U.StellarWallet.PublicKey, project.SeedAssetCode)
		if err != nil {
			balance2 = 0
			// do nothing, since the user hasn't invested in seed assets yet
			// return errors.Wrap(err, "error while retrieving asset balance, quitting")
		}

		balance := balance1 + balance2
		percentageInvestment := balance / project.TotalValue
		project.InvestorMap[investor.U.StellarWallet.PublicKey] = percentageInvestment
	}

	err = project.Save()
	log.Println("INVESTOR MAP: ", project.InvestorMap)
	if err != nil {
		return errors.Wrap(err, "error while saving project, quitting")
	}
	return nil
}

// sendRecipientNotification sends the notification to the recipient requesting them
// to logon to the platform and unlock the project that has just been invested in
func (project *Project) sendRecipientNotification() error {
	recipient, err := RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipient")
	}
	notif.SendUnlockNotifToRecipient(project.Index, recipient.U.Email)
	return nil
}

// UnlockProject unlocks a specific project that has just been invested in
func UnlockProject(username string, pwhash string, projIndex int, seedpwd string) error {
	log.Println("UNLOCKING PROJECT")
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	recipient, err := ValidateRecipient(username, pwhash)
	if err != nil {
		return errors.Wrap(err, "couldn't validate recipient")
	}

	if recipient.U.Index != project.RecipientIndex {
		return errors.New("Recipient Indices don't match, quitting!")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while decrpyting seed")
	}

	checkPubkey, err := wallet.ReturnPubkey(recpSeed)
	if err != nil {
		return errors.Wrap(err, "couldn't get public key from seed")
	}

	if checkPubkey != recipient.U.StellarWallet.PublicKey {
		log.Println("Invalid seed")
		return errors.New("Failed to unlock project")
	}

	if !project.Lock {
		return errors.New("Project not locked")
	}

	project.LockPwd = seedpwd
	project.Lock = false
	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}
	return nil
}

// sendRecipientAssets sends a recipient the debt asset and the payback asset associated with
// the opensolar platform
func sendRecipientAssets(projIndex int) error {
	startTime := utils.Unix()
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	for utils.Unix()-startTime < consts.LockInterval {
		log.Printf("WAITING FOR PROJECT %d TO BE UNLOCKED", projIndex)
		project, err = RetrieveProject(projIndex)
		if err != nil {
			return errors.Wrap(err, "Couldn't retrieve project")
		}
		if !project.Lock {
			log.Println("Project UNLOCKED IN LOOP")
			break
		}
		time.Sleep(10 * time.Second)
	}

	// lock is open, retrieve project and transfer assets
	project, err = RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	recipient, err := RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve recipienrt")
	}

	recpSeed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, project.LockPwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed")
	}

	escrowPubkey, err := escrow.InitEscrow(project.Index, consts.EscrowPwd, recipient.U.StellarWallet.PublicKey, recpSeed, consts.PlatformSeed)
	if err != nil {
		return errors.Wrap(err, "error while initializing issuer")
	}

	log.Println("successfully setup escrow")
	project.EscrowPubkey = escrowPubkey
	err = escrow.TransferFundsToEscrow(project.TotalValue, project.Index, project.EscrowPubkey, consts.PlatformSeed)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "could not transfer funds to the escrow, quitting!")
	}

	log.Println("Transferred funds to escrow!")
	project.LockPwd = "" // set lockpwd to nil immediately after retrieving seed
	metadata := project.Metadata

	project.DebtAssetCode = assets.AssetID(consts.DebtAssetPrefix + metadata)
	project.PaybackAssetCode = assets.AssetID(consts.PaybackAssetPrefix + metadata)

	err = MunibondReceive(consts.OpenSolarIssuerDir, project.RecipientIndex, projIndex, project.DebtAssetCode,
		project.PaybackAssetCode, project.EstimatedAcquisition, recpSeed, project.TotalValue, project.PaybackPeriod)
	if err != nil {
		return errors.Wrap(err, "error while receiving assets from issuer on recipient's end")
	}

	err = project.updateProjectAfterAcceptance()
	if err != nil {
		return errors.Wrap(err, "failed to update project after acceptance of asset")
	}

	return nil
}

// updateProjectAfterAcceptance updates the project after acceptance of investment
// by the recipient
func (project *Project) updateProjectAfterAcceptance() error {

	project.BalLeft = project.TotalValue
	project.Stage = Stage5.Number // set to stage 5 (after the raise is done, we need to wait for people to actually construct the solar panels)

	err := project.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save project")
	}

	go monitorPaybacks(project.RecipientIndex, project.Index)
	return nil
}

// Payback pays the platform back in STABLEUSD and DebtAsset and receives PaybackAssets
// in return. Price to be paid per month depends on the electricity consumed by the recipient
// in the particular time frame
// If we allow a user to hold balances in btc / xlm, we could direct them to exchange the coin for STABLEUSD
// (or we could setup a payment provider which accepts fiat + crypto and do this ourselves)

// Payback is called by the recipient when he chooses to pay towards the project according to the payback interval
func Payback(recpIndex int, projIndex int, assetName string, amount float64, recipientSeed string) error {

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve project")
	}

	if project.InvestmentType != "munibond" {
		return errors.New("other investment models are not supported right now, quitting")
	}

	pct, err := MunibondPayback(consts.OpenSolarIssuerDir, recpIndex, amount,
		recipientSeed, projIndex, assetName, project.InvestorIndices, project.TotalValue, project.EscrowPubkey)
	if err != nil {
		return errors.Wrap(err, "Error while paying back the issuer")
	}

	project.BalLeft -= (1 - pct) * amount // the balance left should be the percenteage paid towards the asset, which is the monthly bill. THe re st goes into  ownership
	project.AmountOwed -= amount          // subtract the amount owed so we can track progress of payments in the monitorPaybacks loop
	project.OwnershipShift += pct
	project.DateLastPaid = utils.Unix()

	if project.BalLeft == 0 {
		log.Println("YOU HAVE PAID OFF THIS ASSET's LOAN, TRANSFERRING FUTURE PAYMENTS AS OWNERSHIP ASSETS OWNERSHIP OF ASSET TO YOU")
		project.Stage = 9
		// ownership shift is complete, so future payments will be made towards what's
	}

	if project.OwnershipShift == 1 {
		// the recipient has paid off the asset completely. TODO: we need to transfer some sort
		// of document to the person identifying that they now own the project
		log.Println("You now own the asset completely, there is no need to pay money in the future towards this particular project")
		project.Stage = 8 // TODO: review where and how this stage transition should occur
		project.BalLeft = 0
		project.AmountOwed = 0
	}

	err = project.Save()
	if err != nil {
		return errors.Wrap(err, "coudln't save project")
	}

	// TODO: we need to distribute funds which were paid back to all the parties involved, but we do so only for the investor here
	err = DistributePayments(recipientSeed, project.EscrowPubkey, projIndex, amount)
	if err != nil {
		return errors.Wrap(err, "error while distributing payments")
	}

	return nil
}

// DistributePayments distributes the return promised as part of the project back to investors and pays the other entities involved in the project
func DistributePayments(recipientSeed string, escrowPubkey string, projIndex int, amount float64) error {
	// this should act as the service which redistributes payments received out to the parties involved
	// amount is the amount that we want to give back to the investors and other entities involved
	project, err := RetrieveProject(projIndex)
	if err != nil {
		errors.Wrap(err, "couldn't retrieve project, quitting!")
	}

	if project.EscrowLock {
		log.Println("project", project.Index, "'s escrow locked, can't send funds")
		return errors.New("project escrow locked, can't send funds")
	}

	var fixedRate float64
	if project.InterestRate != 0 {
		fixedRate = project.InterestRate
	} else {
		fixedRate = 0.05 // 5 % interest rate if rate not defined earlier
	}

	amountGivenBack := fixedRate * amount
	for pubkey, percentage := range project.InvestorMap {
		// send x to this pubkey
		txAmount := percentage * amountGivenBack
		// here we send funds from the 2of2 multisig. Platform signs by default
		err = escrow.SendFundsFromEscrow(project.EscrowPubkey, pubkey, recipientSeed, consts.PlatformSeed, txAmount, "returns")
		if err != nil {
			log.Println(err) // if there is an error with one payback, doesn't mean we should stop and wait for the others
			continue
		}
	}
	return nil
}

// CalculatePayback calculates the amount of payback assets that must be issued in relation
// to the total amount invested in the project
func (project Project) CalculatePayback(amount float64) float64 {
	amountPB := (amount / project.TotalValue) * float64(project.EstimatedAcquisition*12)
	return amountPB
}

// monitorPaybacks monitors whether the user is paying back regularly towards the given project
// thread has to be isolated since if this fails, we stop tracking paybacks by the recipient.
func monitorPaybacks(recpIndex int, projIndex int) {
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println("Couldn't retrieve project")
			time.Sleep(consts.OneWeekInSecond)
		}

		recipient, err := RetrieveRecipient(recpIndex)
		if err != nil {
			log.Println("Couldn't retrieve recipient")
			time.Sleep(consts.OneWeekInSecond)
		}

		guarantor, err := RetrieveEntity(project.GuarantorIndex)
		if err != nil {
			log.Println("couldn't retrieve guarantor")
			time.Sleep(consts.OneWeekInSecond)
		}
		// this will be our payback period and we need to check if the user pays us back

		period := float64(time.Duration(project.PaybackPeriod) * consts.OneWeekInSecond) // in seconds due to the const
		if period == 0 {
			period = 1 // for the test suite
		}
		timeElapsed := utils.Unix() - project.DateLastPaid // this would be in seconds (unix time)
		factor := float64(timeElapsed) / period
		project.AmountOwed += factor * oracle.MonthlyBill() // add the amount owed only if the time elapsed is more than one payback period
		// Reputation adjustments based on payback history:
		if factor <= 1 {
			// don't do anything since the user has been paying back regularly
			log.Println("User: ", recipient.U.Email, "is on track paying towards order: ", projIndex)
			// maybe even update reputation here on a fractional basis depending on a user's timely payments
		} else if factor > NormalThreshold && factor < AlertThreshold {
			// person has not paid back for one-two consecutive period, send gentle reminder
			notif.SendNicePaybackAlertEmail(projIndex, recipient.U.Email)
			time.Sleep(consts.OneWeekInSecond)
		} else if factor >= SternAlertThreshold && factor < DisconnectionThreshold {
			// person has not paid back for four consecutive cycles, send reminder
			notif.SendSternPaybackAlertEmail(projIndex, recipient.U.Email)
			for _, i := range project.InvestorIndices {
				// send an email to recipients to assure them that we're on the issue and will be acting
				// soon if the recipient fails to pay again.
				investor, err := RetrieveInvestor(i)
				if err != nil {
					log.Println(err)
					continue
				}
				if investor.U.Notification {
					notif.SendSternPaybackAlertEmailI(projIndex, investor.U.Email)
				}
			}
			notif.SendSternPaybackAlertEmailG(projIndex, guarantor.U.Email)
			time.Sleep(consts.OneWeekInSecond)
		} else if factor >= DisconnectionThreshold {
			// send a disconnection notice to the recipient and let them know we have redirected
			// power towards the grid. Also maybe email ourselves in this case so that we can
			// contact them personally to resolve the issue as soon as possible.
			notif.SendDisconnectionEmail(projIndex, recipient.U.Email)
			for _, i := range project.InvestorIndices {
				// send an email to recipients to assure them that we're on the issue and will be acting
				// soon if the recipient fails to pay again.
				investor, err := RetrieveInvestor(i)
				if err != nil {
					log.Println(err)
					time.Sleep(consts.OneWeekInSecond)
					continue
				}
				if investor.U.Notification {
					notif.SendDisconnectionEmailI(projIndex, investor.U.Email)
				}
			}
			// we have sent out emails to investors, send an email to the guarantor and cover first losses of investors
			notif.SendDisconnectionEmailG(projIndex, guarantor.U.Email)
			err = CoverFirstLoss(project.Index, guarantor.U.Index, project.AmountOwed)
			if err != nil {
				log.Println(err)
				time.Sleep(consts.OneWeekInSecond)
				continue
			}
		}

		time.Sleep(consts.OneWeekInSecond) // poll every week to check progress on payments
	}
}

func addWaterfallAccount(projIndex int, pubkey string, amount float64) error {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve project, quitting")
	}
	if project.WaterfallMap == nil {
		project.WaterfallMap = make(map[string]float64)
	}
	project.WaterfallMap[pubkey] = amount
	return project.Save()
}

// CoverFirstLoss covers first loss for investors byu sending funds from the guarantor's account
func CoverFirstLoss(projIndex int, entityIndex int, amount float64) error {
	// cover first loss for the project specified
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve projects from database, quitting")
	}

	entity, err := RetrieveEntity(entityIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve entity from database, quitting")
	}

	// we now have the entity and the project under question
	if project.GuarantorIndex != entity.U.Index {
		return errors.New("guarantor index does not match with entity's index in database")
	}

	if entity.FirstLossGuaranteeAmt < amount {
		log.Println("amount required greater than what guarantor agreed to provide, adjusting first loss to cover for what's available")
		amount = entity.FirstLossGuaranteeAmt
	}
	// we now need to send funds from the gurantor's account to the escrow
	seed, err := wallet.DecryptSeed(entity.U.StellarWallet.EncryptedSeed, entity.FirstLossGuarantee) //
	if err != nil {
		return errors.Wrap(err, "could not decrypt seed, quitting!")
	}

	var txhash string
	// we have the escrow's pubkey, transfer funds to the escrow
	if !consts.Mainnet {
		_, txhash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, project.EscrowPubkey, amount, seed, "first loss guarantee")
		if err != nil {
			return errors.Wrap(err, "could not transfer asset to escrow, quitting")
		}
	} else {
		_, txhash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, project.EscrowPubkey, amount, seed, "first loss guarantee")
		if err != nil {
			return errors.Wrap(err, "could not transfer asset to escrow, quitting")
		}
	}

	log.Println("txhash of guarantor kick in:", txhash)

	return nil
}
