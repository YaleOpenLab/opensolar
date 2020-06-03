package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"github.com/spf13/viper"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	loader "github.com/YaleOpenLab/opensolar/loader"
	rpc "github.com/YaleOpenLab/opensolar/rpc"

	// utils "github.com/Varunram/essentials/utils"
	stablecoin "github.com/Varunram/essentials/xlm/stablecoin"
	openxconsts "github.com/YaleOpenLab/openx/consts"
	openxrpc "github.com/YaleOpenLab/openx/rpc"
)

var opts struct {
	Insecure bool   `short:"i" description:"Start the API using http. Not recommended"`
	Port     int    `short:"p" description:"The port on which the server runs on. Default: HTTPS/8081"`
	DemoData bool   `short:"d" description:"Populate project"`
	Sandbox  bool   `short:"s" description:"Populate sandbox"`
	OpenxURL string `short:"o" description:"The URL of the openx instance to connect to. Default: http://localhost:8080"`
	EnvRead  bool   `short:"e" description:"read values from env files"`
}

// parseConfig parses CLI parameters
func parseConfig(args []string) (bool, int, error) {
	port := consts.DefaultRPCPort
	if opts.Port != 0 {
		port = opts.Port
	}
	if opts.OpenxURL != "" {
		consts.OpenxURL = opts.OpenxURL
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error while reading openx access code")
		log.Fatal(err)
	}

	err = checkViperParams("code")
	if err != nil {
		log.Fatal(err)
	}

	consts.TopSecretCode = viper.GetString("code")

	return opts.Insecure, port, nil
}

func parseEnvVars() (int, string, bool, bool, bool, string) {
	log.Println("reading")
	viper.AutomaticEnv()

	port := viper.GetInt("OPENS_PORT")
	code := viper.GetString("OPENX_CODE")
	populate := viper.GetBool("OPENS_POP")
	sandbox := viper.GetBool("OPENS_SB")
	insecure := viper.GetBool("OPENS_INSECURE")
	openxURL := viper.GetString("OPENX_URL")

	return port, code, populate, sandbox, insecure, openxURL
}

func checkViperParams(params ...string) error {
	for _, param := range params {
		if !viper.IsSet(param) {
			return errors.New("required param: " + param + " not found")
		}
	}
	return nil
}

// Mainnet calls openx's API to find out whether its running on testnet or mainnet
func mainnet() bool {
	body := consts.OpenxURL + "/mainnet"
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MAINNET BOOL: ", data[0] == byte(1))
	return data[0] == byte(1)
}

// loadOpenxConsts parses consts by receiving consts from the openx API
func loadOpenxConsts() error {
	body := consts.OpenxURL + "/platform/getconsts?code=" + consts.TopSecretCode
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Fatal(err)
	}

	var x openxrpc.OpensolarConstReturn
	err = json.Unmarshal(data, &x)
	if err != nil || len(string(data)) < 100 { // weird hack to catch 400s
		return errors.New("request to get consts failed")
	}

	consts.PlatformPublicKey = x.PlatformPublicKey
	consts.PlatformSeed = x.PlatformSeed
	consts.PlatformEmail = x.PlatformEmail
	consts.PlatformEmailPass = x.PlatformEmailPass
	consts.StablecoinCode = x.StablecoinCode
	consts.StablecoinPublicKey = x.StablecoinPublicKey
	consts.AnchorUSDCode = x.AnchorUSDCode
	consts.AnchorUSDAddress = x.AnchorUSDAddress
	consts.AnchorUSDTrustLimit = x.AnchorUSDTrustLimit
	consts.AnchorAPI = x.AnchorAPI
	consts.Mainnet = x.Mainnet
	log.Println("X MAINNET: ", x.Mainnet)
	openxconsts.DbDir = x.DbDir // for our retrieve methods

	stablecoin.SetConsts("STABLEUSD", consts.StablecoinPublicKey, "seed", "seedfile", openxconsts.StablecoinTrustLimit,
		consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.AnchorUSDTrustLimit, consts.Mainnet)

	return nil
}

func main() {
	var err error
	//log.Fatal(sandbox.Test())
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var insecure bool
	var port int

	if opts.EnvRead {
		port, consts.TopSecretCode, opts.DemoData,
			opts.Sandbox, insecure, consts.OpenxURL = parseEnvVars()
	} else {
		insecure, port, err = parseConfig(os.Args) // parseconfig should be before StartPlatform to parse the mainnet bool
		if err != nil {
			log.Fatal(err)
		}
	}

	consts.Mainnet = mainnet() // make an API call to openx for the status on this
	openxconsts.SetConsts(consts.Mainnet)

	err = loadOpenxConsts()
	if err != nil {
		log.Fatal(err)
	}

	if consts.Mainnet {
		// set mainnet db to open in spearate folder, no other way to do it than changing it here
		log.Println("initializing mainnet")
		err = loader.Mainnet()
		if err != nil {
			log.Fatal(err)
		}

		project, err := core.RetrieveProject(1)
		if err != nil {
			log.Fatal(err)
		}

		project.Metadata = "MAINNETTEST"
		project.InvestorAssetCode = ""
		project.TotalValue = 1
		project.MoneyRaised = 0
		project.InvestmentType = "munibond"
		project.RecipientIndex = 1
		project.DebtAssetCode = "TESTTELLER"

		err = project.Save()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("initializing testnet")

		err = loader.Testnet()
		if err != nil {
			log.Fatal(err)
		}

		if opts.DemoData {
			err = demoData()
			if err != nil {
				log.Fatal(err)
			}
		}

		if opts.Sandbox {
			err = sandbox()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	fmt.Println(`
	  ██████╗  ██████╗ ███████╗███╗  ██╗███████╗ ██████╗  ██╗      █████╗ ██████╗
	 ██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██╔═══██╗██║     ██╔══██╗██╔══██╗
	 ██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████╗██║   ██║██║     ███████║██████╔╝
	 ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║╚════██║██║   ██║██║     ██╔══██║██╔══██╗
	 ╚██████╔╝██║     ███████╗██║ ╚████║███████║╚██████╔╝███████╗██║  ██║██║  ██║
	  ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
		`)
	fmt.Println(`Starting Opensolar`)

	/*
		recp, err := core.RetrieveRecipient(7)
		if err != nil {
			log.Fatal(err)
		}

		recp.NextPaymentInterval = utils.IntToHumanTime(utils.Unix() + int64(consts.OneWeek))

		err = recp.Save()
		if err != nil {
			log.Fatal(err)
		}

		//go core.MonitorPaybacks(7, 1) // montior test project payback
	*/
	rpc.StartServer(port, insecure)
}
