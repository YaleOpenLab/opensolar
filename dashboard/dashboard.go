package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

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
	Title         string
	Name          string
	OpensStatus   LinkFormat
	OpenxStatus   LinkFormat
	Validate      LinkFormat
	TellerEnergy  LinkFormat
	DateLastStart LinkFormat
	Fruit         [3]string
}

var platformURL = "https://api2.openx.solar"
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

		x.Validate.Text = val
		x.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + token

		x.TellerEnergy.Text, err = utils.ToString(Recipient.TellerEnergy)
		if err != nil {
			log.Fatal(err)
		}

		x.TellerEnergy.Text = "Energy generated till date: " + x.TellerEnergy.Text + " Wh"

		x.DateLastStart.Text = utils.StringToHumanTime(Recipient.DeviceStarts[len(Recipient.DeviceStarts)-1])
		x.DateLastStart.Text = "Teller Last Start Time: " + x.DateLastStart.Text

		templates.Lookup("doc").Execute(w, x)
	})
}

func renderHTML() (string, error) {
	doc, err := ioutil.ReadFile("index.html")
	return string(doc), err
}

func StartServer(portx int, insecure bool) {
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
