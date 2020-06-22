package main

import (
	"encoding/json"
	"log"
	"sync"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
	"github.com/Varunram/essentials/xlm"
	core "github.com/YaleOpenLab/opensolar/core"
)

func validateRecp(wg *sync.WaitGroup, username, token string) {
	defer wg.Done()
	body := "/recipient/validate?username=" + username + "&token=" + token

	data, err := erpc.GetRequest(platformURL + body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(data, &Recipient)
	if err != nil {
		log.Println(err)
		return
	}

	if Recipient.U != nil {
		if Recipient.U.Index != 0 {
			Return.Validate.Text = "Validated Recipient"
			Return.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + Token
		}
	} else {
		Return.Validate.Text = "Could not validate Recipient"
		Return.Validate.Link = platformURL + "/recipient/validate?username=" + username + "&token=" + Token
	}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		primNativeBalance := xlm.GetNativeBalance(Recipient.U.StellarWallet.PublicKey) * XlmUSD
		if primNativeBalance < 0 {
			primNativeBalance = 0
		}

		var err error
		Pnb, err = utils.ToString(primNativeBalance)
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		primUsdBalance := xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, "STABLEUSD")
		if primUsdBalance < 0 {
			primUsdBalance = 0
		}

		var err error
		Pub, err = utils.ToString(primUsdBalance)
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		secNativeBalance := xlm.GetNativeBalance(Recipient.U.SecondaryWallet.PublicKey) * XlmUSD
		if secNativeBalance < 0 {
			secNativeBalance = 0
		}

		var err error
		Snb, err = utils.ToString(secNativeBalance)
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		secUsdBalance := xlm.GetAssetBalance(Recipient.U.SecondaryWallet.PublicKey, "STABLEUSD")
		if secUsdBalance < 0 {
			secUsdBalance = 0
		}

		var err error
		Sub, err = utils.ToString(secUsdBalance)
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		Return.DABalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.DebtAssetCode))
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		Return.PBBalance.Text, err = utils.ToString(xlm.GetAssetBalance(Recipient.U.StellarWallet.PublicKey, Project.PaybackAssetCode))
		if err != nil {
			log.Println(err)
			return
		}
	}(&wg2)

	wg2.Wait()
}

func adminTokenHandler(wg *sync.WaitGroup, index int) {
	defer wg.Done()

	var wg2 sync.WaitGroup

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var projCount length
		data, err := erpc.GetRequest(platformURL + "/admin/getallprojects?username=admin&token=" + AdminToken)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(data, &projCount)
		if err != nil {
			log.Println(err)
			return
		}

		Return.ProjCount.Text, err = utils.ToString(projCount.Length)
		if err != nil {
			log.Println(err)
			return
		}

		Return.ProjCount.Link = platformURL + "/admin/getallprojects?username=admin&token=" + AdminToken
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var userCount length
		data, err := erpc.GetRequest(platformURL + "/admin/getallusers?username=admin&token=" + AdminToken)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(data, &userCount)
		if err != nil {
			log.Println(err)
			return
		}

		Return.UserCount.Text, err = utils.ToString(userCount.Length)
		if err != nil {
			log.Println(err)
			return
		}

		Return.UserCount.Link = platformURL + "/admin/getallusers?username=admin&token=" + AdminToken
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var userCount length
		data, err := erpc.GetRequest(platformURL + "/admin/getallinvestors?username=admin&token=" + AdminToken)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(data, &userCount)
		if err != nil {
			log.Println(err)
			return
		}

		Return.InvCount.Text, err = utils.ToString(userCount.Length)
		if err != nil {
			log.Println(err)
			return
		}

		Return.InvCount.Link = platformURL + "/admin/getallinvestors?username=admin&token=" + AdminToken
	}(&wg2)

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var userCount length
		data, err := erpc.GetRequest(platformURL + "/admin/getallrecipients?username=admin&token=" + AdminToken)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(data, &userCount)
		if err != nil {
			log.Println(err)
			return
		}

		Return.RecpCount.Text, err = utils.ToString(userCount.Length)
		if err != nil {
			log.Println(err)
			return
		}

		Return.RecpCount.Link = platformURL + "/admin/getallrecipients?username=admin&token=" + AdminToken
	}(&wg2)

	indexS, err := utils.ToString(index)
	if err != nil {
		log.Println(err)
		return
	}

	body := "/project/get?index=" + indexS

	data, err := erpc.GetRequest(platformURL + body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(data, &Project)
	if err != nil {
		log.Println(err)
		return
	}

	wg2.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		escrowBalance := xlm.GetAssetBalance(Project.EscrowPubkey, "STABLEUSD")
		if escrowBalance < 0 {
			escrowBalance = 0
		}

		escrowBalanceS, err := utils.ToString(escrowBalance)
		if err != nil {
			log.Println(err)
			return
		}

		Return.EscrowBalance.Text = escrowBalanceS
		Return.EscrowBalance.Link = "https://testnet.steexp.com/account/" + Project.EscrowPubkey
	}(&wg2)

	invIndex, err := utils.ToString(Project.InvestorIndices[0])
	if err != nil {
		log.Println(err)
		return
	}

	devIndex, err := utils.ToString(Project.DeveloperIndices[0])
	if err != nil {
		log.Println(err)
		return
	}

	wg2.Add(1)
	go getInvestor(&wg2, AdminToken, invIndex)
	wg2.Add(1)
	go getDeveloper(&wg2, AdminToken, devIndex)

	wg2.Wait()
}

func getInvestor(wg *sync.WaitGroup, AdminToken string, invIndex string) {
	defer wg.Done()
	data, err := erpc.GetRequest(platformURL + "/admin/getinvestor?username=admin&token=" +
		AdminToken + "&index=" + invIndex)
	if err != nil {
		log.Println(err)
		return
	}

	var investor core.Investor
	err = json.Unmarshal(data, &investor)
	if err != nil {
		log.Println(err)
		return
	}

	Return.Investor.Name = investor.U.Name
	Return.Investor.Username = investor.U.Username
	Return.Investor.Email = investor.U.Email
}

func getDeveloper(wg *sync.WaitGroup, AdminToken string, devIndex string) {
	defer wg.Done()
	data, err := erpc.GetRequest(platformURL + "/admin/getentity?username=admin&token=" +
		AdminToken + "&index=" + devIndex)
	if err != nil {
		log.Println(err)
		return
	}

	var developer core.Entity
	err = json.Unmarshal(data, &developer)
	if err != nil {
		log.Println(err)
		return
	}

	Return.Developer.Name = developer.U.Name
	Return.Developer.Username = developer.U.Username
	Return.Developer.Email = developer.U.Email
}
