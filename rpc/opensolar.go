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
	getActiveProjects()
	getCompletedProjects()
	getFeaturedProjects()
}

// ProjectRPC contains a list of all the project related RPC endpoints
var ProjectRPC = map[int][]string{
	2:  {"/project/all", "GET"},                                           // GET
	3:  {"/project/get", "GET", "index"},                                  // GET
	4:  {"/projects", "GET", "stage"},                                     // GET
	5:  {"/utils/addhash", "GET", "projIndex", "choice", "choicestr"},     // GET
	6:  {"/tellershutdown", "GET", "projIndex", "deviceId", "tx1", "tx2"}, // GET
	7:  {"/tellerpayback", "GET", "deviceId", "projIndex"},                // GET
	8:  {"/project/get/dashboard", "GET", "index"},                        // GET
	9:  {"/explore", "GET"},                                               // GET
	10: {"/project/detail", "GET", "index"},                               // GET
	11: {"/project/active", "GET"},                                        // GET
	12: {"/project/complete", "GET"},                                      // GET
	13: {"/project/featured", "GET"},                                      // GET
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
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve all projects") {
			return
		}
		erpc.MarshalSend(w, allProjects)
	})
}

// getActiveProjects gets a list of active projects
func getActiveProjects() {
	http.HandleFunc(ProjectRPC[11][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		activeProjects, err := core.RetrieveActiveProjects()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve all projects") {
			return
		}

		erpc.MarshalSend(w, activeProjects)
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
		if erpc.Err(w, err, erpc.StatusNotFound) {
			return
		}

		index := r.URL.Query()["index"][0]

		uKey, err := utils.ToInt(index)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}
		contract, err := core.RetrieveProject(uKey)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
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
		if erpc.Err(w, err, erpc.StatusBadRequest, "Passed index not an integer, quitting", messages.ConversionError) {
			return
		}

		if stage > 9 || stage < 0 {
			stage = 0
		}

		allProjects, err := core.RetrieveProjectsAtStage(stage)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
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
		if erpc.Err(w, err, erpc.StatusBadRequest, "passed project index not int, quitting", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(projIndex)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "couldn't retrieve prject index from database") {
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
		if erpc.Err(w, err, erpc.StatusInternalServerError, "error while saving project to db, quitting!") {
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
		deviceID := r.URL.Query()["deviceId"][0]
		tx1 := r.URL.Query()["tx1"][0]
		tx2 := r.URL.Query()["tx2"][0]
		notif.SendTellerShutdownEmail(prepUser.Email, projIndex, deviceID, tx1, tx2)
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
		deviceID := r.URL.Query()["deviceId"][0]
		notif.SendTellerPaymentFailedEmail(prepUser.Email, projIndex, deviceID)
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
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
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
		//log.Println("All demo projects", allProjects)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var arr []ExplorePageStub
		for _, project := range allProjects {

			if project.Complete {
				continue
			}

			var x ExplorePageStub

			stageString, err := utils.ToString(project.Stage)
			if erpc.Err(w, err, erpc.StatusInternalServerError) {
				return
			}

			x.StageDescription = stageString + " | " + core.GetStageDescription(project.Stage)
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
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		project, err := core.RetrieveProject(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		stageString, err := utils.ToString(project.Stage)
		if erpc.Err(w, err, erpc.StatusInternalServerError, "", messages.ConversionError) {
			return
		}

		project.Content.Details["ExploreTab"]["StageDescription"] = stageString + " | " + core.GetStageDescription(project.Stage)
		project.Content.Details["ExploreTab"]["MoneyRaised"] = project.MoneyRaised
		project.Content.Details["ExploreTab"]["TotalValue"] = project.TotalValue

		erpc.MarshalSend(w, project.Content.Details)
	})
}

// getCompletedProjects gets a list of active projects
func getCompletedProjects() {
	http.HandleFunc(ProjectRPC[12][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		activeProjects, err := core.RetrieveCompletedProjects()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve all projects") {
			return
		}

		erpc.MarshalSend(w, activeProjects)
	})
}

// getFeaturedProjects gets a list of featured projects
func getFeaturedProjects() {
	http.HandleFunc(ProjectRPC[13][0], func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		activeProjects, err := core.RetrieveFeaturedProjects()
		if erpc.Err(w, err, erpc.StatusInternalServerError, "did not retrieve all projects") {
			return
		}

		erpc.MarshalSend(w, activeProjects)
	})
}
