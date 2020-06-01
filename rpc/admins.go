package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	"github.com/YaleOpenLab/openx/database"
)

func setupAdminHandlers() {
	flagProject()
	getallProjectsAdmin()
	retrieveRecpAdmin()
	retrieveInvAdmin()
	retrieveEntityAdmin()
	retrieveAllInvestors()
	retrieveAllRecipients()
	projectComplete()
	projectFeatured()
}

// AdminRPC is a list of all the endpoints that can be called by admins
var AdminRPC = map[int][]string{
	1: {"/admin/flag", "GET", "projIndex"},          // GET
	2: {"/admin/getallprojects", "GET"},             // GET
	3: {"/admin/getrecipient", "GET", "index"},      // GET
	4: {"/admin/getinvestor", "GET", "index"},       // GET
	5: {"/admin/getentity", "GET", "index"},         // GET
	6: {"/admin/getallinvestors", "GET"},            // GET
	7: {"/admin/getallrecipients", "GET"},           // GET
	8: {"/admin/project/complete", "POST", "index"}, // POST
	9: {"/admin/project/featured", "POST", "index"}, // POST
}

// validateAdmin validates whether a given user is an admin and returns a bool
func validateAdmin(w http.ResponseWriter, r *http.Request, options []string, method string) (database.User, bool) {
	prepUser, err := userValidateHelper(w, r, options, method)
	if err != nil {
		log.Println(err)
		return prepUser, false
	}

	if !prepUser.Admin {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		return prepUser, false
	}

	return prepUser, true
}

// flagProject flags a project. Flagging a project stops automated signing by the platform.
func flagProject() {
	http.HandleFunc(AdminRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		user, admin := validateAdmin(w, r, AdminRPC[1][2:], AdminRPC[1][1])
		if !admin {
			return
		}

		projIndex, err := utils.ToInt(r.URL.Query()["projIndex"][0])
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ParamError("projIndex")) {
			return
		}

		err = core.MarkFlagged(projIndex, user.Index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

type getProjectsAdmin struct {
	Length int
}

func getallProjectsAdmin() {
	http.HandleFunc(AdminRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[2][2:], AdminRPC[2][1])
		if !adminBool {
			return
		}

		projects, err := core.RetrieveAllProjects()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x getProjectsAdmin
		x.Length = len(projects)

		erpc.MarshalSend(w, x)
	})
}

func retrieveRecpAdmin() {
	http.HandleFunc(AdminRPC[3][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[3][2:], AdminRPC[3][1])
		if !adminBool {
			return
		}

		indexS := r.URL.Query()["index"][0]

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		x, err := core.RetrieveRecipient(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func retrieveInvAdmin() {
	http.HandleFunc(AdminRPC[4][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[4][2:], AdminRPC[4][1])
		if !adminBool {
			return
		}

		indexS := r.URL.Query()["index"][0]

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		x, err := core.RetrieveInvestor(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func retrieveEntityAdmin() {
	http.HandleFunc(AdminRPC[5][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[5][2:], AdminRPC[5][1])
		if !adminBool {
			return
		}

		indexS := r.URL.Query()["index"][0]

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		x, err := core.RetrieveEntity(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.MarshalSend(w, x)
	})
}

// LenReturn is the length return structure
type LenReturn struct {
	Length int
}

func retrieveAllInvestors() {
	http.HandleFunc(AdminRPC[6][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[6][2:], AdminRPC[6][1])
		if !adminBool {
			return
		}

		investors, err := core.RetrieveAllInvestors()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x LenReturn
		x.Length = len(investors)

		erpc.MarshalSend(w, x)
	})
}

func retrieveAllRecipients() {
	http.HandleFunc(AdminRPC[7][0], func(w http.ResponseWriter, r *http.Request) {
		_, adminBool := validateAdmin(w, r, AdminRPC[7][2:], AdminRPC[7][1])
		if !adminBool {
			return
		}

		recipients, err := core.RetrieveAllRecipients()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		var x LenReturn
		x.Length = len(recipients)

		erpc.MarshalSend(w, x)
	})
}

// projectComplete marks a project as completed. Should be manually set by admins from the admin dashboard
func projectComplete() {
	http.HandleFunc(AdminRPC[8][0], func(w http.ResponseWriter, r *http.Request) {
		user, admin := validateAdmin(w, r, AdminRPC[8][2:], AdminRPC[8][1])
		if !admin {
			return
		}

		indexS := r.FormValue("index")

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		project, err := core.RetrieveProject(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		project.Complete = true
		project.CompleteAuth = user.Index
		project.CompleteDate = utils.Timestamp()

		err = project.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// projectFeatured sets the featured flag on a project
func projectFeatured() {
	http.HandleFunc(AdminRPC[9][0], func(w http.ResponseWriter, r *http.Request) {
		_, admin := validateAdmin(w, r, AdminRPC[9][2:], AdminRPC[9][1])
		if !admin {
			return
		}

		indexS := r.FormValue("index")

		index, err := utils.ToInt(indexS)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		project, err := core.RetrieveProject(index)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		project.Featured = true

		err = project.Save()
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
