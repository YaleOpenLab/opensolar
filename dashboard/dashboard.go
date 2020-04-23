package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
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

type PersonFormat struct {
	Name     string
	Username string
	Email    string
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
	EscrowBalance   LinkFormat
	Recipient       PersonFormat
	Investor        PersonFormat
	Developer       PersonFormat
	ProjCount       int
	UserCount       int
}

var platformURL = "https://api2.openx.solar"
var AdminToken string
var Token string
var Project core.Project
var Recipient core.Recipient
var Return Content
var Developer core.Entity

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

func validateRecp(wg *sync.WaitGroup, username, token string) {
	defer wg.Done()
	body := "/recipient/validate?username=" + username + "&token=" + token
	var val string

	data, err := erpc.GetRequest(platformURL + body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &Recipient)
	if err != nil {
		log.Fatal(err)
	}

	if Recipient.U != nil {
		if Recipient.U.Index != 0 {
			val = "Validated Recipient"
		}
	}

	val = "Could not validate Recipient"

	Return.Validate.Text = val
	Return.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + Token
}

func getProject(wg *sync.WaitGroup, index int) error {
	defer wg.Done()
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

func serveStatic() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}

type length struct {
	Length int
}

func getUToken(wg *sync.WaitGroup, username string) {
	defer wg.Done()
	var err error

	Token, err = getToken(username, "password")
	if err != nil {
		log.Fatal("error while fetching recipient token: ", err)
	}
}

func getAToken(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	AdminToken, err = getToken("admin", "password")
	if err != nil {
		log.Fatal("error while fetching recipient token: ", err)
	}
}

func frontend() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		erpc.SetConsts(30)
		w.Header().Add("Content Type", "text/html")
		// The template name "template" does not matter here
		templates := template.New("template")
		doc, err := renderHTML()
		if err != nil {
			log.Fatal(err)
		}

		templates.New("doc").Parse(doc)

		Return.Title = "Opensolar status dashboard"
		Return.Name = "John"

		Return.OpensStatus.Text = "Opensolar is Down"
		Return.OpensStatus.Link = platformURL + "/ping"
		if opensPing() {
			Return.OpensStatus.Text = "Opensolar is Up"
		}

		Return.OpenxStatus.Text = "Openx is Down"
		Return.OpenxStatus.Link = "https://api.openx.solar/ping"
		if openxPing() {
			Return.OpenxStatus.Text = "Openx is Up"
		}

		username := "aibonitoGsIoJ"
		// get token

		var wg1 sync.WaitGroup

		wg1.Add(1)
		go getUToken(&wg1, username)
		wg1.Add(1)
		go getAToken(&wg1)
		wg1.Add(1)
		go validateRecp(&wg1, username, Token)
		wg1.Add(1)
		go getProject(&wg1, 1)
		wg1.Wait()

		invIndex, err := utils.ToString(Project.InvestorIndices[0])
		if err != nil {
			log.Fatal(err)
		}

		devIndex, err := utils.ToString(Project.DeveloperIndices[0])
		if err != nil {
			log.Fatal(err)
		}

		var wg2 sync.WaitGroup
		wg2.Add(1)
		go getInvestor(&wg2, AdminToken, invIndex)
		wg2.Add(1)
		go getDeveloper(&wg2, AdminToken, devIndex)
		wg2.Wait()

		log.Println("first sync group complete")

		if Project.DateLastPaid == 0 {
			Return.DateLastPaid.Text = "Date Last Paid: First Payment not yet made"
		} else {
			Return.DateLastPaid.Text = "Date Last Paid: " + utils.IntToHumanTime(Project.DateLastPaid)
		}

		if Recipient.NextPaymentInterval == "" {
			Return.NextInterval.Text = "Next Payment Interval: First Payment not yet made"
		} else {
			Return.NextInterval.Text = "Next Payment Interval: " + Recipient.NextPaymentInterval
		}

		Return.TellerEnergy.Text, err = utils.ToString(Recipient.TellerEnergy)
		if err != nil {
			log.Fatal(err)
		}

		Return.TellerEnergy.Text = "Energy generated till " + utils.Timestamp() + " is: " + Return.TellerEnergy.Text + " Wh"

		Return.DateLastStart.Text = utils.StringToHumanTime(Recipient.DeviceStarts[len(Recipient.DeviceStarts)-1])
		Return.DateLastStart.Text = "Last Boot Time: " + Return.DateLastStart.Text

		Return.DeviceID.Text = Recipient.DeviceId

		Return.DABalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.DebtAssetCode))
		if err != nil {
			log.Fatal(err)
		}

		Return.PBBalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.PaybackAssetCode))
		if err != nil {
			log.Fatal(err)
		}

		Return.DABalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey
		Return.PBBalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

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

		Return.AccountBalance1.Text = "XLM: " + pnbS + " STABLEUSD: " + pubS
		Return.AccountBalance1.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey
		Return.AccountBalance2.Text = "XLM: " + snbS + " STABLEUSD: " + subS
		Return.AccountBalance2.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

		escrowBalance := xlm.GetAssetBalance(Project.EscrowPubkey, "STABLEUSD")
		if escrowBalance < 0 {
			escrowBalance = 0
		}

		escrowBalanceS, err := utils.ToString(escrowBalance)
		if err != nil {
			log.Fatal(err)
		}

		Return.EscrowBalance.Text = escrowBalanceS
		Return.EscrowBalance.Link = "https://testnet.steexp.com/account/" + Project.EscrowPubkey

		Return.Recipient.Username = Recipient.U.Username
		Return.Recipient.Name = Recipient.U.Name
		Return.Recipient.Email = Recipient.U.Email

		var projCount length

		data, err := erpc.GetRequest("https://api2.openx.solar/admin/getallprojects?username=admin&token=" + AdminToken)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &projCount)
		if err != nil {
			log.Fatal(err)
		}

		Return.ProjCount = projCount.Length

		var userCount length

		data, err = erpc.GetRequest("https://api2.openx.solar/admin/getallusers?username=admin&token=" + AdminToken)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &userCount)
		if err != nil {
			log.Fatal(err)
		}

		Return.UserCount = userCount.Length

		projIndex, err := utils.ToString(Recipient.ReceivedSolarProjectIndices[0])
		if err != nil {
			log.Fatal(err)
		}

		data, err = erpc.GetRequest("https://api2.openx.solar/project/get?index=" + projIndex)
		if err != nil {
			log.Fatal(err)
		}

		var project core.Project
		err = json.Unmarshal(data, &project)
		if err != nil {
			log.Fatal(err)
		}

		templates.Lookup("doc").Execute(w, Return)
	})
}

func getInvestor(wg *sync.WaitGroup, AdminToken string, invIndex string) {
	defer wg.Done()
	data, err := erpc.GetRequest("https://api2.openx.solar/admin/getinvestor?username=admin&token=" +
		AdminToken + "&index=" + invIndex)
	if err != nil {
		log.Fatal(err)
	}

	var investor core.Investor
	err = json.Unmarshal(data, &investor)
	if err != nil {
		log.Fatal(err)
	}

	Return.Investor.Name = investor.U.Name
	Return.Investor.Username = investor.U.Username
	Return.Investor.Email = investor.U.Email
}

func getDeveloper(wg *sync.WaitGroup, AdminToken string, devIndex string) {
	defer wg.Done()
	data, err := erpc.GetRequest("https://api2.openx.solar/admin/getentity?username=admin&token=" +
		AdminToken + "&index=" + devIndex)
	if err != nil {
		log.Fatal(err)
	}

	var developer core.Entity
	err = json.Unmarshal(data, &developer)
	if err != nil {
		log.Fatal(err)
	}

	Return.Developer.Name = developer.U.Name
	Return.Developer.Username = developer.U.Username
	Return.Developer.Email = developer.U.Email
}

func renderHTML() (string, error) {
	doc, err := ioutil.ReadFile("index.html")
	return string(doc), err
}

func StartServer(portx int, insecure bool) {
	xlm.SetConsts(0, false)
	frontend()
	serveStatic()

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
