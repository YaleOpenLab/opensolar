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
}

// AdminRPC is a list of all the endpoints that can be called by admins
var AdminRPC = map[int][]string{
	1: []string{"/admin/flag", "GET", "projIndex"}, // GET
	2: []string{"/admin/getallprojects", "GET"},    // GET
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
