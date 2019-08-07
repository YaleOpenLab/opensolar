package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
)

func setupStagesHandlers() {
	returnAllStages()
	returnSpecificStage()
	promoteStage()
}

// returnAllStages returns all the defined stages for this specific platform.
// Opensolar has 9 stages defined in stages.go
// this is a public function that can be called by anyone, so we don't authenticate
func returnAllStages() {
	http.HandleFunc("/stages/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		var arr []opensolar.Stage
		arr = append(arr, opensolar.Stage0, opensolar.Stage1, opensolar.Stage2, opensolar.Stage3, opensolar.Stage4,
			opensolar.Stage5, opensolar.Stage6, opensolar.Stage7, opensolar.Stage8, opensolar.Stage9)

		erpc.MarshalSend(w, arr)
	})
}

// returnSpecificStage returns details on a specific stage defined in the opensolar platform
func returnSpecificStage() {
	http.HandleFunc("/stages", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["index"] == nil {
			log.Println("User did not pass index to retrieve stage for, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println("Passed index not an integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		var x opensolar.Stage
		switch index {
		case 1:
			x = opensolar.Stage1
		case 2:
			x = opensolar.Stage2
		case 3:
			x = opensolar.Stage3
		case 4:
			x = opensolar.Stage4
		case 5:
			x = opensolar.Stage5
		case 6:
			x = opensolar.Stage6
		case 7:
			x = opensolar.Stage7
		case 8:
			x = opensolar.Stage8
		case 9:
			x = opensolar.Stage9
		default:
			// default is stage0, so we don't have a case defined for it above
			x = opensolar.Stage0
		}

		erpc.MarshalSend(w, x)
	})
}

// promoteStage returns details on a specific stage defined in the opensolar platform
func promoteStage() {
	http.HandleFunc("/stages/promote", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)

		if r.URL.Query()["index"] == nil {
			log.Println("some fields missing to promote from stage x to y, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		err = opensolar.StageXtoY(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
