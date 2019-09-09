package consts

import (
	"os"
	"time"
)

// PlatformPublicKey is the Stellar public key of the openx platform
var PlatformPublicKey string

// PlatformSeed is the Stellar seed corresponding to the above Stellar public key
var PlatformSeed string

// PlatformEmail is the email of the platform used to send notifications related to openx
var PlatformEmail string

// PlatformEmailPass is the password for the emial account linked above
var PlatformEmailPass string

// AdminEmail is the email of the platform's admin that is public
var AdminEmail string

// StablecoinCode is the code of the in house stablecoin that openx possesses
var StablecoinCode string

// StablecoinPublicKey is the public Stellar address of our in house stablecoin
var StablecoinPublicKey string

// AnchorUSDCode is the code for AnchorUSD's stablecoin
var AnchorUSDCode string

// AnchorUSDAddress is the address associated with AnchorUSD
var AnchorUSDAddress string

// AnchorUSDTrustLimit is the trust limit till which an account trusts AnchorUSD's stablecoin
var AnchorUSDTrustLimit float64

// AnchorAPI is the URL of AnchorUSD's API
var AnchorAPI string

// Mainnet denotes if openx is running on Stellar mainnet / testnet
var Mainnet bool

// OpenxURL is the openx URL that opensolar has to connect to
var OpenxURL = "http://localhost:8080"

// TopSecretCode is a test nuclear code that is used for testing
var TopSecretCode = "OPENSOLARTEST"

// HomeDir is the directory where opensolar projects, investors, entities, etc are stored
var HomeDir = os.Getenv("HOME") + "/.opensolar"

// DbDir is the directory where the openx database (boltDB) is stored
var DbDir = ""

// DbName is the name of the openx database
var DbName = "opensolar.db"

// OpenSolarIssuerDir is the directory where project escrow seeds are stored
var OpenSolarIssuerDir = ""

// PlatformSeedFile is the location where PlatformSeedFile is stored and decrypted each time the platform is started
var PlatformSeedFile string

// InvestorAssetPrefix is the prefix that will be hashed to give an investor AssetID
var InvestorAssetPrefix = "InvestorAssets_"

// DebtAssetPrefix is the prefix that will be hashed to give a recipient AssetID
var DebtAssetPrefix = "DebtAssets_"

// SeedAssetPrefix is the prefix that will be hashed to give an ivnestor his seed id
var SeedAssetPrefix = "SeedAssets_"

// PaybackAssetPrefix is the prefix that will be hashed to give a payback AssetID
var PaybackAssetPrefix = "PaybackAssets_"

// IssuerSeedPwd is the password of the issuer's seed
var IssuerSeedPwd = "blah"

// EscrowPwd is the password of the project escrows' seed
var EscrowPwd = "blah"

// Tlsport is the default SSL port on which openx starts
var Tlsport = 443

// DefaultRpcPort is the default Insecure port on which openx starts
var DefaultRpcPort = 8081

// LockInterval is the time a recipient is given to unlock the project and redeem investment, right now at 3 days
var LockInterval = int64(1 * 60 * 60 * 24 * 3)

// OneHour is one hour in seconds
var OneHour = time.Duration(1 * 60 * 60)

// PaybackInterval is the default teller payback interval
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)

// OneWeekInSecond is one week in seconds
var OneWeekInSecond = time.Duration(604800 * time.Second)

// TwoWeeksInSecond is two weeks in seconds
var TwoWeeksInSecond = time.Duration(1209600 * time.Second)

// SixWeeksInSecond is six weeks in seconds
var SixWeeksInSecond = time.Duration(3628800 * time.Second)

// CutDownPeriod is the period where we cut down power to the recipient and instead redirect it to the grid
var CutDownPeriod = time.Duration(4838400 * time.Second)

// TellerHomeDir is the home directory of the teller
var TellerHomeDir = HomeDir + "/teller" // the home directory of the teller executable

// TellerMaxLocalStorageSize is the max file storage limit on the teller before we hash the entire thing to ipfs
var TellerMaxLocalStorageSize = 2000

// TellerPollInterval is the frequency at which we poll the interval
var TellerPollInterval = time.Duration(3600 * 24 * time.Second)

// LoginRefreshInterval is the frequency at which the teller's credentials are updated (ie if you change your password, wait 5 minutes for the teller to disconnect)
var LoginRefreshInterval = time.Duration(5 * 60 * time.Second)

// ProjectReportThreshold is the threshold above which admins are allowed to flag the project
var ProjectReportThreshold = 10

// TellerListenNum is the number of messages the mqtt client will listen for from the broker
var TellerListenNum = 100

// TellerQos is the quality of service that the mqtt client must expect. Set to 0 (worst). Goes up until 2
var TellerQos = 0

// SetTnConsts sets constants that are relevant for staring opensolar on testnet
func SetTnConsts() {
	HomeDir = os.Getenv("HOME") + "/.opensolar/testnet"
	DbDir = HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	OpenSolarIssuerDir = HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	PlatformSeedFile = HomeDir + "/platformseed.hex" // where the platform's seed is stored
}

// SetMnConsts sets constants that are relevant for staring opensolar on mainnet // THIS IS UNUSED
func SetMnConsts() {
	HomeDir = os.Getenv("HOME") + "/.opensolar/mainnet"
	DbDir = HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
	OpenSolarIssuerDir = HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
	PlatformSeedFile = HomeDir + "/platformseed.hex" // where the platform's seed is stored
}
