package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	geo "github.com/martinlindhe/google-geolocate"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	opensolar "github.com/YaleOpenLab/opensolar/core"
	rpc "github.com/YaleOpenLab/opensolar/rpc"
	orpc "github.com/YaleOpenLab/openx/rpc"
)

func baseURL(url string) string {
	return APIURL + "/" + url + "?username=" + LocalRecipient.U.Username + "&token=" + Token
}

func basePostData() url.Values {
	postdata := url.Values{}
	postdata.Set("username", LocalRecipient.U.Username)
	postdata.Set("token", Token)
	return postdata
}

func httpsGet(request []string, xparams ...string) ([]byte, error) {
	endpoint := request[0]
	reqParams := request[2:]

	if len(reqParams) != len(xparams) {
		colorOutput(CyanColor, "length of required params not equal to passed params: ", endpoint, reqParams, xparams)
		return nil, errors.New("length of required params not equal to passed params, quitting")
	}

	var params string
	params += baseURL(endpoint)
	for _, elem := range xparams {
		params += elem
	}

	return erpc.GetRequest(params)
}

// GetLocation gets the teller's location
func getLocation(mapskey string) error {
	// see https://developers.google.com/maps/documentation/geolocation/intro on how
	// to improve location accuracy
	client := geo.NewGoogleGeo(mapskey)
	res, err := client.Geolocate()
	if err != nil {
		colorOutput(RedColor, "Error while getting location: ", err)
		return err
	}
	location := fmt.Sprintf("Lat%fLng%f", res.Lat, res.Lng) // some random format, can be improved upon if necessary
	DeviceLocation = location
	return nil
}

// ping pings the platform to see if its up
func ping() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := erpc.GetRequest(APIURL + "/ping")
	if err != nil {
		return err
	}
	var x erpc.StatusResponse
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	// the result would be the status of the platform
	codeString, err := utils.ToString(x.Code)
	if err != nil {
		return err
	}
	colorOutput(GreenColor, "PLATFORM STATUS: "+codeString)
	return nil
}

// getProjectIndex gets a project's index
func getProjectIndex(assetName string) (int, error) {
	data, err := httpsGet(rpc.ProjectRPC[2])
	if err != nil {
		colorOutput(RedColor, "Error while making get request: ", err)
		return -1, err
	}

	var x []opensolar.Project
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return -1, err
	}
	for _, elem := range x {
		if elem.DebtAssetCode == assetName {
			return elem.Index, nil
		}
	}
	return -1, errors.New("Not found")
}

// LoginReturn is a wrapper around the token string
var LoginReturn struct {
	Token string
}

// login logs on to the platform
func login(username string, pwhash string) error {
	// we first need to login and then get the access token
	postdata := url.Values{}
	postdata.Set("username", username)
	postdata.Set("pwhash", pwhash)

	// Read in the cert file
	data, err := erpc.PostForm(APIURL+orpc.UserRPC[0][0], postdata)
	if err != nil {
		return errors.Wrap(err, "did not make request")
	}

	err = json.Unmarshal(data, &LoginReturn)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}

	log.Println(APIURL+orpc.UserRPC[0][0], postdata)
	// validate that the user is indeed a recipient
	Token = LoginReturn.Token
	data, err = erpc.GetRequest(APIURL + rpc.RecpRPC[3][0] + "?username=" + username + "&token=" + LoginReturn.Token)
	if err != nil {
		return err
	}

	var x core.Recipient
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	log.Println("FINE?", x)

	if x.U.Index == 0 {
		return errors.New("couldn't validate recipient")
	}

	colorOutput(GreenColor, "AUTHENTICATED RECIPIENT")
	LocalRecipient = x
	return nil
}

// ProjectPayback pays back to the platform
func projectPayback(assetName string, amountx float64) error {
	amount, err := utils.ToString(amountx)
	if err != nil {
		return err
	}

	projIndex, err := utils.ToString(LocalProject.Index)
	if err != nil {
		return err
	}

	var data []byte
	// retrieve project index
	if strings.Contains(APIURL, "localhost") {
		postdata := basePostData()
		postdata.Set("projIndex", projIndex)
		postdata.Set("assetName", assetName)
		postdata.Set("seedpwd", LocalSeedPwd)
		postdata.Set("amount", amount)

		data, err = erpc.PostForm(APIURL+rpc.RecpRPC[4][0], postdata)
	} else {
		form := url.Values{}
		form.Set("projIndex", projIndex)
		form.Set("assetName", assetName)
		form.Set("seedpwd", LocalSeedPwd)
		form.Set("amount", amount)

		data, err = erpc.PostForm(APIURL+rpc.RecpRPC[4][0], form)
	}

	if err != nil {
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	colorOutput(CyanColor, "PAYBACK RESPONSE: ", x)
	if x.Code == 200 {
		colorOutput(GreenColor, "PAID!")
		return nil
	}
	return errors.New("Errored out")
}

// SetDeviceId sets the device id of the teller
func setDeviceID(username string, deviceID string) error {

	postdata := url.Values{}
	log.Println(LocalRecipient.U.Username)
	log.Println(Token)
	postdata.Set("username", LocalRecipient.U.Username)
	postdata.Set("token", Token)
	log.Println("deviceid", deviceID)
	postdata.Set("deviceId", deviceID)

	data, err := erpc.PostForm(APIURL+rpc.RecpRPC[5][0], postdata)
	if err != nil {
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	if x.Code == 200 {
		colorOutput(GreenColor, "PAID!")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreStartTime stores that start time of this particular instance
func storeStartTime() error {
	unixString, err := utils.ToString(utils.Unix())
	if err != nil {
		return err
	}

	postdata := basePostData()
	postdata.Set("start", unixString)

	data, err := erpc.PostForm(APIURL+rpc.RecpRPC[6][0], postdata)
	if err != nil {
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	if x.Code == 200 {
		colorOutput(GreenColor, "LOGGED START TIME SUCCESSFULLY!")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreLocation stores the location of the teller
func storeLocation(mapskey string) error {
	err := getLocation(mapskey) // this happens to return null
	if err != nil {
		colorOutput(RedColor, err)
		return err
	}

	postdata := basePostData()
	postdata.Set("location", "l"+DeviceLocation) // handle google API failures this funky way

	data, err := erpc.PostForm(APIURL+rpc.RecpRPC[7][0], postdata)
	if err != nil {
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	if x.Code == 200 {
		colorOutput(GreenColor, "LOGGED LOCATION SUCCESSFULLY!")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// PlatformEmailResponse is a wrapper around the platform's email
type PlatformEmailResponse struct {
	Email string
}

// GetPlatformEmail gets the email of the platform
func getPlatformEmail() error {
	data, err := httpsGet(orpc.UserRPC[13])
	if err != nil {
		colorOutput(RedColor, err)
		return err
	}

	var x PlatformEmailResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}

	colorOutput(YellowColor, "PLATFORMEMAIL: "+x.Email)
	return nil
}

// SendDeviceShutdownEmail sends a shutdown notice to the platform
func sendDeviceShutdownEmail(tx1 string, tx2 string) error {

	projIndex, err := utils.ToString(LocalProject.Index)
	if err != nil {
		return err
	}

	data, err := httpsGet(rpc.ProjectRPC[6], "&projIndex="+projIndex,
		"&deviceId="+DeviceID, "&tx1="+tx1, "&tx2="+tx2)
	if err != nil {
		colorOutput(RedColor, err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	if x.Code == 200 {
		colorOutput(RedColor, "SENT STOP EMAIL SUCCESSFULLY")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// GetLocalProjectDetails gets the details of the local project
func getLocalProjectDetails(projIndexx int) (opensolar.Project, error) {
	var x opensolar.Project

	projIndex, err := utils.ToString(projIndexx)
	if err != nil {
		return x, err
	}

	data, err := httpsGet(rpc.ProjectRPC[3], "&index="+projIndex)
	if err != nil {
		colorOutput(RedColor, err)
		return x, err
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return x, err
	}

	return x, nil
}

// sendDevicePaybackFailedEmail sends a notification if the payback routine breaks in its execution
func sendDevicePaybackFailedEmail() error {

	projIndex, err := utils.ToString(LocalProject.Index)
	if err != nil {
		return err
	}

	data, err := httpsGet(rpc.ProjectRPC[7], "&projIndex="+projIndex, "&deviceId="+DeviceID)
	if err != nil {
		colorOutput(RedColor, err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}
	if x.Code == 200 {
		colorOutput(GreenColor, "SENT FAILED PAYBACK EMAIL")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// storeStateHistory stores state history in the data file
func storeStateHistory(hash string) error {
	postdata := basePostData()
	postdata.Set("hash", hash)

	data, err := erpc.PostForm(APIURL+rpc.RecpRPC[16][0], postdata)
	if err != nil {
		colorOutput(RedColor, err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return err
	}

	if x.Code == 200 {
		colorOutput(GreenColor, "SENT FAILED PAYBACK EMAIL")
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// testSwytch tests whether the swytch workflow works correctly
func testSwytch() {
	data, err := erpc.GetRequest(baseURL("swytch/accessToken") + "&clientId=" + SwytchClientid +
		"&clientSecret=" + SwytchClientSecret + "&username=" + SwytchUsername + "&password=" + SwytchPassword)
	if err != nil {
		colorOutput(RedColor, err)
		return
	}

	var x1 rpc.GetAccessTokenData
	err = json.Unmarshal(data, &x1)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	refreshToken := x1.Data[0].Refreshtoken
	// we have the access token as well but need to refresh it using the refresh token, so
	// might as well store later.
	data, err = erpc.GetRequest(APIURL + "/swytch/refreshToken?clientId=c0fe38566a254a3a80b2a42081b46843&clientSecret=46d10252a4954007af5e2f8941aeeb37&" +
		"refreshToken=" + refreshToken)
	if err != nil {
		colorOutput(RedColor, err)
		return
	}

	var x2 rpc.GetAccessTokenData
	err = json.Unmarshal(data, &x2)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	accessToken := x1.Data[0].Accesstoken

	data, err = erpc.GetRequest(APIURL + "/swytch/getuser?authToken=" + accessToken)
	if err != nil {
		colorOutput(RedColor, err)
		return
	}

	var x3 rpc.GetSwytchUserStruct
	err = json.Unmarshal(data, &x3)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	userID := x3.Data[0].ID
	colorOutput(CyanColor, "USER ID: ", userID)
	// we have the user id, query for assets

	data, err = erpc.GetRequest(APIURL + "/swytch/getassets?authToken=" + accessToken + "&userId=" + userID)
	if err != nil {
		colorOutput(RedColor, err)
		return
	}

	var x4 rpc.GetAssetStruct
	err = json.Unmarshal(data, &x4)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	assetID := x4.Data[0].ID
	colorOutput(CyanColor, "ASSETID: ", assetID)
	// we have the asset id, try to get some info
	data, err = erpc.GetRequest(APIURL + "/swytch/getenergy?authToken=" + accessToken + "&assetId=" + assetID)
	if err != nil {
		colorOutput(RedColor, err)
		return
	}

	var x5 rpc.GetEnergyStruct
	err = json.Unmarshal(data, &x5)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	colorOutput(CyanColor, "Energy data from installed asset: ", x4)

	data, err = erpc.GetRequest(APIURL + "/swytch/getattributes?authToken=" + accessToken + "&assetId=" + assetID)
	if err != nil {
		colorOutput(CyanColor, err)
		return
	}

	var x6 rpc.GetEnergyAttributionData
	err = json.Unmarshal(data, &x6)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return
	}

	colorOutput(CyanColor, "Energy Attribute data: ", x6)
}

func sendXLM(publickey string, amountx float64, memo string) (string, error) {
	amount, err := utils.ToString(amountx)
	if err != nil {
		colorOutput(RedColor, err)
		return "", err
	}

	data, err := httpsGet(orpc.UserRPC[7], "&destination="+
		publickey, "&amount="+amount, "&seedpwd="+LocalSeedPwd)

	if err != nil {
		colorOutput(RedColor, err)
		return "", err
	}

	var txhash string
	err = json.Unmarshal(data, &txhash)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return "", err
	}

	if len(txhash) != 66 { // include the quotes at the start and end
		data, err := httpsGet(orpc.UserRPC[7], "&destination="+
			publickey, "&amount="+amount, "&seedpwd="+LocalSeedPwd)

		if err != nil {
			colorOutput(RedColor, err)
			return "", err
		}

		var txhash string
		err = json.Unmarshal(data, &txhash)
		if err != nil {
			colorOutput(RedColor, string(data), err)
			return "", err
		}

		if len(txhash) != 66 { // include the quotes at the start and end
			return txhash, errors.New("xlm transaction not broadcast")
		}
	}
	return txhash, err
}

func getLatestBlockHash() (string, error) {
	log.Println("COOL?", orpc.UserRPC[33])
	data, err := httpsGet(orpc.UserRPC[33])
	if err != nil {
		colorOutput(RedColor, err)
		return "", err
	}

	var blockhash string
	err = json.Unmarshal(data, &blockhash)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return "", err
	}

	return blockhash, err
}

func getNativeBalance() (float64, error) {
	data, err := httpsGet(orpc.UserRPC[3])
	if err != nil {
		colorOutput(RedColor, err)
		return -1, err
	}

	var balance float64
	err = json.Unmarshal(data, &balance)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return -1, err
	}

	return balance, err
}

func getAssetBalance(asset string) (float64, error) {
	data, err := httpsGet(orpc.UserRPC[4], "&asset="+asset)
	if err != nil {
		colorOutput(RedColor, err)
		return -1, err
	}

	var balance float64
	err = json.Unmarshal(data, &balance)
	if err != nil {
		colorOutput(RedColor, string(data), err)
		return -1, err
	}

	return balance, err
}

func putIpfsData(data []byte) (string, error) {
	// retrieve project index
	postdata := basePostData()
	postdata.Set("data", string(data))

	data, err := erpc.PostForm(APIURL+orpc.UserRPC[34][0], postdata)
	if err != nil {
		return "", err
	}

	return string(data), err // return the hash
}

func getIpfsData(hash string) (string, error) {

	data, err := httpsGet(orpc.UserRPC[5], "&hash="+hash)
	if err != nil {
		return "", err
	}

	return string(data), err
}

func putEnergy(energyx uint32) ([]byte, error) {

	energy, err := utils.ToString(energyx)
	if err != nil {
		colorOutput(RedColor, err)
		return nil, err
	}

	postdata := basePostData()
	postdata.Set("energy", energy)

	data, err := erpc.PostForm(APIURL+rpc.RecpRPC[23][0], postdata)
	if err != nil {
		return nil, err
	}

	return data, err
}
