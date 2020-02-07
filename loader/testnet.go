package loader

import (
	"os"

	// edb "github.com/Varunram/essentials/database"

	// openxconsts "github.com/YaleOpenLab/openx/consts"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
)

// Testnet loads the stuff needed for testnet. Ordering is very important since some consts need the others
// to function correctly
func Testnet() error {
	consts.HomeDir += "/testnet"
	consts.DbDir = consts.HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	consts.OpenSolarIssuerDir = consts.HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	consts.PlatformSeedFile = consts.HomeDir + "/platformseed.hex" // where the platform's seed is stored

	if _, err := os.Stat(consts.HomeDir); os.IsNotExist(err) {
		// no home directory exists, create
		core.CreateHomeDir()
	}
	return nil
}
