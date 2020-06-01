package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	core "github.com/YaleOpenLab/opensolar/core"
)

// setupStagesHandlers sets up all stage related handlers
func setupStagesHandlers() {
	returnAllStages()
	returnSpecificStage()
	promoteStage()
}

// StagesRPC is a list of all stage related RPC endpoints
var StagesRPC = map[int][]string{
	1: {"/stages/all", "GET"},              // GET
	2: {"/stages", "GET", "index"},         // GET
	3: {"/stages/promote", "GET", "index"}, // GET
}

// returnAllStages returns all the defined stages for this platform.  Opensolar
// has 9 stages defined in stages.go
func returnAllStages() {
	http.HandleFunc(StagesRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		var arr []core.Stage
		arr = append(arr, core.Stage0, core.Stage1, core.Stage2, core.Stage3, core.Stage4,
			core.Stage5, core.Stage6, core.Stage7, core.Stage8, core.Stage9)

		erpc.MarshalSend(w, arr)
	})
}

// returnSpecificStage returns details on a stage defined in the opensolar platform
func returnSpecificStage() {
	http.HandleFunc(StagesRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, StagesRPC[2][2:], StagesRPC[2][1])
		if err != nil {
			log.Println(err)
			return
		}

		indexx := r.URL.Query()["index"][0]

		index, err := utils.ToInt(indexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "passed index not an integer", messages.ConversionError) {
			return
		}

		var x core.Stage
		switch index {
		case 1:
			x = core.Stage1
		case 2:
			x = core.Stage2
		case 3:
			x = core.Stage3
		case 4:
			x = core.Stage4
		case 5:
			x = core.Stage5
		case 6:
			x = core.Stage6
		case 7:
			x = core.Stage7
		case 8:
			x = core.Stage8
		case 9:
			x = core.Stage9
		default:
			// default is stage0, so we don't have a case defined for it above
			x = core.Stage0
		}

		erpc.MarshalSend(w, x)
	})
}

// promoteStage returns details on a stage defined in the opensolar platform
func promoteStage() {
	http.HandleFunc(StagesRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, StagesRPC[2][2:], StagesRPC[2][1])
		if err != nil {
			log.Println(err)
			return
		}

		indexx := r.URL.Query()["index"][0]
		index, err := utils.ToInt(indexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "passed index not an integer", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(index)
		if err != nil {
			log.Println("Couldn't retrieve project from the database")
			return
		}

		var recpBool, invBool, entityBool bool
		recpBool, invBool, entityBool = true, true, true

		recp, err := recpValidateHelper(w, r, StagesRPC[2][2:], StagesRPC[2][1])
		if err != nil {
			log.Println("stage promoter not a recipient: ", err)
			recpBool = false
		}

		inv, err := InvValidateHelper(w, r, StagesRPC[2][2:], StagesRPC[2][1])
		if err != nil {
			log.Println("stage promoter not an investor: ", err)
			invBool = false
		}

		entity, err := entityValidateHelper(w, r, StagesRPC[2][2:], StagesRPC[2][1])
		if err != nil {
			log.Println("stage promoter not an entity: ", err)
			entityBool = false
		}

		if !(recpBool && recp.U.Index != project.RecipientIndex) || !recp.U.Admin {
			log.Println("not authorized to upgrade stages, quitting")
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if invBool && !inv.U.Admin {
			log.Println("investors are not allowed to perform stage migrations, exiting")
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		if !(entityBool && entity.U.Index != project.ContractorIndex &&
			entity.U.Index != project.GuarantorIndex &&
			entity.U.Index != project.MainDeveloperIndex) || !entity.U.Admin {
			log.Println("you are not a registered entity related to the project")
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		// we check whether the person is actually associated with the project in question
		err = core.StageXtoY(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
