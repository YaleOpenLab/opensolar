package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/messages"
	"github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	openx "github.com/YaleOpenLab/openx/database"
)

func setupAdminHandlers() {
	flagProject()
}

// AdminRPC is a list of all the endpoints that can be called by admins
var AdminRPC = map[int][]string{
	1: []string{"/admin/flag", "GET", "projIndex"}, // GET
}

// adminValidateHelper is a helper that validates if the caller is an admin, and returns the user struct if so
func adminValidateHelper(w http.ResponseWriter, r *http.Request) (openx.User, error) {
	var user openx.User

	username := r.URL.Query()["username"][0]
	token := r.URL.Query()["token"][0]

	user, err := core.ValidateUser(username, token)
	if err != nil {
		log.Println(err)
		erpc.ResponseHandler(w, erpc.StatusBadRequest, messages.NotAdminError)
		return user, err
	}

	if !user.Admin {
		erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.NotAdminError)
		return user, errors.New("unauthorized")
	}

	return user, nil
}

// flagProject flags a project. Flagging a project stops automated signing by the platform.
func flagProject() {
	http.HandleFunc(AdminRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		err := checkReqdParams(w, r, AdminRPC[1][2:], AdminRPC[1][1])
		if err != nil {
			return
		}

		user, err := adminValidateHelper(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.NotAdminError)
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
