package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	// core "github.com/YaleOpenLab/opensolar/core"
)

func setupGuarantorRPCs() {
	depositXLMGuarantor()
	depositAssetGuarantor()
}

var GuaRPC = map[int][]string{
	1: []string{"/guarantor/deposit/xlm", "POST", "amount", "projIndex", "seedpwd"},                // POST
	2: []string{"/guarantor/deposit/asset", "POST", "amount", "projIndex", "seedpwd", "assetCode"}, // POST
}

// depositXLMGuarantor is called by a guarantor when they wish to refill the escrow account with xlm
func depositXLMGuarantor() {
	http.HandleFunc(GuaRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, GuaRPC[1][2:], GuaRPC[1][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		amountx := r.FormValue("amount")
		projIndexx := r.FormValue("projIndex")
		seedpwd := r.FormValue("seedpwd")

		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		amount, err := utils.ToFloat(amountx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = prepEntity.RefillEscrowXLM(projIndex, amount, seedpwd)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// depositAssetGuarantor is called by a guarantor when they wish to refill the escrow account with an asset
func depositAssetGuarantor() {
	http.HandleFunc(GuaRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, GuaRPC[2][2:], GuaRPC[2][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		amountx := r.FormValue("amount")
		projIndexx := r.FormValue("projIndex")
		seedpwd := r.FormValue("seedpwd")
		asset := r.FormValue("assetCode")

		projIndex, err := utils.ToInt(projIndexx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		amount, err := utils.ToFloat(amountx)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = prepEntity.RefillEscrowAsset(projIndex, asset, amount, seedpwd)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
