package main

import (
	"fmt"
	// "github.com/pkg/errors"
	"log"
	"os"

	consts "github.com/YaleOpenLab/opensolar/consts"
	openxconsts "github.com/YaleOpenLab/openx/consts"
	loader "github.com/YaleOpenLab/opensolar/loader"
	rpc "github.com/YaleOpenLab/openx/rpc"
	erpc "github.com/Varunram/essentials/rpc"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
	Port     int  `short:"p" description:"The port on which the server runs on. Default: HTTPS/8081"`
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
	return opts.Insecure, port, nil
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

	if Mainnet() {
		openxconsts.DbDir = openxconsts.HomeDir + "/mainnet/"                           // set mainnet db to open in spearate folder
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
