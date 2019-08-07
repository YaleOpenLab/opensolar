package opensolar

import (
	"log"
	"time"

	utils "github.com/Varunram/essentials/utils"
	stablecoin "github.com/YaleOpenLab/openx/chains/stablecoin"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	issuer "github.com/YaleOpenLab/openx/chains/xlm/issuer"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	notif "github.com/YaleOpenLab/openx/notif"
	oracle "github.com/YaleOpenLab/openx/oracle"
	"github.com/pkg/errors"
)

// MunibondInvest invests in a specific munibond
func MunibondInvest(issuerPath string, invIndex int, invSeed string, invAmount float64,
	projIndex int, invAssetCode string, totalValue float64, seedInvestmentFactor float64) error {
	// offer user to exchange xlm for stableusd and invest directly if the user does not have stableusd
	// this should be a menu on the Frontend but here we do this automatically
	var err error

	investor, err := RetrieveInvestor(invIndex)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve investor from database")
	}

	err = stablecoin.OfferExchange(investor.U.StellarWallet.PublicKey, invSeed, invAmount)
	if err != nil {
		return errors.Wrap(err, "Unable to offer xlm to STABLEUSD excahnge for investor")
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

	investor.AmountInvested += invAmount //  / seedInvestmentFactor -> figure out after demo
	investor.InvestedSolarProjects = append(investor.InvestedSolarProjects, InvestorAsset.GetCode())
	investor.InvestedSolarProjectsIndices = append(investor.InvestedSolarProjectsIndices, projIndex)
	// keep note of who all invested in this asset (even though it should be easy
	// to get that from the blockchain)
	err = investor.Save()
	if err != nil {
		return err
	}

	if investor.U.Notification {
		notif.SendInvestmentNotifToInvestor(projIndex, investor.U.Email, stableTxHash, invTrustTxHash, invAssetTxHash)
	}
	return nil
}

// MunibondReceive sends assets to the recipient
func MunibondReceive(issuerPath string, recpIndex int, projIndex int, debtAssetId string,
	paybackAssetId string, years int, recpSeed string, totalValue float64, paybackPeriod int) error {

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

	DebtAsset := assets.CreateAsset(debtAssetId, issuerPubkey)
	PaybackAsset := assets.CreateAsset(paybackAssetId, issuerPubkey)

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

// sendPaymentNotif sends a notification every payback period to the recipient to
// kindly remind him to payback towards the project
func sendPaymentNotif(recpIndex int, projIndex int, paybackPeriod int, email string) {
	// setup a payback monitoring routine for monitoring if the recipient pays us back on time
	// the recipient must give his email to receive updates
	paybackTimes := 0
	for {

		_, err := RetrieveRecipient(recpIndex) // need to retrieve to make sure nothing goes awry
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
		paybackTimes += 1
		log.Println("Sent: ", email, "a notification on payments for payment cycle: ", paybackTimes)
		time.Sleep(2 * time.Duration(paybackPeriod) * time.Second)
	}
}

// MunibondPayback is used by the recipient to pay the platform back. Here, we pay the
// project escrow instead of the platform since it would be responsible for redistribution of funds
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

	monthlyBill := oracle.MonthlyBill()
	if err != nil {
		return -1, errors.Wrap(err, "Unable to fetch oracle price, exiting")
	}

	log.Println("Retrieved average price from oracle: ", monthlyBill)

	if amount < monthlyBill {
		return -1, errors.New("amount paid is less than amount needed. Please refill your main account")
	}

	err = stablecoin.OfferExchange(recipient.U.StellarWallet.PublicKey, recipientSeed, amount)
	if err != nil {
		return -1, errors.Wrap(err, "Unable to offer xlm to STABLEUSD exchange for investor")
	}

	StableBalance, err := xlm.GetAssetBalance(recipient.U.StellarWallet.PublicKey, "STABLEUSD")

	if err != nil || (StableBalance < amount) {
		return -1, errors.Wrap(err, "You do not have the required stablecoin balance, please refill")
	}

	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return -1, err
	}

	var stablecoinHash string
	if !consts.Mainnet {
		_, stablecoinHash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, escrowPubkey, amount, recipientSeed, "Opensolar payback: "+projIndexString)
		if err != nil {
			return -1, errors.Wrap(err, "Error while sending STABLEUSD back")
		}
	} else {
		_, stablecoinHash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, escrowPubkey, amount, recipientSeed, "Opensolar payback: "+projIndexString)
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

	ownershipAmt := amount - monthlyBill
	ownershipPct := ownershipAmt / totalValue
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

	return ownershipPct, nil
}

// the models package won't be imported directly in any place but would be imported
// by all the investment models that exist

// SendUSDToPlatform sends STABLEUSD back to the platform for investment
func SendUSDToPlatform(invSeed string, invAmount float64, memo string) (string, error) {
	// send stableusd to the platform (not the issuer) since the issuer will be locked
	// and we can't use the funds. We also need ot be able to redeem the stablecoin for fiat
	// so we can't burn them
	var oldPlatformBalance float64
	var err error
	oldPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
	if err != nil {
		// platform does not have stablecoin, shouldn't arrive here ideally
		oldPlatformBalance = 0
	}

	var txhash string
	if !consts.Mainnet {
		_, txhash, err = assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	} else {
		_, txhash, err = assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress, consts.PlatformPublicKey, invAmount, invSeed, memo)
		if err != nil {
			return txhash, errors.Wrap(err, "sending stableusd to platform failed")
		}
	}

	log.Println("Sent STABLEUSD to platform, confirmation: ", txhash)
	time.Sleep(5 * time.Second) // wait for a block

	var newPlatformBalance float64
	if !consts.Mainnet {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.StablecoinCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	} else {
		newPlatformBalance, err = xlm.GetAssetBalance(consts.PlatformPublicKey, consts.AnchorUSDCode)
		if err != nil {
			return txhash, errors.Wrap(err, "error while getting asset balance")
		}
	}

	if newPlatformBalance-oldPlatformBalance < invAmount-1 {
		return txhash, errors.New("Sent amount doesn't match with investment amount")
	}
	return txhash, nil
}
