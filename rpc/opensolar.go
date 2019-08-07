package rpc

import (
	"errors"
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	platform "github.com/YaleOpenLab/openx/platforms/opensolar"
)

// collect all handlers in one place so that we can assemble them easily
// there are some repeating RPCs that we would like to avoid and maybe there's some
// nice way to group them together

// setupProjectRPCs sets up all the RPC calls related to projects that might be used
func setupProjectRPCs() {
	insertProject()
	getProject()
	getAllProjects()
	getProjectsAtIndex()
}

// parseProject is a helper that is used to validate POST data. This returns a project struct
// on successful parsing of the received form data
func parseProject(r *http.Request) (platform.Project, error) {
	// we need to create an instance of the Project
	// and then map values if they do exist
	// note that we just prepare the project here and don't invest in it
	// for that, we need new a new investor struct and a recipient struct
	var prepProject platform.Project
	err := r.ParseForm()
	if err != nil {
		log.Println("did not parse form", err)
		return prepProject, err
	}
	// if we're inserting this in, we need to get the next index number
	// so that we can set this without causing some weird bugs
	allProjects, err := platform.RetrieveAllProjects()
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
	// this should be a post method since you want to accept an project and then insert
	// that into the database
	// this route does not define an originator and would mostly not be useful, should
	// look into a way where we can define originators in the route as well
	http.HandleFunc("/project/insert", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckPost(w, r)
		erpc.CheckOrigin(w, r)
		var prepProject platform.Project
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

// getAllProjects gets a list of all the projects that registered on the platform.
func getAllProjects() {
	http.HandleFunc("/project/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// make a call to the db to get all projects
		// while making this call, the rpc should not be aware of the db we are using
		// and stuff. So we need to have another route that would open the existing
		// db, without asking for one
		allProjects, err := platform.RetrieveAllProjects()
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
	// we need to read passed the key from the URL that the user calls
	http.HandleFunc("/project/get", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		if r.URL.Query()["index"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		uKey, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		contract, err := platform.RetrieveProject(uKey)
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
	erpc.CheckGet(w, r)
	erpc.CheckOrigin(w, r)
	allProjects, err := platform.RetrieveProjectsAtStage(stage)
	if err != nil {
		log.Println("did not retrieve project at specific stage", err)
		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		return
	}
	erpc.MarshalSend(w, allProjects)
}

// various handlers for fetching projects which are at different stages on the platform
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
