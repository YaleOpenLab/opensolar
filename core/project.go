package core

import (
	// "log"
	platforms "github.com/YaleOpenLab/openx/platforms"
)

// Project defines the project investment structure in opensolar
type Project struct {
	// The project is split into two parts - parts which are used in the smart contract and parts which are not
	// we define them as critparams and noncritparams

	// start crit params
	Index                int                // an Index to keep track of how many projects exist
	TotalValue           float64            // the total money that we need from investors
	Lock                 bool               // lock investment in order to wait for recipient's confirmation
	LockPwd              string             // the recipient's seedpwd. Will be set to null as soon as we use it.
	Chain                string             // the chain on which the project desires to be.
	OneTimeUnlock        string             // a one time unlock password where the recipient can store his seedpwd in (will be decrypted after investment)
	AmountOwed           float64            // the amoutn owed to investors as a cumulative sum. Used in case of a breach
	Reputation           float64            // the positive reputation associated with a given project
	Votes                float64            // the number of votes towards a proposed contract by investors
	OwnershipShift       float64            // the percentage of the project that the recipient now owns
	StageData            []string           // the data associated with stage migrations
	StageChecklist       []map[string]bool  // the checklist that has to be completed before moving on to the next stage
	InvestorMap          map[string]float64 // publicKey: percentage donation
	SeedInvestorMap      map[string]float64 // the list of all seed investors who've invested in the project
	WaterfallMap         map[string]float64 // publickey:amount map in order to pay multiple accounts. A bit ugly, but should work fine. Make map before using
	RecipientIndex       int                // The index of the project's recipient
	OriginatorIndex      int                // the originator of the project
	GuarantorIndex       int                // the person guaranteeing the specific project in question
	ContractorIndex      int                // the person with the proposed contract
	InvestorIndices      []int              // The various investors who have invested in the project
	SeedInvestorIndices  []int              // Investors who took part before the contract was at stage 3
	RecipientIndices     []int              // the indices of the recipient family (offtakers, beneficiaries, etc)
	DateInitiated        string             // date the project was created on the platform
	DateFunded           string             // date that the project completed the stage 4-5 migration
	DateLastPaid         int64              // int64 ie unix time since we need comparisons on this one
	AuctionType          string             // the type of the auction in question. Default is blind auction unless explicitly mentioned
	InvestmentType       string             // the type of investment - equity crowdfunding, municipal bond, normal crowdfunding, etc defined in models
	PaybackPeriod        int                // the frequency in number of weeks that the recipient has to pay the platform
	Stage                int                // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
	InvestorAssetCode    string             // the code of the asset given to investors on investment in the project
	DebtAssetCode        string             // the code of the asset given to recipients on receiving a project
	PaybackAssetCode     string             // the code of the asset given to recipients on receiving a project
	SeedAssetCode        string             // the code of the asset given to seed investors on seed investment in the project
	SeedInvestmentFactor float64            // the factor that a seed investor's investment is multiplied by in case he does invest at the seed stage
	SeedInvestmentCap    float64            // the max amount that a seed investor can put in a project when it is in its seed stages
	EscrowPubkey         string             // the publickey of the escrow we setup after project investment
	EscrowLock           bool               // used to lock the escrow in case someting goes wrong
	MoneyRaised          float64            // total money that has been raised until now
	SeedMoneyRaised      float64            // total seed money that has been raised until now
	EstimatedAcquisition int                // the year in which the recipient is expected to repay the initial investment amount by
	BalLeft              float64            // denotes the balance left to pay by the party, percentage raised is not stored in the database since that can be calculated
	AdminFlagged         bool               // denotes if someone reports the project as flagged
	FlaggedBy            int                // the index of the admin who flagged the project
	UserFlaggedBy        []int              // the indices of the users who flagged the project
	Reports              int                // the number of reports against htis particular project

	// below are all the non critical params
	Name                        string     // the name of the project / the identifier by which its referred to
	State                       string     // the state in which the project has been installed in
	Country                     string     // the country in which the project has been installed in
	PanelSize                   string     // size of the given panel, for diplsaying to the user who wants to bid stuff
	PanelTechnicalDescription   string     // This should talk about '10x 100W Komaes etc'
	Inverter                    string     // the inverter of the installed project
	ChargeRegulator             string     // the charge regulator of the installed project
	ControlPanel                string     // the control panel of the installed project
	CommBox                     string     // the comm box of the installed project
	ACTransfer                  string     // the AC transfer of the installed project
	SolarCombiner               string     // the solar combiner of the installed project
	Batteries                   string     // the batteries of the installed project
	IoTHub                      string     // the IoT Hub installed as part of the project
	Rating                      string     // the rating of the project (Moody's, etc)
	Metadata                    string     // other metadata which does not have an explicit name can be stored here. Used to derive assetIDs
	InterestRate                float64    // the rate of return for investors
	Tax                         string     // the specifications of the tax system associated with this particular project
	ProposedInvestmentCap       float64    // the max amount that an investor can invest in when the project is in its proposed stage (stage 2)
	SelfFund                    float64    // the amount that a beneficiary / recipient puts in a project wihtout asking from other investors. This is not included as a seed investment because this would mean the recipient pays his own investment back in the project
	SecurityIssuer              string     // the issuer of the security
	BrokerDealer                string     // the broker dealer associated with the project
	MainDeveloperIndex          int        // the main developer of the project
	BlendedCapitalInvestorIndex int        // the index of the blended capital investor
	DeveloperIndices            []int      // the indices of the developers involved in the project`
	ContractorFee               float64    // fee paid to the contractor from the total fee of the project
	OriginatorFee               float64    // fee paid to the originator included in the total value of the project
	DeveloperFee                []float64  // the fees charged by the developers
	DebtInvestor1               string     // debt investor index, if any
	DebtInvestor2               string     // debt investor index, if any
	TaxEquityInvestor           string     // tax equity investor if any
	Terms                       []struct { // the terms of the project
		Variable      string
		Value         string
		RelevantParty string
		Note          string
		Status        string
		SupportDoc    string
	}
	ExecutiveSummary struct { // the bigger struct that holds all the executive summary metrics
		Investment            map[string]string
		Financials            map[string]string
		ProjectSize           map[string]string
		SustainabilityMetrics map[string]string
	}
	AutoReloadInterval float64  // the interval in which the user's funds reach zero
	ResilienceRating   float64  // resilience of the project
	ActionsRequired    string   // the action(s) required by the user
	Bullets            struct { // list of bullet points to be displayed on the frontend
		Bullet1 string
		Bullet2 string
		Bullet3 string
	}
	Hashes struct { // list of hashes to be displayed on the project details page
		LegalProjectOverviewHash string
		LegalPPAHash             string
		LegalRECAgreementHash    string
		GuarantorAgreementHash   string
		ContractorAgreementHash  string
		StakeholderAgreementHash string
		CommunityEnergyHash      string
		FinancialReportingHash   string
	}
	PendingDocuments map[int]string // a map of the user id and the document he needs to follow on. This text would match the index map in the user struct
	ContractList     []string       // the list of contracts to be displayed on the frontend
	Architecture     struct {       // the section labeled "Architecture/Project Design"
		SpaceLayoutImage   string
		SolarOutputImage   string
		SolarArray         string
		DailyAvgGeneration string
		System             string
		InverterSize       string
		DesignDescription  string
	}
	Context             string     // the section titled "Context"
	SummaryImage        string     // the url to the image linked in the summary
	CommunityEngagement []struct { // the section labelled "Community Engagement" on the frontend
		Width    int
		Title    string
		ImageURL string
		Content  string
		Link     string
	}
	ExplorePageSummary struct { // the summary on the explore page tab
		Solar   string
		Storage string
		Tariff  string
		Stage   int
		Return  string
		Rating  string
		Tax     string
		ETA     int
	}
	EngineeringLayoutType string
	FEText                map[string]interface{} // put all the fe text in here reading it from the relevant json file(s)
	MapLink               string                 // the google maps link to the installation site
}

// ExplorePageSummaryHelper defines the params that will appear on the frontend's explore page
type ExplorePageSummaryHelper struct {
}

// InvestmentHelper defines the investment specifics of the project
type InvestmentHelper struct {
	Capex              string
	Hardware           float64
	FirstLossEscrow    string
	CertificationCosts string
}

// FinancialHelper defines the financial specifics of the project
type FinancialHelper struct {
	Return    float64
	Insurance string
	Tariff    string
	Maturity  string
}

// ProjectSizeHelper defines size, storage and other params that are part of the project size section
type ProjectSizeHelper struct {
	PVSolar          string
	Storage          string
	Critical         float64
	InverterCapacity string
}

// SustainabilityHelper defines parameters relevant to sustainability that ae important to the project
type SustainabilityHelper struct {
	CarbonDrawdown string
	CommunityValue string
	LCA            string
}

// Feedback defines a structure that can be used for providing feedback about entities
type Feedback struct {
	Content string
	// the content of the feedback, good / bad
	// maybe we could have a rating system baked in? a star based rating system?
	// would be nice, idk
	From Entity
	// who gave the feedback?
	To Entity
	// regarding whom is this feedback about
	Date string
	// time at which this feedback was written
	Contract []Project
	// the contract regarding which this feedback is directed at
}

// Stage is the evolution of the erstwhile static stage integer construction
type Stage struct {
	Number          int
	FriendlyName    string   // the informal name that one can use while referring to the stage (nice for UI as well)
	Name            string   // this is a more formal name to give to the given stage
	Activities      []string // the activities that are covered in this particular stage and need to be fulfilled in order to move to the next stage.
	StateTrigger    []string // trigger state change from n to n+1
	BreachCondition []string // define breach conditions for a particular stage
}

// ContractAuction is an auction struct
type ContractAuction struct {
	// TODO: this struct isn't used yet as it needs handlers and stuff, but when
	// we move off main.go for testing, this must be used in order to make stuff
	// easier for us.
	AllContracts    []Project
	AllContractors  []Entity
	WinningContract Project
}

const (
	// InvestorWeight is the percentage weight of the project's total reputation assigned to the investor
	InvestorWeight = 0.1

	// OriginatorWeight is the percentage weight of the project's total reputation assigned to the originator
	OriginatorWeight = 0.1

	// ContractorWeight is the percentage weight of the project's total reputation assigned to the contractor
	ContractorWeight = 0.3

	// DeveloperWeight is the percentage weight of the project's total reputation assigned to the developer
	DeveloperWeight = 0.2

	// RecipientWeight is the percentage weight of the project's total reputation assigned to the recipient
	RecipientWeight = 0.3

	// NormalThreshold is the normal payback interval of 1 payback period. Regular notifications are sent regardless of whether the user has paid back towards the project.
	NormalThreshold = 1

	// AlertThreshold is the threshold above which the user gets a nice email requesting a quick payback whenever possible
	AlertThreshold = 2

	// SternAlertThreshold is the threshold above when the user gets a warning that services will be disconnected if the user doesn't payback soon.
	SternAlertThreshold = 4

	// DisconnectionThreshold is the threshold above which the user gets a notification telling that services have been disconnected.
	DisconnectionThreshold = 6
)

// SolarProjectArray is an array of Projects
type SolarProjectArray []Project

// InitializePlatform imports handlers from the main platform struct that are necessary for starting the platform
func InitializePlatform() error {
	return platforms.InitializePlatform()
}

// RefillPlatform checks whether the platform has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platforms.RefillPlatform(publicKey)
}
