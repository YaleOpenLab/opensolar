package rpc

import (
	"github.com/pkg/errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	// database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	// opzones "github.com/YaleOpenLab/openx/platforms/ozones"
)

func setupEntityRPCs() {
	validateEntity()
	getStage0Contracts()
	getStage1Contracts()
	getStage2Contracts()
	addCollateral()
	createOpensolarProject()
	proposeOpensolarProject()
	// createOpzonesCBond()
	// createOpzonesLuCoop()
}

// EntityValidateHelper is a helper that helps validate an entity
func EntityValidateHelper(w http.ResponseWriter, r *http.Request) (opensolar.Entity, error) {
	// first validate the investor or anyone would be able to set device ids
	erpc.CheckGet(w, r)
	var prepInvestor opensolar.Entity
	// need to pass the pwhash param here
	if r.URL.Query() == nil || r.URL.Query()["username"] == nil ||
		len(r.URL.Query()["pwhash"][0]) != 128 {
		return prepInvestor, errors.New("Invalid params passed")
	}

	prepEntity, err := opensolar.ValidateEntity(r.URL.Query()["username"][0], r.URL.Query()["pwhash"][0])
	if err != nil {
		return prepEntity, errors.Wrap(err, "Error while validating entity")
	}

	return prepEntity, nil
}

// validateEntity is an endpoint that vlaidates is a specific entity is registered on the platform
func validateEntity() {
	http.HandleFunc("/entity/validate", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		erpc.MarshalSend(w, prepEntity)
	})
}

// getStage0Contracts gets a list of all the pre origianted contracts on the platform
func getStage0Contracts() {
	http.HandleFunc("/entity/stage0", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		x, err := opensolar.RetrieveOriginatorProjects(opensolar.Stage0.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving originator project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// getStage1Contracts gets a list of all the originated contracts on the platform
func getStage1Contracts() {
	http.HandleFunc("/entity/stage1", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		x, err := opensolar.RetrieveOriginatorProjects(opensolar.Stage1.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving originator projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// getStage2Contracts gets a list of all the proposed contracts on the platform
func getStage2Contracts() {
	http.HandleFunc("/entity/stage2", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		x, err := opensolar.RetrieveContractorProjects(opensolar.Stage2.Number, prepEntity.U.Index)
		if err != nil {
			log.Println("Error while retrieving contractor projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, x)
	})
}

// addCollateral is a route that can be used to add collateral to a specific contractor who wishes
// to propose a contract towards a specific originated project.
func addCollateral() {
	//func (contractor *Entity) AddCollateral(amount float64, data string) error {
	http.HandleFunc("/entity/addcollateral", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}
		if r.URL.Query()["amount"] == nil || r.URL.Query()["collateral"] == nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		collateralAmount, err := utils.ToFloat(r.URL.Query()["amount"][0])
		if err != nil {
			log.Println("Error while converting string to float", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		collateralData := r.URL.Query()["collateral"][0]
		err = prepEntity.AddCollateral(collateralAmount, collateralData)
		if err != nil {
			log.Println("Error while adding collateral", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// createOpensolarProject creates a contract which the originator can take to the recipient in order to get validated
// as a level 1 project. This route can only be called by a valid entity (can be any entity though, since we can
// allow a shift between entity roles)
// http://localhost:8080/entity/newproject/opensolar?username=samuel&pwhash=9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758&TotalValue=10500&MoneyRaised=0&Years=5&InterestRate=5.5&Location=Bahams&PanelSize="15x24solarpanels"&Inverter=niceinverter&ChargeRegulator=electro&ControlPanel=nsa&CommBox=satellite&ACTransfer=tesla&SolarCombiner=solarcity&Batteries=siemens&IoTHub=rpi3&Metadata=innovaitveprojectoestablishspacestations&OriginatorFee=105.12&recpIndex=1&AuctionType=blind&PaybackPeriod=2
func createOpensolarProject() {
	http.HandleFunc("/entity/newproject/opensolar", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		// here we nered to validate the user but before that we need to decied which parameters are
		// essential and need to be present in order to create a contract. Since this in the test phase,
		// we will try to keep this as minimal as possible

		if r.URL.Query()["TotalValue"] == nil || r.URL.Query()["Years"] == nil || r.URL.Query()["InterestRate"] == nil ||
			r.URL.Query()["Location"] == nil || r.URL.Query()["PanelSize"] == nil || r.URL.Query()["Inverter"] == nil ||
			r.URL.Query()["ChargeRegulator"] == nil || r.URL.Query()["ControlPanel"] == nil || r.URL.Query()["CommBox"] == nil ||
			r.URL.Query()["ACTransfer"] == nil || r.URL.Query()["SolarCombiner"] == nil || r.URL.Query()["Batteries"] == nil ||
			r.URL.Query()["IoTHub"] == nil || r.URL.Query()["Metadata"] == nil || r.URL.Query()["OriginatorFee"] == nil ||
			r.URL.Query()["recpIndex"] == nil || r.URL.Query()["AuctionType"] == nil || r.URL.Query()["PaybackPeriod"] == nil {
			log.Println("Bad request, required params missing!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		allProjects, err := opensolar.RetrieveAllProjects()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// parse the parameters whhich are db intensive first and error out
		// so we don't make  too many db calls
		var x opensolar.Project
		x.TotalValue, err = utils.ToFloat(r.URL.Query()["TotalValue"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.EstimatedAcquisition, err = utils.ToInt(r.URL.Query()["Years"][0])
		if err != nil {
			log.Println("param passed not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.InterestRate, err = utils.ToFloat(r.URL.Query()["InterestRate"][0])
		if err != nil {
			log.Println("param passed not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.OriginatorFee, err = utils.ToFloat(r.URL.Query()["OriginatorFee"][0])
		if err != nil {
			log.Println("ORiginator fee not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		x.RecipientIndex, err = utils.ToInt(r.URL.Query()["recpIndex"][0])
		if err != nil {
			log.Println("passed recipient index not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		_, err = opensolar.RetrieveRecipient(x.RecipientIndex)
		if err != nil {
			log.Println("could not retrieve recipient, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		x.AuctionType = r.URL.Query()["AuctionType"][0]
		x.PaybackPeriod, err = utils.ToInt(r.URL.Query()["PaybackPeriod"][0])
		if err != nil {
			log.Println("payback period not integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// store other params which don't require db calls
		x.OriginatorIndex = prepEntity.U.Index
		x.State = r.URL.Query()["Location"][0]
		x.PanelSize = r.URL.Query()["PanelSize"][0]
		x.Inverter = r.URL.Query()["Inverter"][0]
		x.ChargeRegulator = r.URL.Query()["ChargeRegulator"][0]
		x.ControlPanel = r.URL.Query()["ControlPanel"][0]
		x.CommBox = r.URL.Query()["CommBox"][0]
		x.ACTransfer = r.URL.Query()["ACTransfer"][0]
		x.SolarCombiner = r.URL.Query()["SolarCombiner"][0]
		x.Batteries = r.URL.Query()["Batteries"][0]
		x.IoTHub = r.URL.Query()["IoTHub"][0]
		x.Metadata = r.URL.Query()["Metadata"][0]

		// autogen params that are easy to generate
		x.Index = len(allProjects) + 1
		x.Stage = 0
		x.MoneyRaised = 0
		x.BalLeft = x.TotalValue
		x.Votes = 0
		x.Reputation = x.TotalValue
		x.InvestmentType = "Munibond" // hardcode for now, expand if we have other investment models later down the road
		x.DateInitiated = utils.Timestamp()

		err = x.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// createOpensolarProject creates a contract which the originator can take to the recipient in order to get validated
// as a level 1 project. This route can only be called by a valid entity (can be any entity though, since we can
// allow a shift between entity roles)
func proposeOpensolarProject() {
	http.HandleFunc("/entity/proposeproject/opensolar", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		prepEntity, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if r.URL.Query()["projIndex"] == nil || r.URL.Query()["fee"] == nil {
			log.Println("missing required params, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("project idnex not int, quitting!")
			return
		}

		x, err := opensolar.RetrieveProject(projIndex)
		if err != nil {
			log.Println("couldn't retrieve project with index")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		fee, err := utils.ToFloat(r.URL.Query()["fee"][0])
		if err != nil {
			log.Println("fee passed not integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		// the below are the parameters we need to change. Add more here after consulting
		// what parameters a contractor can ideally change
		x.TotalValue += fee
		x.OriginatorFee = fee
		x.OriginatorIndex = prepEntity.U.Index
		x.Stage = 2 // set proposed contract stage

		err = x.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

/*
// http://localhost:8080/entity/newproject/opzone/constructionbond?username=samuel&pwhash=9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758&Title=opzonetest&Location=SFBay&Description=Mocksecription&InstrumentType=OpZoneConstruction&Amount=10million&CostOfUnit=200000&NoOfUnits=50&SecurityType=SEC1&Tax=10pcofffed&MaturationDate=2040&InterestRate=5.5&Rating=AAA&BondIssuer=FEDGOV&BondHolders=BHolder&Underwriter=WellsFargo
func createOpzonesCBond() {
	http.HandleFunc("/entity/newproject/opzone/constructionbond", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		_, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if r.URL.Query()["Title"] == nil || r.URL.Query()["Location"] == nil || r.URL.Query()["Description"] == nil ||
			r.URL.Query()["InstrumentType"] == nil || r.URL.Query()["Amount"] == nil || r.URL.Query()["CostOfUnit"] == nil ||
			r.URL.Query()["NoOfUnits"] == nil || r.URL.Query()["SecurityType"] == nil || r.URL.Query()["Tax"] == nil ||
			r.URL.Query()["MaturationDate"] == nil || r.URL.Query()["InterestRate"] == nil || r.URL.Query()["Rating"] == nil ||
			r.URL.Query()["BondIssuer"] == nil || r.URL.Query()["BondHolders"] == nil || r.URL.Query()["Underwriter"] == nil {
			log.Println("required params missing, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		var x opzones.ConstructionBond

		x.CostOfUnit, err = utils.ToFloat(r.URL.Query()["CostOfUnit"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.NoOfUnits, err = utils.ToInt(r.URL.Query()["NoOfUnits"][0])
		if err != nil {
			log.Println("param passed not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.InterestRate, err = utils.ToFloat(r.URL.Query()["InterestRate"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		allCBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("error while retreiveing all construction bonds, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		x.Index = len(allCBonds) + 1

		x.Title = r.URL.Query()["Title"][0]
		x.Location = r.URL.Query()["Location"][0]
		x.Description = r.URL.Query()["Description"][0]
		x.InstrumentType = r.URL.Query()["InstrumentType"][0]
		x.Amount = r.URL.Query()["Amount"][0]
		x.SecurityType = r.URL.Query()["SecurityType"][0]
		x.Tax = r.URL.Query()["Tax"][0]
		x.MaturationDate = r.URL.Query()["MaturationDate"][0]
		x.Rating = r.URL.Query()["Rating"][0]
		x.BondIssuer = r.URL.Query()["BondIssuer"][0]
		x.BondHolders = r.URL.Query()["BondHolders"][0]
		x.Underwriter = r.URL.Query()["Underwriter"][0]
		x.DateInitiated = utils.Timestamp()
		x.AmountRaised = 0

		err = x.Save()
		if err != nil {
			log.Println("error while saving project")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// http://localhost:8080/entity/newproject/opzone/lucoop?username=samuel&pwhash=9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758&Title=lucoop&Location=SFBay&Description=adfemolivingunitcoop&TypeOfUnit=transformable&Amount=300&SecurityType=SEC1&MaturationDate=2040&MonthlyPayment=3000&MemberRights=memberrights&InterestRate=5.5&Rating=AAA&BondIssuer=BWriter&Underwriter=WellsFargo&recpIndex=1
func createOpzonesLuCoop() {
	http.HandleFunc("/entity/newproject/opzone/lucoop", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		_, err := EntityValidateHelper(w, r)
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		var x opzones.LivingUnitCoop

		if r.URL.Query()["Title"] == nil || r.URL.Query()["Description"] == nil || r.URL.Query()["TypeOfUnit"] == nil ||
			r.URL.Query()["SecurityType"] == nil || r.URL.Query()["MaturationDate"] == nil || r.URL.Query()["MonthlyPayment"] == nil ||
			r.URL.Query()["MemberRights"] == nil || r.URL.Query()["InterestRate"] == nil || r.URL.Query()["Rating"] == nil ||
			r.URL.Query()["BondIssuer"] == nil || r.URL.Query()["Underwriter"] == nil || r.URL.Query()["recpIndex"] == nil {
			log.Println("required params not passed, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		x.Title = r.URL.Query()["Title"][0]
		x.Location = r.URL.Query()["Location"][0]
		x.Description = r.URL.Query()["Description"][0]
		x.TypeOfUnit = r.URL.Query()["TypeOfUnit"][0]
		x.Amount, err = utils.ToFloat(r.URL.Query()["Amount"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.SecurityType = r.URL.Query()["SecurityType"][0]
		x.MaturationDate = r.URL.Query()["MaturationDate"][0]
		x.MonthlyPayment, err = utils.ToFloat(r.URL.Query()["MonthlyPayment"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.MemberRights = r.URL.Query()["MemberRights"][0]
		x.InterestRate, err = utils.ToFloat(r.URL.Query()["InterestRate"][0])
		if err != nil {
			log.Println("param passed not float, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		x.Rating = r.URL.Query()["Rating"][0]
		x.BondIssuer = r.URL.Query()["BondIssuer"][0]
		x.Underwriter = r.URL.Query()["Underwriter"][0]

		x.RecipientIndex, err = utils.ToInt(r.URL.Query()["recpIndex"][0])
		if err != nil {
			log.Println("recpIndex not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		allLuCoops, err := opzones.RetrieveAllLivingUnitCoops()
		if err != nil {
			log.Println("Couldn't retriev all living unit coops, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		x.Index = len(allLuCoops) + 1
		x.DateInitiated = utils.Timestamp()
		x.UnitsSold = 0

		err = x.Save()
		if err != nil {
			log.Println("error while saving project")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}
*/
