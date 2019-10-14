package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
)

func setupDeveloperRPCs() {
	withdrawdeveloper()
}

var DevRPC = map[int][]string{
	1: []string{"/developer/withdraw", "POST", "amount", "projIndex"}, // GET
}

// getStage1Contracts gets a list of all the originated contracts on the platform
func withdrawdeveloper() {
	http.HandleFunc(DevRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepDev, err := entityValidateHelper(w, r, DevRPC[1][2:], DevRPC[1][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		amountx := r.FormValue("amount")
		projIndexx := r.FormValue("projIndex")

		amount, err := utils.ToFloat(amountx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = core.RequestWaterfallWithdrawal(prepDev.U.Index, projIndex, amount)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
		return
	})
}
