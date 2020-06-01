package rpc

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"

	"github.com/YaleOpenLab/opensolar/messages"
	// core "github.com/YaleOpenLab/opensolar/core"
)

func setupGuarantorRPCs() {
	depositXLMGuarantor()
	depositAssetGuarantor()
}

// GuaRPC contains a list of all guarantor related RPC endpoints
var GuaRPC = map[int][]string{
	1: {"/guarantor/deposit/xlm", "POST", "amount", "projIndex", "seedpwd"},                // POST
	2: {"/guarantor/deposit/asset", "POST", "amount", "projIndex", "seedpwd", "assetCode"}, // POST
}

// depositXLMGuarantor is called by a guarantor when they wish to refill
// the escrow account with xlm
func depositXLMGuarantor() {
	http.HandleFunc(GuaRPC[1][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, GuaRPC[1][2:], GuaRPC[1][1])
		if err == nil {
			if !prepEntity.Guarantor {
				erpc.ResponseHandler(w, erpc.StatusUnauthorized, messages.NotGuarantorError)
				return
			}
		} else {
			log.Println("Error while validating entity", err)
			return
		}

		amountx := r.FormValue("amount")
		projIndexx := r.FormValue("projIndex")
		seedpwd := r.FormValue("seedpwd")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		amount, err := utils.ToFloat(amountx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		err = prepEntity.RefillEscrowXLM(projIndex, amount, seedpwd)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// depositAssetGuarantor is called by a guarantor when they wish to refill
// the escrow account with an asset
func depositAssetGuarantor() {
	http.HandleFunc(GuaRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, GuaRPC[2][2:], GuaRPC[2][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			return
		}

		amountx := r.FormValue("amount")
		projIndexx := r.FormValue("projIndex")
		seedpwd := r.FormValue("seedpwd")
		asset := r.FormValue("assetCode")

		projIndex, err := utils.ToInt(projIndexx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		amount, err := utils.ToFloat(amountx)
		if erpc.Err(w, err, erpc.StatusBadRequest, "", messages.ConversionError) {
			return
		}

		err = prepEntity.RefillEscrowAsset(projIndex, asset, amount, seedpwd)
		if erpc.Err(w, err, erpc.StatusInternalServerError) {
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}
