package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"text/template"

	tickers "github.com/Varunram/essentials/exchangetickers"
	"github.com/Varunram/essentials/xlm"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
	flags "github.com/jessevdk/go-flags"
)

func serveStatic() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
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

		username := "fabideasHGJgf"

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
		wg1.Wait()

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

		Return.DeviceID.Text = Recipient.DeviceID

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

func renderHTML() (string, error) {
	doc, err := ioutil.ReadFile("index.html")
	return string(doc), err
}

func startServer(portx int, insecure bool) {
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
	runtime.GOMAXPROCS(4)
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	startServer(opts.Port, opts.Insecure)
}
