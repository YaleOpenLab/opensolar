package consts

import (
	"os"
	"time"
)

// Start repeated params, use only for testing
var PlatformPublicKey = ""
var PlatformSeed = ""
var PlatformEmail = ""
var PlatformEmailPass = ""
var StablecoinCode = ""
var StablecoinPublicKey = ""
var AnchorUSDCode = ""
var AnchorUSDAddress = ""
var AnchorUSDTrustLimit = float64(10000)
var AnchorAPI = ""
var Mainnet = false

// End repeated params

var OpenxURL = "http://localhost:8080" // default openx instance to connect to
var TopSecretCode = "OPENSOLARTEST"     // code for requesting stuff from openx

// directories
var HomeDir = os.Getenv("HOME") + "/.opensolar" // home directory where we store everything

var DbName = "opensolar.db" // the name of the db that we want to store stuff in
var DbDir = ""              // the directory where the database is stored (project info, user info, etc)
var OpenSolarIssuerDir = "" // the directory where we store opensolar projects' issuer seeds
var PlatformSeedFile = ""   // where the platform's seed is stored

// prefixes
var InvestorAssetPrefix = "InvestorAssets_" // the prefix that will be hashed to give an investor AssetID
var BondAssetPrefix = "BondAssets_"         // the prefix that will be hashed to give a bond asset
var CoopAssetPrefix = "CoopAsset_"          // the prefix that will be hashed to give the cooperative asset
var DebtAssetPrefix = "DebtAssets_"         // the prefix that will be hashed to give a recipient AssetID
var SeedAssetPrefix = "SeedAssets_"         // the prefix that will be hashed to give an ivnestor his seed id
var PaybackAssetPrefix = "PaybackAssets_"   // the prefix that will be hashed to give a payback AssetID
var IssuerSeedPwd = "blah"                  // the password for unlocking the encrypted file. This must be modified a compile time and kept secret
var EscrowPwd = "blah"                      // the password used for locking the seed used by the escrow. This must be modified a compile time and kept secret

// ports + number consts
var Tlsport = 443                                           // default port for ssl
var DefaultRpcPort = 8081                                   // the default port on which the rpc server of the platform starts. Defaults to HTTPS
var LockInterval = int64(1 * 60 * 60 * 24 * 3)              // time a recipient is given to unlock the project and redeem investment, right now at 3 days
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)  // second * minute * hour * day * number, 30 days right now
var OneWeekInSecond = time.Duration(604800 * time.Second)   // one week in seconds
var TwoWeeksInSecond = time.Duration(1209600 * time.Second) // one week in seconds, easier to have it here than call it in multiple places
var SixWeeksInSecond = time.Duration(3628800 * time.Second) // six months in seconds, send notification
var CutDownPeriod = time.Duration(4838400 * time.Second)    // period when we direct power to the grid

// teller related consts
var TellerHomeDir = HomeDir + "/teller"                        // the home directory of the teller executable
var TellerMaxLocalStorageSize = 2000                           // in bytes, tweak this later to something like 10M after testing
var TellerPollInterval = time.Duration(30000 * time.Second)    // frequency with which the teller of a particular system is polled
var LoginRefreshInterval = time.Duration(5 * 60 * time.Second) // every 5 minutes we refresh the teller to import the changes on the platform

func SetTnConsts() {
	HomeDir = os.Getenv("HOME") + "/.opensolar/testnet"
	DbDir = HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	OpenSolarIssuerDir = HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	PlatformSeedFile = HomeDir + "/platformseed.hex" // where the platform's seed is stored
}

func SetMnConsts() {
	HomeDir = os.Getenv("HOME") + "/.opensolar/mainnet"
	DbDir = HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	OpenSolarIssuerDir = HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	PlatformSeedFile = HomeDir + "/platformseed.hex" // where the platform's seed is stored
}
