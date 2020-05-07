package stablecoin

import (
	"encoding/json"
	"log"
	"net/url"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/opensolar/consts"
)

// TokenResponse is the token return endpoint provided by openx
type TokenResponse struct {
	Token string `json:"Token"`
}

// GetTestStablecoin gets test stablecoin from openx, amountx in USD
func GetTestStablecoin(username string, pubkey string, seedpwd string, amountx float64) error {
	form := url.Values{}
	form.Add("username", username)
	form.Add("pwhash", utils.SHA3hash("password"))

	retdata, err := erpc.PostForm(consts.OpenxURL+"/token", form)
	if err != nil {
		log.Println(err)
		return err
	}

	var x TokenResponse

	err = json.Unmarshal(retdata, &x)
	if err != nil {
		log.Println(err)
		return err
	}

	exchangeValue := 10000000.0 // 1XLM = 10**7 USD, hardcoded
	amountx = amountx / exchangeValue

	amount, err := utils.ToString(amountx)
	if err != nil {
		log.Println(err)
		return err
	}

	body := consts.OpenxURL + "/stablecoin/get?username=" +
		username + "&token=" + x.Token + "&seedpwd=" + seedpwd + "&amount=" + amount

	go erpc.GetRequest(body)
	return nil
}
