package loader

import (
	"log"
	"os"

	xlm "github.com/Varunram/essentials/xlm"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
)

// Mainnet loads the stuff needed for mainnet. Ordering is very important since some consts need the others
// to function correctly
func Mainnet() error {
	consts.HomeDir += "/mainnet"
	consts.DbDir = consts.HomeDir + "/database/"
	consts.OpenSolarIssuerDir = consts.HomeDir + "/projects/"
	consts.PlatformSeedFile = consts.HomeDir + "/platformseed.hex"
	xlm.SetConsts(0, consts.Mainnet)

	if _, err := os.Stat(consts.HomeDir); os.IsNotExist(err) {
		// nothing exists, create dbs and buckets
		log.Println("creating mainnet home dir")
		// database.CreateHomeDir()
		core.CreateHomeDir()
		log.Println("created mainnet home dir")
		// Create an admin investor
		inv, err := core.NewInvestor("dci", "p", "x", "dci")
		if err != nil {
			return err
		}

		inv.U.Inspector = true
		inv.U.Kyc = true
		inv.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		inv.U.Reputation = 100000
		inv.U.Notification = true
		err = inv.U.AddEmail("varunramganesh@gmail.com") // change this to something more official later
		if err != nil {
			return err
		}

		// Create an admin recipient
		recp, err := core.NewRecipient("vg", "p", "x", "vg")
		if err != nil {
			return err
		}

		recp.U.Inspector = true
		recp.U.Kyc = true
		recp.U.Admin = true // no handlers for the admin bool, just set it wherever needed.
		recp.U.Reputation = 100000
		recp.U.Notification = true
		err = recp.U.AddEmail("varunramganesh@gmail.com")
		if err != nil {
			return err
		}

		orig, err := core.NewOriginator("martin", "p", "x", "MartinWainstein")
		if err != nil {
			return err
		}

		contractor, err := core.NewContractor("samuel", "p", "x", "SamuelVisscher")
		if err != nil {
			return err
		}

		var project core.Project
		project.Index = 1
		project.TotalValue = 5 // 5USD for testing
		project.Name = "SU Pasto School, Aibonito"
		project.Metadata = "MIT/Yale Pilot 2"
		project.OriginatorIndex = orig.U.Index
		project.ContractorIndex = contractor.U.Index
		project.EstimatedAcquisition = 5
		project.Stage = 4
		project.MoneyRaised = 0
		// add stuff in here as necessary
		err = project.Save()
		if err != nil {
			return err
		}

		log.Println("Please seed DCI pubkey: ", inv.U.StellarWallet.PublicKey, " with funds")
		log.Println("Please seed Varunram's pubkey: ", recp.U.StellarWallet.PublicKey, " with funds")
		log.Println("Please seed Martin's pubkey: ", orig.U.StellarWallet.PublicKey, " with funds")
		log.Println("Please seed Samuel's pubkey: ", contractor.U.StellarWallet.PublicKey, " with funds")
	}
	return nil
}
