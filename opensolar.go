package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	loader "github.com/YaleOpenLab/opensolar/loader"
	rpc "github.com/YaleOpenLab/opensolar/rpc"

	openxconsts "github.com/YaleOpenLab/openx/consts"
)

var opts struct {
	Insecure bool   `short:"i" description:"Start the API using http. Not recommended"`
	Port     int    `short:"p" description:"The port on which the server runs on. Default: HTTPS/8081"`
	OpenxURL string `short:"o" description:"The URL of the openx instance to connect to. Default: http://localhost:8080"`
}

// ParseConfig parses CLI parameters passed
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

func ParseConsts() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading platform email from config file")
		return err
	}

	err = checkViperParams("PlatformPublicKey", "PlatformSeed", "PlatformEmail",
		"PlatformEmailPass", "StablecoinCode", "StablecoinPublicKey", "AnchorUSDCode",
		"AnchorUSDAddress", "AnchorUSDTrustLimit", "AnchorAPI", "Mainnet")
	if err != nil {
		return err
	}

	consts.PlatformPublicKey = viper.GetString("PlatformPublicKey")
	consts.PlatformSeed = viper.GetString("PlatformSeed")
	consts.PlatformEmail = viper.GetString("PlatformEmail")
	consts.PlatformEmailPass = viper.GetString("PlatformEmailPass")
	consts.StablecoinCode = viper.GetString("StablecoinCode")
	consts.StablecoinPublicKey = viper.GetString("StablecoinPublicKey")
	consts.AnchorUSDCode = viper.GetString("AnchorUSDCode")
	consts.AnchorUSDAddress = viper.GetString("AnchorUSDAddress")
	consts.AnchorUSDTrustLimit = viper.GetInt("AnchorUSDTrustLimit")
	consts.AnchorAPI = viper.GetString("AnchorAPI")
	consts.Mainnet = viper.GetBool("Mainnet")

	return nil
}

func Mainnet() bool {
	body := consts.OpenxURL + "/mainnet"
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Fatal(err)
	}

	return data[0] == byte(0)
}

func main() {
	var err error
	insecure, port, err := ParseConfig(os.Args) // parseconfig should be before StartPlatform to parse the mainnet bool
	if err != nil {
		log.Fatal(err)
	}

	err = ParseConsts()
	if err != nil {
		log.Fatal(err)
	}

	if Mainnet() {
		openxconsts.DbDir = openxconsts.HomeDir + "/mainnet/"
		// set mainnet db to open in spearate folder, no other way to do it than changing it here
		log.Println("MAINNET INIT")
		err = loader.Mainnet()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = loader.Testnet()
		if err != nil {
			log.Fatal(err)
		}
	}

	// rpc.KillCode = "NUKE" // compile time nuclear code
	// run this only when you need to monitor the tellers. Not required for local testing.
	// go opensolar.MonitorTeller(1)
	fmt.Println(`Starting Opensolar`)
	rpc.StartServer(port, insecure)
}
