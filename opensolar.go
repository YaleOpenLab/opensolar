package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	loader "github.com/YaleOpenLab/opensolar/loader"
	rpc "github.com/YaleOpenLab/opensolar/rpc"

	openxconsts "github.com/YaleOpenLab/openx/consts"
	openxrpc "github.com/YaleOpenLab/openx/rpc"
)

var opts struct {
	Insecure bool   `short:"i" description:"Start the API using http. Not recommended"`
	Port     int    `short:"p" description:"The port on which the server runs on. Default: HTTPS/8081"`
	OpenxURL string `short:"o" description:"The URL of the openx instance to connect to. Default: http://localhost:8080"`
}

// ParseConfig parses CLI parameters
func ParseConfig(args []string) (bool, int, error) {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return false, -1, err
	}
	port := consts.DefaultRpcPort
	if opts.Port != 0 {
		port = opts.Port
	}
	if opts.OpenxURL != "" {
		consts.OpenxURL = opts.OpenxURL
	}
	return opts.Insecure, port, nil
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
func Mainnet() bool {
	body := consts.OpenxURL + "/mainnet"
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Fatal(err)
	}

	return data[0] == byte(0)
}

// ParseConsts parses consts by receiving consts from the openx API
func ParseConsts() error {
	body := consts.OpenxURL + "/platform/getconsts?code=" + consts.TopSecretCode
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Fatal(err)
	}

	var x openxrpc.OpensolarConstReturn
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
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
	openxconsts.DbDir = x.DbDir // for our retrieve methods

	return nil
}

func main() {
	var err error
	insecure, port, err := ParseConfig(os.Args) // parseconfig should be before StartPlatform to parse the mainnet bool
	if err != nil {
		log.Fatal(err)
	}

	if Mainnet() {
		consts.Mainnet = true
		openxconsts.SetConsts(true)
		// set mainnet db to open in spearate folder, no other way to do it than changing it here
		log.Println("initializing mainnet")

		err = ParseConsts()
		if err != nil {
			log.Fatal(err)
		}

		err = loader.Mainnet()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		openxconsts.SetConsts(false)
		log.Println("initializing testnet")

		err = ParseConsts()
		if err != nil {
			log.Fatal(err)
		}

		err = loader.Testnet()
		if err != nil {
			log.Fatal(err)
		}
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

	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	go core.MonitorTeller(1, "https://localhost:80")
	fmt.Println(`
		██████╗ ██████╗ ███████╗███╗   ██╗███████╗ ██████╗ ██╗      █████╗ ██████╗
	 ██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██╔═══██╗██║     ██╔══██╗██╔══██╗
	 ██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████╗██║   ██║██║     ███████║██████╔╝
	 ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║╚════██║██║   ██║██║     ██╔══██║██╔══██╗
	 ╚██████╔╝██║     ███████╗██║ ╚████║███████║╚██████╔╝███████╗██║  ██║██║  ██║
	  ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
		`)
	fmt.Println(`Starting Opensolar`)
	rpc.StartServer(port, insecure)
}
