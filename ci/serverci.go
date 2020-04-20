package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	erpc "github.com/Varunram/essentials/rpc"
	// osrpc "github.com/YaleOpenLab/opensolar/rpc"
	"github.com/Varunram/essentials/utils"
	flags "github.com/jessevdk/go-flags"
)

type Content struct {
	Title              string
	Name               string
	PlatformStatus     string
	PlatformStatusLink string
	Fruit              [3]string
}

var platforumURL = "https://api2.openx.solar"

func platPing() bool {
	data, err := erpc.GetRequest(platforumURL + "/ping")
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

		ping := "Platform is Down"
		platformStatusLink := platforumURL + "/ping"
		if platPing() {
			ping = "Platform is Up"
		}

		arr := [3]string{"Apple", "Lemon", "Orange"}
		context := Content{
			Title:              "My Fruits",
			Name:               "John",
			PlatformStatus:     ping,
			PlatformStatusLink: platformStatusLink,
			Fruit:              arr,
		}
		templates.Lookup("doc").Execute(w, context)
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
