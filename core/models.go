package core

import (
	"log"
	"time"

	tickers "github.com/Varunram/essentials/exchangetickers"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"

	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	issuer "github.com/Varunram/essentials/xlm/issuer"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	stablecoin "github.com/YaleOpenLab/opensolar/stablecoin"

	consts "github.com/YaleOpenLab/opensolar/consts"
	notif "github.com/YaleOpenLab/opensolar/notif"
	oracle "github.com/YaleOpenLab/opensolar/oracle"
)

// MunibondInvest invests in a munibond. Sends USD to the platform, receives INVAssets
// in return, and sends an email to the investor's email id confirming investment if it succeeds.
func MunibondInvest(issuerPath string, invIndex int, invSeed string, invAmount float64,
	projIndex int, invAssetCode string, totalValue float64, seedInvestmentFactor float64, seed bool) error {

	var err error

	investor, err := RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve investor from database")
	}

	if !consts.Mainnet {
		usdBalance := xlm.GetAssetBalance(investor.U.StellarWallet.PublicKey, "STABLEUSD")
		if usdBalance < invAmount {
			// need to exchange stablecoin equivalent to the difference in balance plus some change
			amount := invAmount - usdBalance + 10
			err = stablecoin.GetTestStablecoin(investor.U.Username, investor.U.StellarWallet.PublicKey, invSeed, amount)
			if err != nil {
				return errors.Wrap(err, "Unable to offer xlm to STABLEUSD excahnge for investor")
			}
			time.Sleep(30 * time.Second)
		}
	}

	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}

	stableTxHash, err := SendUSDToPlatform(invSeed, invAmount, "Opensolar investment: "+projIndexString)
	if err != nil {
		return errors.Wrap(err, "Unable to send STABLEUSD to platform")
	}

	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.GetPath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve seed")
	}

	InvestorAsset := assets.CreateAsset(invAssetCode, issuerPubkey)

	invTrustTxHash, err := assets.TrustAsset(InvestorAsset.GetCode(), issuerPubkey, totalValue, invSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting investor asset")
	}

	log.Printf("Investor trusts InvAsset %s with txhash %s", InvestorAsset.GetCode(), invTrustTxHash)
	_, invAssetTxHash, err := assets.SendAssetFromIssuer(InvestorAsset.GetCode(), investor.U.StellarWallet.PublicKey, invAmount, issuerSeed, issuerPubkey)
	if err != nil {
		return errors.Wrap(err, "Error while sending out investor asset")
	}

	log.Printf("Sent InvAsset %s to investor %s with txhash %s", InvestorAsset.GetCode(), investor.U.StellarWallet.PublicKey, invAssetTxHash)

	investor.AmountInvested += invAmount

	if seed {
		investor.SeedInvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.GetCode())
		investor.SeedInvestedSolarProjectsIndices = append(investor.InvestedSolarProjectsIndices, projIndex)
	} else {
		investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.GetCode())
		investor.InvestedSolarProjectsIndices = append(investor.InvestedSolarProjectsIndices, projIndex)
	}

	err = investor.Save()
	if err != nil {
		return err
	}

	if investor.U.Notification {
		notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}
	return nil
}

// MunibondReceive sends Debt and Payback assets to the recipient. Sends a notification email
// to the recipient containing the tx hashes of all transactions involved.
func MunibondReceive(issuerPath string, recpIndex int, projIndex int, debtAssetID string,
	paybackAssetID string, years int, recpSeed string, totalValue float64, paybackPeriod time.Duration) error {

	log.Println("Retrieving recipient")
	recipient, err := RetrieveRecipient(recpIndex)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve recipient from database")
	}

	log.Println("Retrieving issuer")
	issuerPubkey, issuerSeed, err := wallet.RetrieveSeed(issuer.GetPath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	DebtAsset := assets.CreateAsset(debtAssetID, issuerPubkey)
	PaybackAsset := assets.CreateAsset(paybackAssetID, issuerPubkey)

	if years == 0 {
		years = 1
	}

	pbAmtTrust := float64(years * 12 * 2)

	paybackTrustHash, err := assets.TrustAsset(PaybackAsset.GetCode(), issuerPubkey, pbAmtTrust, recpSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting Payback Asset")
	}
	log.Printf("Recipient Trusts Payback asset %s with txhash %s", PaybackAsset.GetCode(), paybackTrustHash)

	_, paybackAssetHash, err := assets.SendAssetFromIssuer(PaybackAsset.GetCode(), recipient.U.StellarWallet.PublicKey, pbAmtTrust, issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		return errors.Wrap(err, "Error while sending payback asset from issue")
	}

	log.Printf("Sent PaybackAsset to recipient %s with txhash %s", recipient.U.StellarWallet.PublicKey, paybackAssetHash)

	debtTrustHash, err := assets.TrustAsset(DebtAsset.GetCode(), issuerPubkey, totalValue*2, recpSeed)
	if err != nil {
		return errors.Wrap(err, "Error while trusting debt asset")
	}
	log.Printf("Recipient Trusts Debt asset %s with txhash %s", DebtAsset.GetCode(), debtTrustHash)

	_, recpDebtAssetHash, err := assets.SendAssetFromIssuer(DebtAsset.GetCode(), recipient.U.StellarWallet.PublicKey, totalValue, issuerSeed, issuerPubkey) // same amount as debt
	if err != nil {
		return errors.Wrap(err, "Error while sending debt asset")
	}

	log.Printf("Sent DebtAsset to recipient %s with txhash %s\n", recipient.U.StellarWallet.PublicKey, recpDebtAssetHash)
	recipient.ReceivedSolarProjects = append(recipient.ReceivedSolarProjects, DebtAsset.GetCode())
	recipient.ReceivedSolarProjectIndices = append(recipient.ReceivedSolarProjectIndices, projIndex)
	err = recipient.Save()
	if err != nil {
		return errors.Wrap(err, "couldn't save recipient")
	}

	txhash, err := issuer.FreezeIssuer(issuerPath, projIndex, "blah")
	if err != nil {
		return errors.Wrap(err, "Error while freezing issuer")
	}

	log.Printf("Tx hash for freezing issuer is: %s", txhash)
	log.Printf("PROJECT %d's INVESTMENT CONFIRMED!", projIndex)

	if recipient.U.Notification {
		notif.SendInvestmentNotifToRecipient(projIndex, recipient.U.Email, paybackTrustHash, paybackAssetHash, debtTrustHash, recpDebtAssetHash)
	}

	go sendPaymentNotif(recipient.U.Index, projIndex, paybackPeriod, recipient.U.Email)
	return nil
}

// sendPaymentNotif sends a notification every payback period to the recipient to remind them
// to payback towards the project.
func sendPaymentNotif(recpIndex int, projIndex int, paybackPeriod time.Duration, email string) {
	paybackTimes := 0
	for {
		_, err := RetrieveRecipient(recpIndex)
		if err != nil {
			log.Println("Error while retrieving recipient from database", err)
			message := "Error while retrieving your account details, please contact help as soon as you receive this message " + err.Error()
			notif.SendAlertEmail(message, email) // don't catch the error here
			time.Sleep(time.Second * 2 * 604800)
		}

		if paybackTimes == 0 {
			// sleep and bother during the next cycle
			time.Sleep(time.Second * 2 * 604800)
		}

		// PAYBACK TIME!!
		// we don't know if the user has paid, but we send an email anyway
		notif.SendPaybackAlertEmail(projIndex, email)
		// sleep until the next payment is due
		paybackTimes++
		log.Println("Sent: ", email, "a notification on payments for payment cycle: ", paybackTimes)
		time.Sleep(paybackPeriod * consts.OneWeekInSecond)
	}
}

// MunibondPayback is used by the recipient to pay the platform back. Pays the
// project escrow USD, and the project issuer DebtAsset and Payback Asset.
func MunibondPayback(issuerPath string, recpIndex int, amount float64, recipientSeed string, projIndex int,
	assetName string, projectInvestors []int, totalValue float64, escrowPubkey string) (float64, error) {

	recipient, err := RetrieveRecipient(recpIndex)
	if err != nil {
		return -1, errors.Wrap(err, "Error while retrieving recipient from database")
	}

	issuerPubkey, _, err := wallet.RetrieveSeed(issuer.GetPath(issuerPath, projIndex), consts.IssuerSeedPwd)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to retrieve issuer seed")
	}

	monthlyBill := oracle.MonthlyBill() * float64(recipient.TellerEnergy) / 1000000
	if err != nil {
		return -1, errors.Wrap(err, "Unable to fetch oracle price, exiting")
	}

	log.Println("YOUR BILL: ", monthlyBill)

	if amount < monthlyBill {
		return -1, errors.New("amount paid is less than amount needed. Please refill your main account")
	}

	var StableBalance float64
	var xlmBalance float64

	if consts.Mainnet {
		StableBalance = xlm.GetAssetBalance(recipient.U.StellarWallet.PublicKey, consts.AnchorUSDCode)
	} else {
		StableBalance = xlm.GetAssetBalance(recipient.U.StellarWallet.PublicKey, consts.StablecoinCode)
	}

	xlmBalance = xlm.GetNativeBalance(recipient.U.StellarWallet.PublicKey)
	xlmUSD, err := tickers.BinanceTicker()
	if err != nil {
		return -1, errors.Wrap(err, "unable to fetch ticker price from binance")
	}

	balance := StableBalance + xlmUSD*xlmBalance

	if balance < amount {
		return -1, errors.Wrap(err, "You do not have the required stablecoin balance, please refill")
	}

	if StableBalance < amount {
		if consts.Mainnet {
			return -1, errors.New("need more stablecoin, exiting")
		}
		// need to exchange some XLM for stablecoin
		balNeeded := amount - StableBalance + 10 // some more for change, fees, etc
		err := stablecoin.GetTestStablecoin(recipient.U.Username, recipient.U.StellarWallet.PublicKey, recipientSeed, balNeeded)
		if err != nil {
			log.Println(err)
			return -1, errors.Wrap(err, "could not exchange xlm for stablecoin")
		}
		time.Sleep(30 * time.Second) // wait for the stablecoin daemon to give stablecoin
	}

	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return -1, err
	}

	var stablecoinHash string
	if !consts.Mainnet {
		_, stablecoinHash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey,
			escrowPubkey, amount, recipientSeed, "Opensolar payback: "+projIndexString)
		if err != nil {
			return -1, errors.Wrap(err, "Error while sending STABLEUSD back")
		}
	} else {
		_, stablecoinHash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress,
			escrowPubkey, amount, recipientSeed, "Opensolar payback: "+projIndexString)
		if err != nil {
			return -1, errors.Wrap(err, "Error while sending STABLEUSD back")
		}
	}

	log.Println("Paid", amount, " back to platform in stableUSD, txhash", stablecoinHash)

	_, debtPaybackHash, err := assets.SendAssetToIssuer(assetName, issuerPubkey, amount, recipientSeed)
	if err != nil {
		return -1, errors.Wrap(err, "Error while sending debt asset back")
	}
	log.Println("Paid", amount, " back to platform in DebtAsset, txhash", debtPaybackHash)

	if recipient.U.Notification {
		notif.SendPaybackNotifToRecipient(projIndex, recipient.U.Email, stablecoinHash, debtPaybackHash)
	}

	for _, i := range projectInvestors {
		investor, err := RetrieveInvestor(i)
		if err != nil {
			log.Println("Error while retrieving investor from list of investors", err)
			continue
		}
		if investor.U.Notification {
			notif.SendPaybackNotifToInvestor(projIndex, investor.U.Email, stablecoinHash, debtPaybackHash)
		}
	}

	ownershipAmt := amount - monthlyBill
	ownershipPct := ownershipAmt / totalValue

	recipient.NextPaymentInterval = utils.IntToHumanTime(utils.Unix() + 2419200)
	recipient.TellerEnergy = 0
	err = recipient.Save()
	if err != nil {
		log.Println(err)
	}

	return ownershipPct, nil
}

// SendUSDToPlatform sends STABLEUSD to the platform. Used by investors investing in projects.
func SendUSDToPlatform(invSeed string, invAmount float64, memo string) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them
	var oldPlatformBalance float64
	var err error
	var txhash string

	if !consts.Mainnet {
		oldPlatformBalance = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
		_, txhash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	} else {
		oldPlatformBalance = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.AnchorUSDCode)
		_, txhash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	}

	log.Println("Sent USD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	var newPlatformBalance float64
	if !consts.Mainnet {
		newPlatformBalance = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
	} else {
		newPlatformBalance = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.AnchorUSDCode)
	}

	if newPlatformBalance-oldPlatformBalance < invAmount-1 {
		return txhash, errors.New("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}
