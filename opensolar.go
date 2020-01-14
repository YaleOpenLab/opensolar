package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	loader "github.com/YaleOpenLab/opensolar/loader"
	rpc "github.com/YaleOpenLab/opensolar/rpc"
	// utils "github.com/Varunram/essentials/utils"
	//sandbox "github.com/YaleOpenLab/opensolar/sandboxv2"

	openxconsts "github.com/YaleOpenLab/openx/consts"
	openxrpc "github.com/YaleOpenLab/openx/rpc"
)

var opts struct {
	Insecure bool   `short:"i" description:"Start the API using http. Not recommended"`
	Port     int    `short:"p" description:"The port on which the server runs on. Default: HTTPS/8081"`
	OpenxURL string `short:"o" description:"The URL of the openx instance to connect to. Default: http://localhost:8080"`
}

// parseConfig parses CLI parameters
func parseConfig(args []string) (bool, int, error) {
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

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
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

	return nil
}

func main() {
	var err error
	//log.Fatal(sandbox.Test())
	insecure, port, err := parseConfig(os.Args) // parseconfig should be before StartPlatform to parse the mainnet bool
	if err != nil {
		log.Fatal(err)
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

		project.ExploreStub.Name = "This is a sample project Name"
		project.ExploreStub.Location = "Puerto Rico"
		project.ExploreStub.DonationType = "donation"
		project.ExploreStub.Originator = "Project Originator"
		project.ExploreStub.Description = "Project Description"
		project.ExploreStub.Bullet1 = "This is a sample bullet"
		project.ExploreStub.Bullet2 = "This is a sample bullet"
		project.ExploreStub.Bullet3 = "This is a sample bullet"
		project.ExploreStub.Solar = "Solar"
		project.ExploreStub.Storage = "This is a sample storage ipsum"
		project.ExploreStub.Tariff = "Unlimited tariff"
		project.ExploreStub.Return = "Unlimited Return"
		project.ExploreStub.Rating = "AAA"
		project.ExploreStub.Tax = "1000"
		project.ExploreStub.Acquisition = "Sample Acquisition"

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

		var project core.Project
		project.Index = 1
		project.Metadata = "TESTPROJECT"
		project.InvestorAssetCode = ""
		project.TotalValue = 1
		project.MoneyRaised = 0
		project.InvestmentType = "munibond"
		project.RecipientIndex = 1
		project.DebtAssetCode = "TESTTELLER"
		project.Metadata = "MAINNETTEST"
		project.InvestorAssetCode = ""
		project.TotalValue = 1
		project.MoneyRaised = 0
		project.InvestmentType = "munibond"
		project.RecipientIndex = 1
		project.DebtAssetCode = "TESTTELLER"
		project.ExploreStub.Name = "This is a sample project Name"
		project.ExploreStub.Location = "Puerto Rico"
		project.ExploreStub.DonationType = "donation"
		project.ExploreStub.Originator = "Project Originator"
		project.ExploreStub.Description = "Project Description"
		project.ExploreStub.Bullet1 = "This is a sample bullet"
		project.ExploreStub.Bullet2 = "This is a sample bullet"
		project.ExploreStub.Bullet3 = "This is a sample bullet"
		project.ExploreStub.Solar = "Solar"
		project.ExploreStub.Storage = "This is a sample storage ipsum"
		project.ExploreStub.Tariff = "Unlimited tariff"
		project.ExploreStub.Return = "Unlimited Return"
		project.ExploreStub.Rating = "AAA"
		project.ExploreStub.Tax = "1000"
		project.ExploreStub.Acquisition = "Sample Acquisition"

		err = project.Save()
		if err != nil {
			log.Fatal(err)
		}
	}

	/*
		errs := make(chan error, 1)
		go core.TrackProject(1, "localhost:1883", "test", errs)
		err = <-errs
		if err != nil {
			log.Fatal(err)
		}
	*/
	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	// go core.MonitorTeller(1, "https://localhost:80")
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
