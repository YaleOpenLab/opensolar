package core

import (
	"github.com/pkg/errors"
	"log"
)

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

// Stage0 is the Handshake stage
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

// Stage1 is the engagement stage
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

// Stage2 is the quote stage
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

// Stage3 is the signing stage
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

// Stage4 is the raise stage
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

// Stage5 is the construction stage
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

// Stage6 is the connection stage
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

// Stage7 is the legacy stage
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

// Stage8 is the legacy stage
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

// Stage9 is the end of life stage
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
