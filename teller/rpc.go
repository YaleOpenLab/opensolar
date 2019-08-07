package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	core "github.com/YaleOpenLab/opensolar/core"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	geo "github.com/martinlindhe/google-geolocate"
)

// GetLocation gets the teller's location
func GetLocation(mapskey string) string {
	// see https://developers.google.com/maps/documentation/geolocation/intro on how
	// to improve location accuracy
	client := geo.NewGoogleGeo(mapskey)
	res, err := client.Geolocate()
	if err != nil {
		log.Println("Error while getting location: ", err)
		return ""
	}
	location := fmt.Sprintf("Lat%fLng%f", res.Lat, res.Lng) // some random format, can be improved upon if necessary
	DeviceLocation = location
	return location
}

// PingRpc pings the platform to see if its up
func PingRpc() error {
	// make a curl request out to lcoalhost and get the ping response
	data, err := erpc.GetRequest(ApiUrl + "/ping")
	if err != nil {
		return err
	}
	var x erpc.StatusResponse
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	// the result would be the status of the platform
	codeString, err := utils.ToString(x.Code)
	if err != nil {
		return err
	}
	ColorOutput("PLATFORM STATUS: "+codeString, GreenColor)
	return nil
}

// GetProjectIndex gets a specific project's index
func GetProjectIndex(assetName string) (int, error) {
	data, err := erpc.GetRequest(ApiUrl + "/projects?index=7")
	if err != nil {
		log.Println("Error while making get request: ", err)
		return -1, err
	}
	var x solar.SolarProjectArray
	err = json.Unmarshal(data, &x)
	if err != nil {
		return -1, err
	}
	for _, elem := range x {
		if elem.DebtAssetCode == assetName {
			return elem.Index, nil
		}
	}
	return -1, errors.New("Not found")
}

// LoginToPlatform logs on to the platform
func LoginToPlatform(username string, pwhash string) error {
	data, err := erpc.GetRequest(ApiUrl + "/recipient/validate?" + "username=" + username + "&pwhash=" + pwhash)
	if err != nil {
		return err
	}
	var x core.Recipient
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	ColorOutput("AUTHENTICATED RECIPIENT", GreenColor)
	LocalRecipient = x
	return nil
}

// ProjectPayback pays back to the platform
func ProjectPayback(assetName string, amountx float64) error {
	amount, err := utils.ToString(amountx)
	if err != nil {
		return err
	}
	// retrieve project index
	log.Println("PAYMENT BODY: ", ApiUrl+"/recipient/payback?"+"username="+LocalRecipient.U.Username+
		"&pwhash="+LocalRecipient.U.Pwhash+"&projIndex="+LocalProjIndex+"&assetName="+LocalProject.DebtAssetCode+"&seedpwd="+
		LocalSeedPwd+"&amount="+amount)
	data, err := erpc.GetRequest(ApiUrl + "/recipient/payback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&assetName=" + LocalProject.DebtAssetCode + "&seedpwd=" +
		LocalSeedPwd + "&amount=" + amount)
	if err != nil {
		return err
	}
	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	log.Println("PAYBACK RESPONSE: ", x)
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return errors.New("Errored out")
}

// SetDeviceId sets the device id of the teller
func SetDeviceId(username string, pwhash string, deviceId string) error {
	data, err := erpc.GetRequest(ApiUrl + "/recipient/deviceId?" + "username=" + username +
		"&pwhash=" + pwhash + "&deviceid=" + deviceId)
	if err != nil {
		return err
	}
	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("PAID!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreStartTime stores that start time of this particular instance
func StoreStartTime() error {
	unixString, err := utils.ToString(utils.Unix())
	if err != nil {
		return err
	}
	data, err := erpc.GetRequest(ApiUrl + "/recipient/startdevice?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&start=" + unixString)
	if err != nil {
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED START TIME SUCCESSFULLY!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreLocation stores the location of the teller
func StoreLocation(mapskey string) error {
	location := GetLocation(mapskey) // this happens to return null
	log.Println("MAPSKEY: ", mapskey, location)
	log.Println(ApiUrl + "/recipient/storelocation?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&location=" + location)
	data, err := erpc.GetRequest(ApiUrl + "/recipient/storelocation?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&location=" + location)
	if err != nil {
		log.Println("RPC ERROR IN STORELOCATION ENDPOINT")
		return err
	}

	log.Println("DATA: ", data)
	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("LOGGED LOCATION SUCCESSFULLY!", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// GetPlatformEmail gets the email of the platform
func GetPlatformEmail() error {
	data, err := erpc.GetRequest(ApiUrl + "/platformemail?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x rpc.PlatformEmailResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return err
	}
	PlatformEmail = x.Email
	ColorOutput("PLATFORMEMAIL: "+PlatformEmail, GreenColor)
	return nil
}

// SendDeviceShutdownEmail sends a shutdown notice to the platform
func SendDeviceShutdownEmail(tx1 string, tx2 string) error {

	data, err := erpc.GetRequest(ApiUrl + "/tellershutdown?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId +
		"&tx1=" + tx1 + "&tx2=" + tx2)
	if err != nil {
		log.Println(err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT STOP EMAIL SUCCESSFULLY", GreenColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// GetLocalProjectDetails gets the details of the local project
func GetLocalProjectDetails(projIndex string) (solar.Project, error) {

	var x solar.Project
	body := ApiUrl + "/project/get?index=" + projIndex
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return x, err
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return x, err
	}

	return x, nil
}

// SendDevicePaybackFailedEmail sends a notification if the payback routine breaks in its execution
func SendDevicePaybackFailedEmail() error {

	data, err := erpc.GetRequest(ApiUrl + "/tellerpayback?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&projIndex=" + LocalProjIndex + "&deviceId=" + DeviceId)
	if err != nil {
		log.Println(err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// StoreStateHistory stores state history in the data file
func StoreStateHistory(hash string) error {
	data, err := erpc.GetRequest(ApiUrl + "/recipient/ssh?" + "username=" + LocalRecipient.U.Username +
		"&pwhash=" + LocalRecipient.U.Pwhash + "&hash=" + hash)
	if err != nil {
		log.Println(err)
		return err
	}

	var x erpc.StatusResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x.Code == 200 {
		ColorOutput("SENT FAILED PAYBACK EMAIL", RedColor)
		return nil
	}
	return errors.New("Errored out, didn't receive 200")
}

// testSwytch tests whether the swytch workflow works correctly
func testSwytch() {
	body := ApiUrl + ApiUrl + "swytch/accessToken?" +
		"clientId=" + SwytchClientid + "&clientSecret=" + SwytchPassword + "&username=" + SwytchPassword +
		"&password=" + SwytchPassword
	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return
	}

	var x1 rpc.GetAccessTokenData
	err = json.Unmarshal(data, &x1)
	if err != nil {
		log.Fatal(err)
	}

	refresh_token := x1.Data[0].Refresh_token
	// we have the access token as well but need to refresh it using the refresh token, so
	// might as well store later.
	data, err = erpc.GetRequest("swytch/refreshToken?clientId=c0fe38566a254a3a80b2a42081b46843&clientSecret=46d10252a4954007af5e2f8941aeeb37&" +
		"refreshToken=" + refresh_token)
	if err != nil {
		log.Println(err)
		return
	}

	var x2 rpc.GetAccessTokenData
	err = json.Unmarshal(data, &x2)
	if err != nil {
		log.Fatal(err)
	}

	access_token := x1.Data[0].Access_token

	data, err = erpc.GetRequest(ApiUrl + "swytch/getuser?authToken=" + access_token)
	if err != nil {
		log.Println(err)
		return
	}

	var x3 rpc.GetSwytchUserStruct
	err = json.Unmarshal(data, &x3)
	if err != nil {
		log.Fatal(err)
	}

	user_id := x3.Data[0].Id
	log.Println("USER ID: ", user_id)
	// we have the user id, query for assets

	data, err = erpc.GetRequest(ApiUrl + "swytch/getassets?authToken=" + access_token + "&userId=" + user_id)
	if err != nil {
		log.Println(err)
		return
	}

	var x4 rpc.GetAssetStruct
	err = json.Unmarshal(data, &x4)
	if err != nil {
		log.Fatal(err)
	}

	asset_id := x4.Data[0].Id
	log.Println("ASSETID: ", asset_id)
	// we have the asset id, try to get some info
	data, err = erpc.GetRequest(ApiUrl + "swytch/getenergy?authToken=" + access_token + "&assetId=" + asset_id)
	if err != nil {
		log.Println(err)
		return
	}

	var x5 rpc.GetEnergyStruct
	err = json.Unmarshal(data, &x5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Energy data from installed asset: ", x4)

	data, err = erpc.GetRequest(ApiUrl + "swytch/getattributes?authToken=" + access_token + "&assetId=" + asset_id)
	if err != nil {
		log.Println(err)
		return
	}

	var x6 rpc.GetEnergyAttributionData
	err = json.Unmarshal(data, &x6)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Energy Attribute data: ", x6)
}
