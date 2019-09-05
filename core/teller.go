package core

import (
	"encoding/json"
	"log"
	"time"

	erpc "github.com/Varunram/essentials/rpc"
	//	consts "github.com/YaleOpenLab/opensolar/consts"
	notif "github.com/YaleOpenLab/opensolar/notif"
)

const tellerUrl = "https://localhost"

type statusResponse struct {
	Code   int
	Status string
}

// MonitorTeller monitors a teller and checks whether its live. If not, send an email to platform admins
func MonitorTeller(projIndex int, tellerUrl string) {
	// call this function only after a specific order has been accepted by the recipient
	log.Println("monitoring the teller")
	for {
		project, err := RetrieveProject(projIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		data, err := erpc.GetRequest(tellerUrl + "/ping")
		if err != nil {
			log.Println("did not create new GET request", err)
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(5 * time.Second)
			continue
		}

		var x statusResponse
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("error while unmarshalling data", err, string(data))
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
			time.Sleep(5 * time.Second)
			continue
		}

		if x.Code != 200 || x.Status != "HEALTH OK" {
			notif.SendTellerDownEmail(project.Index, project.RecipientIndex)
		}

		time.Sleep(5 * time.Second)
	}
}
