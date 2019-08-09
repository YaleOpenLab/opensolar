package loader

import (
	"github.com/pkg/errors"
	"log"
	"os"

	// edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	// openxconsts "github.com/YaleOpenLab/openx/consts"

	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
)

func testSolarProject(index int, panelsize string, totalValue float64, location string, moneyRaised float64,
	metadata string, invAssetCode string, debtAssetCode string, pbAssetCode string, years int, recpIndex int,
	contractor core.Entity, originator core.Entity, stage int, pbperiod int, auctionType string) (core.Project, error) {

	var project core.Project
	project.Index = index
	project.PanelSize = panelsize
	project.TotalValue = totalValue
	project.State = location
	project.MoneyRaised = moneyRaised
	project.Metadata = metadata
	project.InvestorAssetCode = invAssetCode
	project.DebtAssetCode = debtAssetCode
	project.PaybackAssetCode = pbAssetCode
	project.DateInitiated = utils.Timestamp()
	project.EstimatedAcquisition = years
	project.RecipientIndex = recpIndex
	project.ContractorIndex = contractor.U.Index
	project.OriginatorIndex = originator.U.Index
	project.Stage = stage
	project.PaybackPeriod = pbperiod
	project.AuctionType = auctionType
	project.InvestmentType = "munibond"

	err := project.Save()
	if err != nil {
		return project, errors.New("Error inserting project into db")
	}
	return project, nil
}

func Testnet() error {
	consts.HomeDir += "/testnet"
	consts.DbDir = consts.HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	consts.OpenSolarIssuerDir = consts.HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	consts.PlatformSeedFile = consts.HomeDir + "/platformseed.hex" // where the platform's seed is stored

	if _, err := os.Stat(consts.HomeDir); os.IsNotExist(err) {
		// no home directory exists, create
		var err error
		core.CreateHomeDir()
		// populate database with dummy data
		log.Println("populating db with test data for testnet")

		// there is no recipient right now, so create a dummy recipient
		recp, err := core.NewRecipient("martin", "p", "x", "Martin")
		if err != nil {
			return err
		}

		_, err = core.NewInvestor("john", "p", "x", "John")
		if err != nil {
			return err
		}

		originator, err := core.NewOriginator("samuel", "p", "x", "SamuelOriginator", "ABC Street, London", "I am an originator")
		if err != nil {
			return err
		}

		contractor, err := core.NewContractor("sam", "p", "x", "SamuelContractor", "14 ABC Street London", "This is a competing contractor")
		if err != nil {
			return err
		}

		_, err = testSolarProject(1, "100 1000 sq.ft homes each with their own private spaces for luxury", 14000, "India Basin, San Francisco",
			0, "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate", "", "", "",
			3, recp.U.Index, contractor, originator, 4, 2, "blind")

		if err != nil {
			return err
		}

		_, err = testSolarProject(2, "180 1200 sq.ft homes in a high rise building 0.1mi from Kendall Square", 30000, "Kendall Square, Boston",
			0, "Kendall Square is set in the heart of Cambridge and is a popular startup IT hub", "", "", "",
			5, recp.U.Index, contractor, originator, 4, 2, "blind")

		if err != nil {
			return err
		}

		_, err = testSolarProject(3, "260 1500 sq.ft homes set in a medieval cathedral style construction", 40000, "Trafalgar Square, London",
			0, "Trafalgar Square is set in the heart of London's financial district, with big banks all over", "", "", "",
			7, recp.U.Index, contractor, originator, 4, 2, "blind")

		if err != nil {
			return err
		}
	}
	return nil
}
