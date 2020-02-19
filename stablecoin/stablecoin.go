package stablecoin

import (
	"encoding/json"
	"log"
	"net/url"

	tickers "github.com/Varunram/essentials/exchangetickers"
	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/opensolar/consts"
)

// TokenResponse is the token return endpoint provided by openx
type TokenResponse struct {
	Token string `json:"Token"`
}

// GetTestStablecoin gets test stablecoin from openx, amountx in USD
func GetTestStablecoin(username string, pubkey string, seed string, amountx float64) error {
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

	rate := tickers.ExchangeXLMforUSD(1) // 1 XLM = rate USD
	amountx = amountx / rate

	amount, err := utils.ToString(amountx)
	if err != nil {
		log.Println(err)
		return err
	}

	body := consts.OpenxURL + "/stablecoin/get?username=" + username + "&token=" + x.Token + "&seedpwd=x&amount=" + amount
	log.Println("STABLECOIN REQ: ", body)
	go erpc.GetRequest(body)
	return nil
}
