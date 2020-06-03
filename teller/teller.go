package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
	core "github.com/YaleOpenLab/opensolar/core"
	solar "github.com/YaleOpenLab/opensolar/core"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

// env GOOS=linux GOARCH=arm GOARM=5 go build

var opts struct {
	Daemon     bool   `short:"d" description:"Run the teller in daemon mode"`
	Port       int    `short:"p" description:"The port on which the teller runs on (default: 443)"`
	TestSwytch bool   `long:"ts" description:"Test swytch API workflow"`
	URL        string `short:"u" description:"The URL of the remote opensolar instance"`
}

var (
	// Token is the access token used to logon to the platform
	Token string
	// Username is the Username used to logon to any openx based platform
	loginUsername string
	// Pwhash is the Pwhash used to logon to any openx based platform
	loginPwhash string
	// loginProjIndex is the project index used for fetching the project initially
	loginProjIndex int
	// LocalRecipient is the recipient struct associated with the project the teller is installed for
	LocalRecipient core.Recipient
	// LocalProject is the project that the teller is associated with
	LocalProject solar.Project
	// LocalSeedPwd contains the seed password of a user
	LocalSeedPwd string
	// APIURL is the API of the remote opensolar instance
	APIURL string
	// AssetName is the asset for which this teller has been installed towards
	AssetName string

	// STATE variables

	// DeviceID contains the device's id
	DeviceID string
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
	// EnergyValue is the net amount of energy accumulated during a day
	EnergyValue uint32

	// GOOGLE MAPS VARIABLES

	// Mapskey is the API key of google maps needed to store the location of the teller
	Mapskey string
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

// SubscribeMessage subscribes to a broker
func SubscribeMessage(mqttopts *mqtt.ClientOptions, topic string, qos int, num int) error {
	colorOutput(CyanColor, "starting mqtt subscriber", mqttopts)
	receiveCount := 0
	receiver := make(chan [2]string)
	var messages []string

	mqttopts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		receiver <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(mqttopts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		colorOutput(CyanColor, "MQTT SUBSCRIBE ERROR: ", token.Error())
		return token.Error()
	}

	token := client.Subscribe(topic, byte(qos), nil)
	if token.Wait() && token.Error() != nil {
		colorOutput(CyanColor, "MQTT SUBSCRIBE ERROR: ", token.Error())
		return token.Error()
	}

	for receiveCount < num {
		incoming := <-receiver
		messages = append(messages, incoming[1:]...)

		fPath := "data.txt"
		f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return errors.Wrap(err, "could not open data file for reading")
		}

		defer f.Close()

		for _, data := range incoming[1:] {
			_, err = f.WriteString(data)
			if err != nil {
				return errors.Wrap(err, "could not write data to the hc file")
			}
		}
		colorOutput(YellowColor, "RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}

	client.Disconnect(250)
	colorOutput(CyanColor, "Subscriber Disconnected")
	colorOutput(CyanColor, "MESSAGES: ", messages)
	return nil
}

func parseConfig() error {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal("Failed to parse arguments / Help command")
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Error while reading values from config file")
	}

	requiredParams := []string{"platformPublicKey", "seedpwd", "username",
		"password", "apiurl", "mapskey", "projIndex", "assetName",
		"mqttbroker", "mqttusername", "mqttpassword", "mqtttopic"}

	for _, param := range requiredParams {
		if !viper.IsSet(param) {
			return errors.New("required param: " + param + " not found")
		}
	}

	LocalSeedPwd = viper.GetString("seedpwd")
	loginUsername = viper.GetString("username")
	loginPwhash = utils.SHA3hash(viper.GetString("password"))
	loginProjIndex = viper.GetInt("projIndex")
	APIURL = viper.GetString("apiurl")
	Mapskey = viper.GetString("mapskey")
	AssetName = viper.GetString("assetName")

	// parse optional params
	SwytchUsername = viper.GetString("susername")
	SwytchPassword = viper.GetString("spassword")
	SwytchClientid = viper.GetString("sclientid")
	SwytchClientSecret = viper.GetString("sclientsecret")

	// parse params needed by the subscriber
	mqttopts := mqtt.NewClientOptions()
	mqttopts.AddBroker(viper.GetString("mqttbroker"))
	mqttopts.SetClientID(viper.GetString("mqttusername"))
	mqttopts.SetUsername(viper.GetString("mqttusername"))
	mqttopts.SetPassword(viper.GetString("mqttpassword"))
	topic := viper.GetString("mqtttopic")
	qos := 0
	num := 10000000 // set this to a very high number

	go SubscribeMessage(mqttopts, topic, qos, num)

	if opts.Port == 0 {
		opts.Port = consts.Tlsport
	}
	if opts.TestSwytch {
		testSwytch()
	}
	if opts.URL != "" {
		APIURL = opts.URL
	}

	return nil
}

func main() {
	var err error
	erpc.SetConsts(60) // set rpc timeout to 60s to allow for slower RPC connections
	// this is inline with the https clisent setup for remote RPC calls

	err = parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	colorOutput(YellowColor, "---------------WELCOME TO THE TELLER INTERFACE---------------")
	defer recoverPanic() // catch any panics that occur during the teller's runtime
	err = StartTeller()  // login to the platform, set device id, etc
	if err != nil {
		log.Fatal(err)
	}

	colorOutput(YellowColor, "TELLER PUBKEY: "+LocalRecipient.U.StellarWallet.PublicKey)
	colorOutput(YellowColor, "DEVICE ID: "+DeviceID)

	signalChan := make(chan os.Signal, 1)
	cleanupDone = make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	var StartHash string
	var balance float64
	var usdBalance float64

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		StartHash, err = getLatestBlockHash()
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		balance, err = getNativeBalance()
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	if consts.Mainnet {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			usdBalance, err = getAssetBalance("USD")
			if err != nil {
				log.Fatal(err)
			}
		}(&wg)

	} else {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			usdBalance, err = getAssetBalance("STABLEUSD")
			if err != nil {
				log.Fatal(err)
			}
		}(&wg)
	}

	wg.Wait()

	colorOutput(MagentaColor, "XLM BALANCE: ", balance)
	colorOutput(MagentaColor, "USD BALANCE: ", usdBalance)
	colorOutput(MagentaColor, "START HASH: ", StartHash)

	go func() {
		// run goroutines in the background to routinely check for payback, state updates and stuff
		readEnergyData()
		time.Sleep(15 * time.Second)
		checkPayback()
		time.Sleep(15 * time.Second)
		updateState(true)
	}()

	if opts.Daemon {
		colorOutput(CyanColor, "Running teller in daemon mode")
		go func() {
			<-signalChan
			colorOutput(CyanColor, "\nSigint received, calling endhandler!")
			err = endHandler()
			for err != nil {
				colorOutput(CyanColor, err)
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
		colorOutput(CyanColor, "\nSigint received, not quitting without closing endhandler!")
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
		msg, err := rl.Readline()
		if err != nil {
			err := endHandler() // error, user wants to quit
			for err != nil {
				colorOutput(CyanColor, err)
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
		colorOutput(YellowColor, "entered command: ", msg)

		ParseInput(cmdslice)
	}
}

// recoverPanic captures any unexpected panics that might occur and cause the teller to quit.
// even in such a situation, we would like to be warned so we can take some action
func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error) // recover the panic as an error
		colorOutput(CyanColor, "unexpected error, invoking EndHandler", err)
		err = endHandler()
		for err != nil { // run this loop until all endhandler functions are called
			colorOutput(CyanColor, err)
			err = endHandler()
			<-cleanupDone // to prevent user from quitting when endhandler is running
		}
		os.Exit(1)
	}
}
