package main

import (
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/signal"
	"strings"

	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	solar "github.com/YaleOpenLab/opensolar/core"
)

// package teller contains the remote client code that would be run on the client's
// side and communicate information with us and with atonomi and other partners.
// that it belongs, the contract, recipient, and eg. person who installed it.
// Consider doing this with IoT partners, eg. Atonomi.

// Teller authenticates with the platform using a remote API and then retrieves
// credentials once authenticated. Both the teller and the project recipient on the
// platform are the same entity, just that the teller is associated with the hw device.
// hw device needs an id and stuff, hopefully Atonomi can give us that. Else, we have a deviceId
// generated using a crypto random soruce,  hopefully should be sufficient.

// Teller tracks whenever the device starts and goes off, so we know when exactly the device was
// switched off. This is enough as proof that the device was running in between. This also
// avoids needing to poll the blockchain often and saves on the (minimal, still) tx fee.

// Since we can't compile this directly on the raspberry pi, we need to cross compile the
// go executable and transfer it over to the raspberry pi:
// env GOOS=linux GOARCH=arm GOARM=5 go build
// advisable to build off the pi and transport the executable since I don't think we want to be running
// go get on a raspberry pi with the stellar dependencies.
// one should run an ipfs node on the raspberry pi to ensure the teller can commit to ipfs without relying
// on the platform

var opts struct {
	Daemon     bool `short:"d" description:"Run the teller in daemon mode"`
	Port       int  `short:"p" description:"The port on which the teller runs on (default: 443)"`
	TestSwytch bool `long:"ts" description:"Test swytch API workflow"`
}

var (
	// LocalRecipient is the recipient struct associated with the project the teller is installed for
	LocalRecipient core.Recipient
	// LocalProject is the project that the teller is associated with
	LocalProject solar.Project
	// LocalProjIndex contains the project index the teller is associated with
	LocalProjIndex string
	// LocalSeedPwd contains the seed password of a user
	LocalSeedPwd string
	// RecpSeed stores the seed and PublicKey for easy vanity use
	RecpSeed string
	// RecpPublicKey is the receipient's PublicKey used to authenticate the teller
	RecpPublicKey string
	// PlatformPublicKey contains the platform parameters for interfacing with the platform
	PlatformPublicKey string
	// PlatformEmail is the platform's email address
	PlatformEmail string
	// ApiUrl is the API of the remote openx node
	ApiUrl string
	// DeviceId contains the device's id
	DeviceId string
	// DeviceLocation contains the device's location
	DeviceLocation string
	// DeviceInfo contains information on the user's device
	DeviceInfo string
	// StartHash records the blockhash when the teller starts and NowHash stores the blockhash at a particular instant
	StartHash string
	// NowHash is the hashchain has right now
	NowHash string
	// HashChainHeader is the header of the ipfs hash chain
	HashChainHeader string
	// SwytchUsername is the username that the teller has on the swytch platform
	SwytchUsername string
	// SwytchPassword is the password that the teller has on the swytch platform
	SwytchPassword string
	// SwytchClientid is the clientId associated with the given IoT Hub on swytch
	SwytchClientid string
	// SwytchClientSecret is the password associated with the given IoT Hub on swytch
	SwytchClientSecret string
	// AssetName is the asset for which this teller has been installed towards
	AssetName string
	// Token is the access token used to logon to the platform
	Token string
)

var cleanupDone chan struct{}

func autoComplete() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("help",
			readline.PcItem("update"),
			readline.PcItem("ping"),
			readline.PcItem("receive"),
			readline.PcItem("display"),
			readline.PcItem("update"),
			readline.PcItem("qq"),
			readline.PcItem("hh"),
		),
		readline.PcItem("display",
			readline.PcItem("balance",
				readline.PcItem("xlm"),
				readline.PcItem("asset"),
			),
			readline.PcItem("info"),
		),
	)
}

func main() {
	var err error
	xlm.SetConsts(10, false)
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal("Failed to parse arguments / Help command")
	}
	if opts.Port == 0 {
		opts.Port = consts.Tlsport
	}
	if opts.TestSwytch {
		testSwytch()
	}

	log.Println("---------------WELCOME TO THE TELLER INTERFACE---------------")
	defer recoverPanic() // catch any panics that may occur during the teller's runtime
	err = StartTeller()  // login to the platform, set device id, etc
	if err != nil {
		log.Fatal(err)
	}
	ColorOutput("TELLER PUBKEY: "+RecpPublicKey, GreenColor)
	ColorOutput("DEVICE ID: "+DeviceId, GreenColor)
	// testSwytch() tests the endpoints associated with the swytch platform
	// channels for preventing immediate sigint. Need this so that the action of any party which attempts
	// to close the teller would still be reported to the platform and emailed to the recipient
	signalChan := make(chan os.Signal, 1)
	cleanupDone = make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	StartHash, err = BlockStamp()
	if err != nil {
		log.Fatal(err)
	}
	// run goroutines in the background to routinely check for payback, state updates and stuff
	go checkPayback()
	// go updateState()
	// go storeDataLocal()

	if opts.Daemon {
		log.Println("Running teller in daemon mode")
		go func() {
			<-signalChan
			log.Println("\nSigint received, calling endhandler!")
			err = endHandler()
			for err != nil {
				log.Println(err)
				err = endHandler()
				<-cleanupDone
			}
			os.Exit(1)
		}()

		startServer(opts.Port) // run a daemon and listen for connections
		return                 // shouldn't come here, even if it does, we should be good
	}

	// non daemon mode, CLI available.
	go func() {
		<-signalChan
		log.Println("\nSigint received, not quitting wihtout closing endhandler!")
		close(cleanupDone)
	}()

	go startServer(opts.Port)

	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       promptColor("teller") + whiteColor("# "),
		HistoryFile:  consts.TellerHomeDir + "/history.txt",
		AutoComplete: autoComplete(),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			err := endHandler() // error, user wants to quit
			for err != nil {
				log.Println(err)
				err = endHandler()
				<-cleanupDone // to prevent user from quitting when endhandler is running
			}
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		ParseInput(cmdslice)
	}
}

// recoverPanic captures any unexpected panics that might occur and cause the teller to quit.
// even in such a situation, we would like to be warned so we can take some action
func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error) // recover the panic as an error
		log.Println("unexpected error, invoking EndHandler", err)
		err = endHandler()
		for err != nil { // run this loop until all endhandler functions are called
			log.Println(err)
			err = endHandler()
			<-cleanupDone // to prevent user from quitting when endhandler is running
		}
		os.Exit(1)
	}
}
