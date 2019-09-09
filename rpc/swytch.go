package rpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

func setupSwytchApis() {
	getAccessToken()
	getRefreshToken()
	getSwytchUser()
	getAssets()
	getEnergy()
	getEnergyAttribution()
}

type getAccessTokenDataHelper struct {
	Accesstoken  string `json:"access_token"`
	Issuedat     int64  `json:"issued_at"`
	Refreshtoken string `json:"refresh_token"`
	Tokentype    string `json:"token_type"`
	Expiresin    int64  `json:"expires_in"`
}

// GetAccessTokenData is a helper struct for the swytch API
type GetAccessTokenData struct {
	Data []getAccessTokenDataHelper `json:"data"`
}

func getAccessToken() {
	http.HandleFunc("/swytch/accessToken", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["clientId"] == nil || r.URL.Query()["clientSecret"] == nil ||
			r.URL.Query()["username"] == nil || r.URL.Query()["password"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
		url := "https://platformapi-staging.swytch.io/v1/oauth/token"
		pwd := "password"

		clientId := r.URL.Query()["clientId"][0]
		clientSecret := r.URL.Query()["clientSecret"][0]
		username := r.URL.Query()["username"][0]
		password := r.URL.Query()["password"][0]

		a := `{
			"grant_type":"` + pwd + `",
			"client_id":"` + clientId + `",
			"client_secret":"` + clientSecret + `",
			"username":"` + username + `",
			"password":"` + password + `"
		}`
		log.Println(a)
		reqbody := strings.NewReader(a)
		req, err := http.NewRequest("POST", url, reqbody)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, _ := ioutil.ReadAll(res.Body)

		var x GetAccessTokenData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func getRefreshToken() {
	http.HandleFunc("/swytch/refreshToken", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["clientId"] == nil || r.URL.Query()["clientSecret"] == nil ||
			r.URL.Query()["refreshToken"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		url := "https://platformapi-staging.swytch.io/v1/oauth/token"
		pwd := "refresh_token"

		clientId := r.URL.Query()["clientId"][0]
		clientSecret := r.URL.Query()["clientSecret"][0]
		refreshToken := r.URL.Query()["refreshToken"][0]

		a := `
		{
			"grant_type":"` + pwd + `",
			"client_id":"` + clientId + `",
			"client_secret":"` + clientSecret + `",
			"refresh_token": "` + refreshToken + `"
		}`

		reqbody := strings.NewReader(a)
		req, err := http.NewRequest("POST", url, reqbody)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, _ := ioutil.ReadAll(res.Body)

		var x GetAccessTokenData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

type getSwytchUserStructToken struct {
	Validtokenbalance bool   `json:"valid_token_balance"`
	Checkedon         string `json:"checked_on"`
	Tokenhash         string `json:"token_hash"`
}

type getSwytchUserStructHelper struct {
	Id           string                   `json:"id"`
	Firstname    string                   `json:"first_name"`
	Lastname     string                   `json:"last_name"`
	Name         string                   `json:"name"`
	Email        string                   `json:"email"`
	Username     string                   `json:"username"`
	Roles        []string                 `json:"roles"`
	Tokenstaking getSwytchUserStructToken `json:"token_staking"`
	Wallet       string                   `json:"wallet"`
}

// GetSwytchUserStruct is a helper struct for the swytch API
type GetSwytchUserStruct struct {
	Data []getSwytchUserStructHelper `json:"data"`
}

func getSwytchUser() {
	http.HandleFunc("/swytch/getuser", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["authToken"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		url := "https://platformapi-staging.swytch.io/v1/auth/user"
		authToken := r.URL.Query()["authToken"][0]

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+authToken)
		req.Header.Add("cache-control", "no-cache")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		log.Println(res)

		var x GetSwytchUserStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type gA2Meta struct {
	Manufacturer    string  `json:"manufacturer"`
	NameplateRating float64 `json:"nameplateRating"`
	SerialNO        string  `json:"serialNO"`
	ThingName       string  `json:"thingName"`
	ThingArn        string  `json:"thingArn"`
	ThingId         string  `json:"thingId"`
}

type gA2Position struct {
	Id          string    `json:"_id"`
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type gA2 struct {
	Id         string      `json:"_id"`
	UpdatedAt  string      `json:"updatedAt"`
	CreatedAt  string      `json:"createdAt"`
	Position   gA2Position `json:"position"`
	Arn        string      `json:"arn"`
	Assetid    string      `json:"asset_id"`
	Ownerid    string      `json:"owner_id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Location   string      `json:"location"`
	Meta       gA2Meta     `json:"meta"`
	Country    string      `json:"country"`
	Status     string      `json:"status"`
	Generating bool        `json:"generating"`
	Nodetype   string      `json:"node_type"`
}

// GetAssetStruct is a helper struct for the swytch API
type GetAssetStruct struct {
	Data []gA2 `json:"data"`
}

func getAssets() {
	http.HandleFunc("/swytch/getassets", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["userId"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		authToken := r.URL.Query()["authToken"][0]
		userId := r.URL.Query()["userId"][0]

		url := "https://platformapi-staging.swytch.io/v1/users/" + userId + "/assets"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+authToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		log.Println(string(body))
		var x GetAssetStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type gEMetadata struct {
	Status          string  `json:"status"`
	Elevation       string  `json:"elevation"`
	Longitude       string  `json:"longitude"`
	Latitude        string  `json:"latitude"`
	Field8          string  `json:"field8"`
	Field7          string  `json:"field7"`
	Field6          string  `json:"field6"`
	Field5          string  `json:"field5"`
	Field4          float64 `json:"field4"`
	Field3          string  `json:"field3"`
	Field2          string  `json:"field2"`
	Field1          string  `json:"field1"`
	Entryid         float64 `json:"entry_id"`
	Createdat       string  `json:"created_at"`
	Manufacturer    string  `json:"manufacturer"`
	NameplateRating float64 `json:"nameplateRating"`
	SerialNO        string  `json:"serialNO"`
	ThingName       string  `json:"thingName"`
	ThingArn        string  `json:"thingArn"`
	ThingId         string  `json:"thingId"`
	Sourcetimestamp string  `json:"source_timestamp"`
}

type getEnergyHelper struct {
	Id              string     `json:"_id"`
	Assetid         string     `json:"asset_id"`
	Assettype       string     `json:"asset_type"`
	Source          string     `json:"source"`
	Value           float64    `json:"value"`
	Unit            string     `json:"unit"`
	Lat             string     `json:"lat"`
	Lng             string     `json:"lng"`
	Energytimestamp string     `json:"energy_timestamp"`
	Timestamp       string     `json:"timestamp"`
	Metadata        gEMetadata `json:"meta"`
	Hash            string     `json:"hash"`
	Blockid         string     `json:"block_id"`
	Blockhash       string     `json:"block_hash"`
	Blocktime       string     `json:"block_time"`
	CreatedAt       string     `json:"createdAt"`
	UpdatedAt       string     `json:"updatedAt"`
}

// GetEnergyStruct is a helper struct for the swytch API
type GetEnergyStruct struct {
	Data []getEnergyHelper `json:"data"`
}

func getEnergy() {
	http.HandleFunc("/swytch/getenergy", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["assetId"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		authToken := r.URL.Query()["authToken"][0]
		assetId := r.URL.Query()["assetId"][0]

		url := "https://platformapi-staging.swytch.io/v1/assets/" + assetId + "/energy?limit=100&offset=0"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+authToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		var x GetEnergyStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type getEnergyAttributionOrigin struct {
	Id          string    `json:"_id"`
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type getEnergyAttributionI struct {
	NextTokenFactorEffectiveAfter string
	NextTokenFactor               string
	CurrentTokenFactor            float64
}

type getEnergyAttributionInputs struct {
	TokenFactor   getEnergyAttributionI `json:"tokenFactor"`
	CarbonOffsets []string              `json:"carbonOffsets"`
}

type getEnergyAttributionHelper struct {
	Id                   string                     `json:"_id"`
	Assetid              string                     `json:"asset_id"`
	Attributionholder    string                     `json:"attribution_holder"`
	Carbonoffset         string                     `json:"carbon_offset"`
	Energyproduced       string                     `json:"energy_produced"`
	Actualenergyproduced string                     `json:"actual_energy_produced"`
	Tokenaward           string                     `json:"token_award"`
	Version              string                     `json:"version"`
	Origin               getEnergyAttributionOrigin `json:"origin"`
	Assettype            string                     `json:"asset_type"`
	Productionperiod     string                     `json:"production_period"`
	Timestamp            string                     `json:"timestamp"`
	Epoch                string                     `json:"epoch"`
	Blockhash            string                     `json:"block_hash"`
	Validationauthority  string                     `json:"validation_authority"`
	Signature            string                     `json:"signature"`
	Tokenid              float64                    `json:"token_id"`
	CreatedAt            string                     `json:"createdAt"`
	UpdatedAt            string                     `json:"updatedAt"`
	Inputs               getEnergyAttributionInputs `json:"inputs"`
	Tags                 string                     `json:"tags"`
	Txhistory            string                     `json:"tx_history"`
	Processingstatus     string                     `json:"processing_status"`
	Transactions         []string                   `json:"transactions"`
	Redeemable           bool                       `json:"redeemable"`
	Claimed              bool                       `json:"claimed"`
	Confirmed            bool                       `json:"confirmed"`
}

// GetEnergyAttributionData is a helper struct for the swytch API
type GetEnergyAttributionData struct {
	Data []getEnergyAttributionHelper `json:"data"`
}

func getEnergyAttribution() {
	http.HandleFunc("/swytch/geteattributes", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["assetId"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		authToken := r.URL.Query()["authToken"][0]
		assetId := r.URL.Query()["assetId"][0]

		url := "https://platformapi-staging.swytch.io/v1/assets/" + assetId + "/attributions?limit=100&offset=0"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+authToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		defer func() {
			if ferr := res.Body.Close(); ferr != nil {
				err = ferr
			}
		}()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		log.Println(string(body))
		var x GetEnergyAttributionData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}
