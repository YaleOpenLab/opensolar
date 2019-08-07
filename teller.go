package opensolar

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	consts "github.com/YaleOpenLab/openx/consts"
	notif "github.com/YaleOpenLab/openx/notif"
)

const tellerUrl = "https://localhost"

type statusResponse struct {
	Code   int
	Status string
}

// MonitorTeller monitors a teller and checks whether its live. If not, send an email to platform admins
func MonitorTeller(projIndex int) {
	// call this function only after a specific order has been accepted by the recipient
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("GET", tellerUrl+"/ping", nil)
		if err != nil {
			log.Println("did not create new GET request", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval)
			continue
		}

		req.Header.Set("Origin", "localhost")
		res, err := client.Do(req)
		if err != nil {
			log.Println("did not make request", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval)
			continue
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error while reading response body", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval)
			continue
		}

		var x statusResponse
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("error while unmarshalling data", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(consts.TellerPollInterval)
			continue
		}

		if x.Code != 200 || x.Status != "HEALTH OK" {
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
		}

		res.Body.Close()
		time.Sleep(consts.TellerPollInterval)
	}
}
