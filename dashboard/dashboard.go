package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	tickers "github.com/Varunram/essentials/exchangetickers"
	"github.com/Varunram/essentials/xlm"

	erpc "github.com/Varunram/essentials/rpc"
	// osrpc "github.com/YaleOpenLab/opensolar/rpc"
	"github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	flags "github.com/jessevdk/go-flags"
)

type LinkFormat struct {
	Link string
	Text string
}

type Content struct {
	Title           string
	Name            string
	OpensStatus     LinkFormat
	OpenxStatus     LinkFormat
	Validate        LinkFormat
	NextInterval    LinkFormat
	TellerEnergy    LinkFormat
	DateLastPaid    LinkFormat
	DateLastStart   LinkFormat
	DeviceID        LinkFormat
	DABalance       LinkFormat
	PBBalance       LinkFormat
	AccountBalance1 LinkFormat
	AccountBalance2 LinkFormat
	Fruit           [3]string
}

var platformURL = "https://api2.openx.solar"
var Project core.Project
var Recipient core.Recipient

func opensPing() bool {
	data, err := erpc.GetRequest(platformURL + "/ping")
	if err != nil {
		log.Println(err)
		return false
	}

	var x erpc.StatusResponse

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return false
	}

	return x.Code == 200
}

func openxPing() bool {
	data, err := erpc.GetRequest("https://api.openx.solar/ping")
	if err != nil {
		log.Println(err)
		return false
	}

	var x erpc.StatusResponse

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return false
	}

	return x.Code == 200
}

func getToken(username, password string) (string, error) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("pwhash", utils.SHA3hash(password))

	retdata, err := erpc.PostForm(platformURL+"/token", form)
	if err != nil {
		log.Println(err)
		return "", err
	}

	type tokenResponse struct {
		Token string `json:"Token"`
	}

	var x tokenResponse

	err = json.Unmarshal(retdata, &x)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return x.Token, nil
}

func validateRecp(username, token string) (string, error) {
	body := "/recipient/validate?username=" + username + "&token=" + token
	data, err := erpc.GetRequest(platformURL + body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = json.Unmarshal(data, &Recipient)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if Recipient.U != nil {
		if Recipient.U.Index != 0 {
			return "Validated Recipient", nil
		}
	}

	return "Could not validate Recipient", nil
}

func getProject(index int) error {
	indexS, err := utils.ToString(index)
	if err != nil {
		log.Println(err)
		return err
	}

	body := "/project/get?index=" + indexS

	data, err := erpc.GetRequest(platformURL + body)
	if err != nil {
		log.Println(err)
		return err
	}

	return json.Unmarshal(data, &Project)
}

func frontend() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content Type", "text/html")
		// The template name "template" does not matter here
		templates := template.New("template")
		doc, err := renderHTML()
		if err != nil {
			log.Fatal(err)
		}

		templates.New("doc").Parse(doc)

		var x Content

		x.Title = "Opensolar status dashboard"
		x.Fruit = [3]string{"Apple", "Lemon", "Orange"}
		x.Name = "John"

		x.OpensStatus.Text = "Opensolar is Down"
		x.OpensStatus.Link = platformURL + "/ping"
		if opensPing() {
			x.OpensStatus.Text = "Opensolar is Up"
		}

		x.OpenxStatus.Text = "Openx is Down"
		x.OpenxStatus.Link = "https://api.openx.solar/ping"
		if openxPing() {
			x.OpenxStatus.Text = "Openx is Up"
		}

		username := "aibonitoGsIoJ"
		// get token
		token, err := getToken(username, "password")
		if err != nil {
			log.Fatal(err)
		}

		val, err := validateRecp(username, token)
		if err != nil {
			log.Fatal(err)
		}

		err = getProject(1)
		if err != nil {
			log.Fatal(err)
		}

		x.Validate.Text = val
		x.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + token

		if Project.DateLastPaid == 0 {
			x.DateLastPaid.Text = "Date Last Paid: First Payment not yet made"
		} else {
			x.DateLastPaid.Text = "Date Last Paid: " + utils.IntToHumanTime(Project.DateLastPaid)
		}

		if Recipient.NextPaymentInterval == 0 {
			x.NextInterval.Text = "Next Payment Interval: First Payment not yet made"
		} else {
			npiS, err := utils.ToString(Recipient.NextPaymentInterval)
			if err != nil {
				log.Fatal(err)
			}
			x.NextInterval.Text = "Next Payment Interval: " + npiS
		}

		x.TellerEnergy.Text, err = utils.ToString(Recipient.TellerEnergy)
		if err != nil {
			log.Fatal(err)
		}

		x.TellerEnergy.Text = "Energy generated till " + utils.Timestamp() + " is: " + x.TellerEnergy.Text + " Wh"

		x.DateLastStart.Text = utils.StringToHumanTime(Recipient.DeviceStarts[len(Recipient.DeviceStarts)-1])
		x.DateLastStart.Text = "Last Boot Time: " + x.DateLastStart.Text

		x.DeviceID.Text = "Device ID: " + Recipient.DeviceId

		x.DABalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.DebtAssetCode))
		if err != nil {
			log.Fatal(err)
		}

		x.PBBalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.PaybackAssetCode))
		if err != nil {
			log.Fatal(err)
		}

		x.DABalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey
		x.PBBalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

		xlmUSD, err := tickers.BinanceTicker()
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

		primNativeBalance := xlm.GetNativeBalance(Recipient.U.StellarWallet.PublicKey) * xlmUSD
		if primNativeBalance < 0 {
			primNativeBalance = 0
		}

		secNativeBalance := xlm.GetNativeBalance(Recipient.U.SecondaryWallet.PublicKey) * xlmUSD
		if secNativeBalance < 0 {
			secNativeBalance = 0
		}

		primUsdBalance := xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, "STABLEUSD")
		if primUsdBalance < 0 {
			primUsdBalance = 0
		}

		secUsdBalance := xlm.GetAssetBalance(Recipient.U.SecondaryWallet.PublicKey, "STABLEUSD")
		if secUsdBalance < 0 {
			secUsdBalance = 0
		}

		pnbS, err := utils.ToString(primNativeBalance)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

		snbS, err := utils.ToString(secNativeBalance)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

		pubS, err := utils.ToString(primUsdBalance)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

		subS, err := utils.ToString(secUsdBalance)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

		x.AccountBalance1.Text = "XLM: " + pnbS + " STABLEUSD: " + pubS
		x.AccountBalance2.Text = "XLM: " + snbS + " STABLEUSD: " + subS
		templates.Lookup("doc").Execute(w, x)
	})
}

func renderHTML() (string, error) {
	doc, err := ioutil.ReadFile("index.html")
	return string(doc), err
}

func StartServer(portx int, insecure bool) {
	xlm.SetConsts(0, false)
	frontend()

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Println("starting server in insecure mode")
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}

var opts struct {
	Port     int  `short:"p" description:"The port on which the server runs on" default:"8081"`
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	StartServer(opts.Port, opts.Insecure)
}
