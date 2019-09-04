package rpc

import (
	"errors"
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

// parseProject is a helper that is used to validate POST data. This returns a project struct
// on successful parsing of the received form data
func parseProject(r *http.Request) (core.Project, error) {
	var prepProject core.Project
	err := r.ParseForm()
	if err != nil {
		log.Println("did not parse form", err)
		return prepProject, err
	}

	allProjects, err := core.RetrieveAllProjects()
	if err != nil {
		log.Println("did not retrieve all projects", err)
		return prepProject, errors.New("error in assigning index")
	}
	prepProject.Index = len(allProjects) + 1
	if r.FormValue("PanelSize") == "" || r.FormValue("TotalValue") == "" || r.FormValue("Location") == "" || r.FormValue("Metadata") == "" || r.FormValue("Stage") == "" {
		return prepProject, errors.New("one of given params is missing: PanelSize, TotalValue, Location, Metadata")
	}

	prepProject.PanelSize = r.FormValue("PanelSize")
	prepProject.TotalValue, err = utils.ToFloat(r.FormValue("TotalValue"))
	if err != nil {
		return prepProject, err
	}
	prepProject.State = r.FormValue("Location")
	prepProject.Metadata = r.FormValue("Metadata")
	prepProject.Stage, err = utils.ToInt(r.FormValue("Stage"))
	if err != nil {
		return prepProject, err
	}
	prepProject.MoneyRaised = 0
	prepProject.BalLeft = float64(0)
	prepProject.DateInitiated = utils.Timestamp()
	return prepProject, nil
}

// insertProject inserts a project into the database.
func insertProject() {
	http.HandleFunc("/project/insert", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckPost(w, r)
		erpc.CheckOrigin(w, r)

		var prepProject core.Project
		prepProject, err := parseProject(r)
		if err != nil {
			log.Println("did not parse project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
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
	http.HandleFunc("/project/all", func(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/project/get", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		if r.URL.Query()["index"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		uKey, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		contract, err := core.RetrieveProject(uKey)
		if err != nil {
			log.Println("did not retrieve project", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, contract)
	})
}

// projectHandler gets proejcts at a specific stage from the database
func projectHandler(w http.ResponseWriter, r *http.Request, stage int) {
	err := erpc.CheckGet(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	allProjects, err := core.RetrieveProjectsAtStage(stage)
	if err != nil {
		log.Println("did not retrieve project at specific stage", err)
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		return
	}
	erpc.MarshalSend(w, allProjects)
}

// getProjectsAtIndex gets projects at a specific stage
func getProjectsAtIndex() {
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["index"] == nil {
			log.Println("No stage number passed, not returning anything!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println("Passed index not an integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if index > 9 || index < 0 {
			index = 0
		}

		projectHandler(w, r, index)
	})
}
