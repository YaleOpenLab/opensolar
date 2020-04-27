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
}

// AdminRPC is a list of all the endpoints that can be called by admins
var AdminRPC = map[int][]string{
	1: []string{"/admin/flag", "GET", "projIndex"},     // GET
	2: []string{"/admin/getallprojects", "GET"},        // GET
	3: []string{"/admin/getrecipient", "GET", "index"}, // GET
	4: []string{"/admin/getinvestor", "GET", "index"},  // GET
	5: []string{"/admin/getentity", "GET", "index"},    // GET
	6: []string{"/admin/getallinvestors", "GET"},       // GET
	7: []string{"/admin/getallrecipients", "GET"},      // GET
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.ParamError("projIndex"))
			return
		}

		err = core.MarkFlagged(projIndex, user.Index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		x, err := core.RetrieveRecipient(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		x, err := core.RetrieveInvestor(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		x, err := core.RetrieveEntity(index)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
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
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x LenReturn
		x.Length = len(recipients)

		erpc.MarshalSend(w, x)
	})
}
