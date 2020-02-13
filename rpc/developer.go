package rpc

import (
	"log"
	"net/http"

	"github.com/YaleOpenLab/opensolar/consts"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	core "github.com/YaleOpenLab/opensolar/core"
)

func setupDeveloperRPCs() {
	withdrawdeveloper()
	developerDashboard()
}

var DevRPC = map[int][]string{
	1: []string{"/developer/withdraw", "POST", "amount", "projIndex"}, // POST
	2: []string{"/developer/dashboard", "GET"},                        // GET
}

// withdrawdeveloper can be called by a developer wishing to withdraw funds from the platfomr
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

type entityDashboardHelper struct {
	YourProfile struct {
		Name           string `json:"Name"`
		ActiveProjects int    `json:"Active Projects"`
	} `json:"Your Profile"`

	YourEnergy struct {
		TiCP    string `json:"Total in Current Period"`
		AllTime string `json:"All Time"`
	} `json:"Your Energy"`

	YourWallet struct {
		ProjectWalletBalance float64 `json:"Project Wallet Balance"`
		AutoReload           string  `json:"Auto Reload"`
	} `json:"Your Wallet"`

	NActions struct {
		Notification    string `json:"Notification"`
		ActionsRequired string `json:"Actions Required"`
	} `json:"Notifications & Actions"`

	YourProjects []entityDashboardData
}

type entityDashboardData struct {
	Index      int
	ExploreTab map[string]interface{} `json:"Explore Tab"`
	Role       string
	PSA        struct {
		Stage   string
		Actions []string
	} `json:"Project Stage & Actions"`
	ProjectWallets struct {
		Wallets [][]string `json:"Project Wallets"`
	}
	PendingPayments []string               `json:"Payments"`
	Documents       map[string]interface{} `json:"Documentation and Smart Contracts"`
}

// developerDashboard returns the stuff that should be there on the contractor dashboard
func developerDashboard() {
	http.HandleFunc(DevRPC[2][0], func(w http.ResponseWriter, r *http.Request) {
		prepEntity, err := entityValidateHelper(w, r, DevRPC[2][2:], DevRPC[2][1])
		if err != nil {
			log.Println("Error while validating entity", err)
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
			return
		}

		var ret entityDashboardHelper

		present := false
		proposed := false

		if len(prepEntity.PresentContractIndices) != 0 {
			present = true
		}
		if len(prepEntity.ProposedContractIndices) != 0 {
			proposed = true
		}

		if !present && !proposed {
			log.Println("Contractor not part of any project")
			erpc.MarshalSend(w, ret)
			return
		}

		var projects []core.Project

		ret.YourProfile.Name = prepEntity.U.Name
		ret.YourProfile.ActiveProjects = len(prepEntity.PresentContractIndices)
		ret.YourEnergy.TiCP = "845 kWh"
		ret.YourEnergy.AllTime = "10,150 MWh"
		ret.YourWallet.AutoReload = "On"
		ret.NActions.Notification = "None"
		ret.NActions.ActionsRequired = "None"

		if present {
			for _, i := range prepEntity.PresentContractIndices {
				project, err := core.RetrieveProject(i)
				if err != nil {
					log.Println("could not retrieve project from db:", err)
					erpc.MarshalSend(w, erpc.StatusInternalServerError)
					return
				}
				projects = append(projects, project)
			}
		}

		if proposed {
			for _, i := range prepEntity.ProposedContractIndices {
				project, err := core.RetrieveProject(i)
				if err != nil {
					log.Println("could not retrieve project from db:", err)
					erpc.MarshalSend(w, erpc.StatusInternalServerError)
					return
				}
				projects = append(projects, project)
			}
		}

		ret.YourProjects = make([]entityDashboardData, len(projects))

		for i, project := range projects {
			var x entityDashboardData
			x.ExploreTab = make(map[string]interface{})
			x.ExploreTab = project.Content.Details["Explore Tab"]
			x.Role = "You are an Offtaker"
			sStage, err := utils.ToString(project.Stage)
			if err != nil {
				log.Println(err)
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.PSA.Stage = "Project is in Stage: " + sStage
			x.PSA.Actions = []string{"Contractor Actions", "No pending action"}
			x.ProjectWallets.Wallets = make([][]string, 2)

			var escrowBalance string
			if consts.Mainnet {
				escrowBalance, err = utils.ToString(xlm.GetAssetBalance(project.EscrowPubkey, consts.AnchorUSDCode))
				if err != nil {
					log.Println(err)
					erpc.MarshalSend(w, erpc.StatusInternalServerError)
					return
				}
				ret.YourWallet.ProjectWalletBalance += xlm.GetAssetBalance(project.EscrowPubkey, consts.AnchorUSDCode)
			} else {

				escrowBalance, err = utils.ToString(xlm.GetAssetBalance(project.EscrowPubkey, consts.StablecoinCode))
				if err != nil {
					log.Println(err)
					erpc.MarshalSend(w, erpc.StatusInternalServerError)
					return
				}
				ret.YourWallet.ProjectWalletBalance += xlm.GetAssetBalance(project.EscrowPubkey, consts.StablecoinCode)
			}

			x.ProjectWallets.Wallets[0] = []string{"Project Escrow Wallet: " + project.EscrowPubkey, escrowBalance}
			x.ProjectWallets.Wallets[1] = []string{"Renewable Energy Certificates (****BBDJL)", "10"}
			x.PendingPayments = []string{"Your Pending Payment", "$203 due on April 30"}
			x.Documents = project.Content.Details["Documents"]

			ret.YourProjects[i] = x
		}
		erpc.MarshalSend(w, ret)
	})
}
