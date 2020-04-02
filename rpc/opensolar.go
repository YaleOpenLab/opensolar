package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	core "github.com/YaleOpenLab/opensolar/core"
	notif "github.com/YaleOpenLab/opensolar/notif"
)

// setupProjectRPCs sets up all project related RPC calls
func setupProjectRPCs() {
	getProject()
	getAllProjects()
	getProjectsAtIndex()
	addContractHash()
	sendTellerShutdownEmail()
	sendTellerFailedPaybackEmail()
	explore()
	projectDetail()
}

// ProjectRPC contains a list of all the project related RPC endpoints
var ProjectRPC = map[int][]string{
	2:  []string{"/project/all", "GET"},                                           // GET
	3:  []string{"/project/get", "GET", "index"},                                  // GET
	4:  []string{"/projects", "GET", "stage"},                                     // GET
	5:  []string{"/utils/addhash", "GET", "projIndex", "choice", "choicestr"},     // GET
	6:  []string{"/tellershutdown", "GET", "projIndex", "deviceId", "tx1", "tx2"}, // GET
	7:  []string{"/tellerpayback", "GET", "deviceId", "projIndex"},                // GET
	8:  []string{"/project/get/dashboard", "GET", "index"},                        // GET
	9:  []string{"/explore", "GET"},                                               // GET
	10: []string{"/project/detail", "GET", "index"},                               // GET
}

// getAllProjects gets a list of all projects
func getAllProjects() {
	http.HandleFunc(ProjectRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
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

// getProject retrieves details of a project when passed the index
func getProject() {
	http.HandleFunc(ProjectRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["index"] == nil {
			log.Println("index not passed")
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ParamError("index"))
			return
		}

		err := erpc.CheckGet(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusNotFound)
			return
		}

		index := r.URL.Query()["index"][0]

		uKey, err := utils.ToInt(index)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
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

// getProjectsAtIndex gets projects belonging to one stage
func getProjectsAtIndex() {
	http.HandleFunc(ProjectRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		err = checkReqdParams(w, r, ProjectRPC[4][2:], ProjectRPC[4][1])
		if err != nil {
			log.Println(err)
			return
		}

		stagex := r.URL.Query()["stage"][0]

		stage, err := utils.ToInt(stagex)
		if err != nil {
			log.Println("Passed index not an integer, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
			return
		}

		if stage > 9 || stage < 0 {
			stage = 0
		}

		allProjects, err := core.RetrieveProjectsAtStage(stage)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, allProjects)
	})
}

// addContractHash adds a contract hash to a project's stage data
func addContractHash() {
	http.HandleFunc(ProjectRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = userValidateHelper(w, r, ProjectRPC[5][2:], ProjectRPC[5][1])
		if err != nil {
			return
		}

		choice := r.URL.Query()["choice"][0]
		hashString := r.URL.Query()["choicestr"][0]
		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if err != nil {
			log.Println("passed project index not int, quitting!")
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
			return
		}

		project, err := core.RetrieveProject(projIndex)
		if err != nil {
			log.Println("couldn't retrieve prject index from database")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		// there are in total 5 types of hashes: OriginatorMoUHash, ContractorContractHash,
		// InvPlatformContractHash, RecPlatformContractHash, SpecSheetHash
		// lets have a fixed set of strings that we can map on here so we have a single endpoint for storing all these hashes

		// TODO: read from the pending docs map here and store this only if we need to.
		switch choice {
		case "omh":
			if project.Stage == 0 {
				project.StageData = append(project.StageData, hashString)
			}
		case "cch":
			if project.Stage == 2 {
				project.StageData = append(project.StageData, hashString)
			}
		case "ipch":
			if project.Stage == 4 {
				project.StageData = append(project.StageData, hashString)
			}
		case "rpch":
			if project.Stage == 4 {
				project.StageData = append(project.StageData, hashString)
			}
		case "ssh":
			if project.Stage == 5 {
				project.StageData = append(project.StageData, hashString)
			}
		default:
			log.Println("invalid choice passed, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = project.Save()
		if err != nil {
			log.Println("error while saving project to db, quitting!")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendTellerShutdownEmail sends a teller shutdown email
func sendTellerShutdownEmail() {
	http.HandleFunc(ProjectRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepUser, err := userValidateHelper(w, r, ProjectRPC[6][2:], ProjectRPC[6][1])
		if err != nil {
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		tx1 := r.URL.Query()["tx1"][0]
		tx2 := r.URL.Query()["tx2"][0]
		notif.SendTellerShutdownEmail(prepUser.Email, projIndex, deviceId, tx1, tx2)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// sendTellerFailedPaybackEmail sends a teller failed payback email
func sendTellerFailedPaybackEmail() {
	http.HandleFunc(ProjectRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		prepUser, err := userValidateHelper(w, r, ProjectRPC[7][2:], ProjectRPC[7][1])
		if err != nil {
			log.Println(err)
			return
		}

		projIndex := r.URL.Query()["projIndex"][0]
		deviceId := r.URL.Query()["deviceId"][0]
		notif.SendTellerPaymentFailedEmail(prepUser.Email, projIndex, deviceId)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// getProjectDashboard gets the project details stub that is displayed on the
// explore page of the frontend
func getProjectDashboard() {
	http.HandleFunc(ProjectRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		err := checkReqdParams(w, r, ProjectRPC[8][2:], ProjectRPC[8][1])
		if err != nil {
			log.Println(err)
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
			return
		}

		project, err := core.RetrieveProject(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, project)
	})
}

// ExplorePageStub is used to show brief descriptions of the project on
// the frontend's explore page
type ExplorePageStub struct {
	StageDescription interface{}
	Name             interface{}
	Location         interface{}
	ProjectType      interface{}
	OriginatorName   interface{}
	Description      interface{}
	Bullet1          interface{}
	Bullet2          interface{}
	Bullet3          interface{}
	Solar            interface{}
	Storage          interface{}
	Tariff           interface{}
	Return           interface{}
	Rating           interface{}
	Tax              interface{}
	Acquisition      interface{}
	MainImage        interface{}
	Stage            int
	Index            int
	Raised           float64
	Total            float64
	Backers          int
}

// explore is the endpoint called on the frontend to show a comprehensive
// description of all projects
func explore() {
	http.HandleFunc(ProjectRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		allProjects, err := core.RetrieveAllProjects()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var arr []ExplorePageStub
		for _, project := range allProjects {
			var x ExplorePageStub
			x.StageDescription = project.Content.Details["Explore Tab"]["stage description"]
			x.Stage = project.Stage
			x.Name = project.Content.Details["Explore Tab"]["name"]
			x.Location = project.Content.Details["Explore Tab"]["location"]
			x.ProjectType = project.Content.Details["Explore Tab"]["project type"]
			x.OriginatorName = project.Content.Details["Explore Tab"]["originator name"]
			x.Description = project.Content.Details["Explore Tab"]["description"]
			x.Bullet1 = project.Content.Details["Explore Tab"]["bullet 1"]
			x.Bullet2 = project.Content.Details["Explore Tab"]["bullet 2"]
			x.Bullet3 = project.Content.Details["Explore Tab"]["bullet 3"]
			x.Solar = project.Content.Details["Explore Tab"]["solar"]
			x.Storage = project.Content.Details["Other Details"]["storage"]
			x.Tariff = project.Content.Details["Other Details"]["tariff"]
			x.Return = project.Content.Details["Explore Tab"]["return"]
			x.Rating = project.Content.Details["Explore Tab"]["rating"]
			x.Tax = project.Content.Details["Other Details"]["tax"]
			x.MainImage = project.Content.Details["Explore Tab"]["mainimage"]
			x.Acquisition = project.Content.Details["Explore Tab"]["acquisition"]
			x.Index = project.Index
			x.Raised = project.MoneyRaised
			x.Total = project.TotalValue
			x.Backers = len(project.InvestorMap)
			arr = append(arr, x)
		}

		erpc.MarshalSend(w, arr)
	})
}

// projectDetail is an endpoint that fetches details required by the frontend
func projectDetail() {
	http.HandleFunc(ProjectRPC[10][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["index"] == nil {
			log.Println("project index not passed")
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ParamError("index"))
			return
		}

		index, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ConversionError)
			return
		}

		project, err := core.RetrieveProject(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		stageString, err := utils.ToString(project.Stage)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError, messages.ConversionError)
			return
		}

		project.Content.Details["ExploreTab"]["StageDescription"] = stageString + " | " + core.GetStageDescription(project.Stage)
		project.Content.Details["ExploreTab"]["MoneyRaised"] = project.MoneyRaised
		project.Content.Details["ExploreTab"]["TotalValue"] = project.TotalValue

		erpc.MarshalSend(w, project.Content.Details)
	})
}
