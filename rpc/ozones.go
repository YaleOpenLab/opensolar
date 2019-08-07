package rpc

/*
import (
	"log"
	"net/http"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	opzones "github.com/YaleOpenLab/openx/platforms/ozones"
)

func setupCoopRPCs() {
	getCoopDetails()
	GetAllCoops()
}

func setupBondRPCs() {
	getBondDetails()
	Search()
	GetAllBonds()
}

// GetAllCoops gets a list of all the coops  that are registered on the platform
func GetAllCoops() {
	http.HandleFunc("/coop/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allBonds, err := opzones.RetrieveAllLivingUnitCoops()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, allBonds)
	})
}

// getCoopDetails gets the details of a particular coop
func getCoopDetails() {
	http.HandleFunc("/coop/get", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		uKey, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		bond, err := opzones.RetrieveLivingUnitCoop(uKey)
		if err != nil {
			log.Println("did not retrieve coop", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		erpc.MarshalSend(w, bond)
	})
}

// getBondDetails gets the details of a particular bond
func getBondDetails() {
	http.HandleFunc("/bond/get", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// get the details of a specific bond by key
		if r.URL.Query()["index"] == nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		uKey, err := utils.ToInt(r.URL.Query()["index"][0])
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		bond, err := opzones.RetrieveConstructionBond(uKey)
		if err != nil {
			log.Println("did not retrieve bond", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, bond)
	})
}

// GetAllBonds gets the list of all bonds that are registered on the platfomr
func GetAllBonds() {
	http.HandleFunc("/bond/all", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		allBonds, err := opzones.RetrieveAllConstructionBonds()
		if err != nil {
			log.Println("did not retrieve all bonds", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
		erpc.MarshalSend(w, allBonds)
	})
}

// Search can be used for searching bonds and coops to a limited capacity
func Search() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		erpc.CheckGet(w, r)
		erpc.CheckOrigin(w, r)
		// search for coop / bond  and return accordingly
		if r.URL.Query()["q"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		searchString := r.URL.Query()["q"][0]
		if strings.Contains(searchString, "bond") {
			allBonds, err := opzones.RetrieveAllConstructionBonds()
			if err != nil {
				log.Println("did not retrieve all bonds", err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, allBonds)
			// do bond stuff
		} else if strings.Contains(searchString, "coop") {
			// do coop stuff
			allCoops, err := opzones.RetrieveAllLivingUnitCoops()
			if err != nil {
				log.Println("did not retrieve bond", err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, allCoops)
		} else {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}
	})
}
*/
