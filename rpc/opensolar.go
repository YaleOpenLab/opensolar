package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	core "github.com/YaleOpenLab/opensolar/core"
)

// setupProjectRPCs sets up all project related RPC calls
func setupProjectRPCs() {
	insertProject()
	getProject()
	getAllProjects()
	getProjectsAtIndex()
}

var ProjRpc = map[int][]string{
	1: []string{"/project/insert", "PanelSize", "TotalValue", "Location", "Metadata", "Stage"}, // POST
	2: []string{"/project/all"},                                                                // GET
	3: []string{"/project/get", "index"},                                                       // GET
	4: []string{"/projects", "index"},                                                          // GET
}

// insertProject inserts a project into the database.
func insertProject() {
	http.HandleFunc(ProjRpc[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ProjRpc[1][1:])
		if err != nil {
			log.Println(err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		panelSize := r.FormValue("PanelSize")
		totalValue := r.FormValue("TotalValue")
		location := r.FormValue("Location")
		metadata := r.FormValue("Metadata")
		stage := r.FormValue("Stage")

		allProjects, err := core.RetrieveAllProjects()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		var prepProject core.Project

		prepProject.Index = len(allProjects) + 1
		prepProject.PanelSize = panelSize
		prepProject.TotalValue, err = utils.ToFloat(totalValue)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		prepProject.State = location
		prepProject.Metadata = metadata
		prepProject.Stage, err = utils.ToInt(stage)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		prepProject.MoneyRaised = 0
		prepProject.BalLeft = float64(0)
		prepProject.DateInitiated = utils.Timestamp()

		err = prepProject.Save()
		if err != nil {
			log.Println("did not save project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// getAllProjects gets a list of all the projects in the database
func getAllProjects() {
	http.HandleFunc(ProjRpc[2][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		allProjects, err := core.RetrieveAllProjects()
		if err != nil {
			log.Println("did not retrieve all projects", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, allProjects)
	})
}

// getProject gets the details of a specific project.
func getProject() {
	http.HandleFunc(ProjRpc[3][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ProjRpc[3][1:])
		if err != nil {
			log.Println(err)
			return
		}

		index := r.URL.Query()["index"][0]

		uKey, err := utils.ToInt(index)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		contract, err := core.RetrieveProject(uKey)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, contract)
	})
}

// projectHandler gets proejcts at a specific stage from the database
func projectHandler(w http.ResponseWriter, r *http.Request, stage int) {
}

// getProjectsAtIndex gets projects at a specific stage
func getProjectsAtIndex() {
	http.HandleFunc(ProjRpc[4][0], func(w http.ResponseWriter, r *http.Request) {

		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ProjRpc[4][1:])
		if err != nil {
			log.Println(err)
			return
		}

		indexx := r.URL.Query()["index"][0]

		index, err := utils.ToInt(indexx)
		if err != nil {
			log.Println("Passed index not an integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if index > 9 || index < 0 {
			index = 0
		}

		allProjects, err := core.RetrieveProjectsAtStage(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, allProjects)
	})
}
