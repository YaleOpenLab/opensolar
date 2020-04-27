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
	"time"

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

type PingFormat struct {
	Link string
	Text string
	URL  string
}

type PersonFormat struct {
	Name     string
	Username string
	Email    string
}

type AdminFormat struct {
	Username   string
	Password   string
	AdminToken string
	RecpToken  string
}

type Content struct {
	Title            string
	Name             string
	OpensStatus      PingFormat
	OpenxStatus      PingFormat
	BuildsStatus     PingFormat
	WebStatus        PingFormat
	Validate         LinkFormat
	NextInterval     LinkFormat
	TellerEnergy     LinkFormat
	DateLastPaid     LinkFormat
	DateLastStart    LinkFormat
	DeviceID         LinkFormat
	DABalance        LinkFormat
	PBBalance        LinkFormat
	AccountBalance1  LinkFormat
	AccountBalance2  LinkFormat
	EscrowBalance    LinkFormat
	Recipient        PersonFormat
	Investor         PersonFormat
	Developer        PersonFormat
	PastEnergyValues []uint32
	DeviceLocation   string
	StateHashes      []string
	PaybackPeriod    time.Duration
	BalanceLeft      float64
	OwnershipShift   float64
	DateInitiated    string
	Stage            int
	DateFunded       string
	InvAssetCode     string
	ProjCount        LinkFormat
	UserCount        LinkFormat
	InvCount         LinkFormat
	RecpCount        LinkFormat
	Admin            AdminFormat
	Date             string
}

var platformURL = "https://api2.openx.solar"
var AdminToken string
var Token string
var Pnb string
var Pub string
var Snb string
var Sub string
var XlmUSD float64
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

func buildsPing() bool {
	data, err := erpc.GetRequest("https://builds.openx.solar/ping")
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

func websitePing() bool {
	data, err := erpc.GetRequest("https://openx.solar")
	if err != nil {
		log.Println(err)
		return false
	}

	return string(data)[2:14] == "doctype html"
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
			Return.Validate.Text = "Validated Recipient"
			Return.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + Token
		}
	} else {
		Return.Validate.Text = "Could not validate Recipient"
		Return.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + Token
	}

	var wg3 sync.WaitGroup
	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		primNativeBalance := xlm.GetNativeBalance(Recipient.U.StellarWallet.PublicKey) * XlmUSD
		if primNativeBalance < 0 {
			primNativeBalance = 0
		}

		var err error
		Pnb, err = utils.ToString(primNativeBalance)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		primUsdBalance := xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, "STABLEUSD")
		if primUsdBalance < 0 {
			primUsdBalance = 0
		}

		var err error
		Pub, err = utils.ToString(primUsdBalance)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		secNativeBalance := xlm.GetNativeBalance(Recipient.U.SecondaryWallet.PublicKey) * XlmUSD
		if secNativeBalance < 0 {
			secNativeBalance = 0
		}

		var err error
		Snb, err = utils.ToString(secNativeBalance)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		secUsdBalance := xlm.GetAssetBalance(Recipient.U.SecondaryWallet.PublicKey, "STABLEUSD")
		if secUsdBalance < 0 {
			secUsdBalance = 0
		}

		var err error
		Sub, err = utils.ToString(secUsdBalance)
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		Return.DABalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.DebtAssetCode))
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		Return.PBBalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.PaybackAssetCode))
		if err != nil {
			log.Fatal(err)
		}
	}(&wg3)

	wg3.Wait()
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

		var wgPre sync.WaitGroup

		wgPre.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.OpensStatus.Text = "Opensolar is Down"
			Return.OpensStatus.Link = platformURL + "/ping"
			if opensPing() {
				Return.OpensStatus.Text = "Opensolar is Up"
			}
			Return.OpensStatus.URL = "api2.openx.solar"
		}(&wgPre)

		wgPre.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.OpenxStatus.Text = "Openx is Down"
			Return.OpenxStatus.Link = "https://api.openx.solar/ping"
			if openxPing() {
				Return.OpenxStatus.Text = "Openx is Up"
			}
			Return.OpenxStatus.URL = "api.openx.solar"
		}(&wgPre)

		wgPre.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.BuildsStatus.Text = "Builds is Down"
			Return.BuildsStatus.Link = "https://builds.openx.solar/ping"
			if buildsPing() {
				Return.BuildsStatus.Text = "Builds is Up"
			}
			Return.BuildsStatus.URL = "builds.openx.solar"
		}(&wgPre)

		wgPre.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			Return.WebStatus.Text = "Website is Down"
			Return.WebStatus.Link = "https://openx.solar"
			if websitePing() {
				Return.WebStatus.Text = "Website is Up"
			}
			Return.WebStatus.URL = "openx.solar"
		}(&wgPre)

		wgPre.Wait()

		username := "aibonitoGsIoJ"

		var wg1 sync.WaitGroup
		wg1.Add(1)
		go getUToken(&wg1, username)
		wg1.Add(1)
		go getAToken(&wg1)
		wg1.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var err error
			XlmUSD, err = tickers.BinanceTicker()
			if err != nil {
				log.Fatal(err)
			}
		}(&wg1)
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
		go validateRecp(&wg2, username, Token)
		wg2.Add(1)
		go getInvestor(&wg2, AdminToken, invIndex)
		wg2.Add(1)
		go getDeveloper(&wg2, AdminToken, devIndex)

		wg2.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
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
		}(&wg2)

		wg2.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var projCount length
			data, err := erpc.GetRequest(platformURL + "/admin/getallprojects?username=admin&token=" + AdminToken)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(data, &projCount)
			if err != nil {
				log.Fatal(err)
			}

			Return.ProjCount.Text, err = utils.ToString(projCount.Length)
			if err != nil {
				log.Fatal(err)
			}

			Return.ProjCount.Link = platformURL + "/admin/getallprojects?username=admin&token=" + AdminToken
		}(&wg2)

		wg2.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var userCount length
			data, err := erpc.GetRequest(platformURL + "/admin/getallusers?username=admin&token=" + AdminToken)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(data, &userCount)
			if err != nil {
				log.Fatal(err)
			}

			Return.UserCount.Text, err = utils.ToString(userCount.Length)
			if err != nil {
				log.Fatal(err)
			}

			Return.UserCount.Link = platformURL + "/admin/getallusers?username=admin&token=" + AdminToken
		}(&wg2)

		wg2.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var userCount length
			data, err := erpc.GetRequest(platformURL + "/admin/getallinvestors?username=admin&token=" + AdminToken)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(data, &userCount)
			if err != nil {
				log.Fatal(err)
			}

			Return.InvCount.Text, err = utils.ToString(userCount.Length)
			if err != nil {
				log.Fatal(err)
			}

			Return.InvCount.Link = platformURL + "/admin/getallinvestors?username=admin&token=" + AdminToken
		}(&wg2)

		wg2.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var userCount length
			data, err := erpc.GetRequest(platformURL + "/admin/getallrecipients?username=admin&token=" + AdminToken)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(data, &userCount)
			if err != nil {
				log.Fatal(err)
			}

			Return.RecpCount.Text, err = utils.ToString(userCount.Length)
			if err != nil {
				log.Fatal(err)
			}

			Return.RecpCount.Link = platformURL + "/admin/getallrecipients?username=admin&token=" + AdminToken
		}(&wg2)

		wg2.Wait()

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

		Return.PastEnergyValues = Recipient.PastTellerEnergy
		Return.DeviceLocation = Recipient.DeviceLocation
		Return.DeviceLocation = Recipient.DeviceLocation
		Return.StateHashes = Recipient.StateHashes

		Return.DABalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey
		Return.PBBalance.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

		Return.AccountBalance1.Text = "XLM: " + Pnb + " STABLEUSD: " + Pub
		Return.AccountBalance1.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

		Return.AccountBalance2.Text = "XLM: " + Snb + " STABLEUSD: " + Sub
		Return.AccountBalance2.Link = "https://testnet.steexp.com/account/" + Recipient.U.StellarWallet.PublicKey

		Return.Recipient.Username = Recipient.U.Username
		Return.Recipient.Name = Recipient.U.Name
		Return.Recipient.Email = Recipient.U.Email

		Return.PaybackPeriod = Project.PaybackPeriod
		Return.BalanceLeft = Project.BalLeft
		Return.OwnershipShift = Project.OwnershipShift
		Return.DateInitiated = Project.DateInitiated
		Return.Stage = Project.Stage
		Return.DateFunded = Project.DateFunded
		Return.InvAssetCode = Project.InvestorAssetCode

		Return.Admin.Username = "admin"
		Return.Admin.Password = "password"
		Return.Admin.AdminToken = AdminToken
		Return.Admin.RecpToken = Token

		Return.Date = utils.Timestamp()

		templates.Lookup("doc").Execute(w, Return)
	})
}

func getInvestor(wg *sync.WaitGroup, AdminToken string, invIndex string) {
	defer wg.Done()
	data, err := erpc.GetRequest(platformURL + "/admin/getinvestor?username=admin&token=" +
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
	data, err := erpc.GetRequest(platformURL + "/admin/getentity?username=admin&token=" +
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
