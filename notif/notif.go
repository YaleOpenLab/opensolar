package notif

import (
	email "github.com/Varunram/essentials/email"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// package notif is used to send out notifications regarding important events that take
// place with respect to a specific project / investment

// footerString is a common footer string that is used by all emails
var footerString = "Have a nice day!\n\nWarm Regards, \nThe OpenSolar Team\n\n\n\n" +
	"You're receiving this email because your contact was given" +
	" on the opensolar platform for receiving notifications on orders in which you're a party.\n\n\n"

// SendInvestmentNotifToRecipient sends a notification to the recipient when an investor
// invests in an order he's the recipient of
func SendInvestmentNotifToRecipient(projIndex int, to string, recpPbTrustHash string, recpAssetHash string, recpDebtTrustHash string, recpDebtAssetHash string) error {
	// this is sent to the recipient on investment from an investor
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that project number: " + projIndexString + " has been invested in.\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your payback trusted asset hash is: https://testnet.steexp.com/tx/" + recpPbTrustHash + "\n" +
		"Your payback asset hash is: https://testnet.steexp.com/tx/" + recpAssetHash + "\n" +
		"Your debt trusted asset hash is: https://testnet.steexp.com/tx/" + recpDebtTrustHash + "\n" +
		"Your debt asset hash is: https://testnet.steexp.com/tx/" + recpDebtAssetHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendInvestmentNotifToRecipientOZ sends a notification to the recipient as part of the opzones platform
func SendInvestmentNotifToRecipientOZ(projIndex int, to string, recpDebtTrustHash string, recpDebtAssetHash string) error {
	// this is sent to the recipient on investment from an investor
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that project number: " + projIndexString + " has been invested in.\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your debt trusted asset hash is: https://testnet.steexp.com/tx/" + recpDebtTrustHash + "\n" +
		"Your debt asset hash is: https://testnet.steexp.com/tx/" + recpDebtAssetHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendInvestmentNotifToInvestor sends a notification to the investor when he invests
// in a particular project
func SendInvestmentNotifToInvestor(projIndex int, to string, stableHash string, trustHash string, assetHash string) error {
	// this is sent to the investor on investment
	// this should ideally contain all the information he needs for a concise proof of
	// investment
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have invested in project number: " + projIndexString + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableHash + "\n" +
		"Your trusted asset hash is: https://testnet.steexp.com/tx/" + trustHash + "\n" +
		"Your investment asset hash is: https://testnet.steexp.com/tx/" + assetHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendSeedInvestmentNotifToInvestor sends a notification to the user after seed investment
func SendSeedInvestmentNotifToInvestor(projIndex int, to string, stableHash string, trustHash string, assetHash string) error {
	// this is sent to the investor on investment
	// this should ideally contain all the information he needs for a concise proof of
	// investment
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have invested in the seed round of project: " + projIndexString + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableHash + "\n" +
		"Your trusted asset hash is: https://testnet.steexp.com/tx/" + trustHash + "\n" +
		"Your investment asset hash is: https://testnet.steexp.com/tx/" + assetHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendPaybackNotifToRecipient sends a notification email to the recipient when he
// pays back towards a particular order
func SendPaybackNotifToRecipient(projIndex int, to string, stableUSDHash string, debtPaybackHash string) error {
	// this is sent to the recipient
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have paid back towards project number: " + projIndexString + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableUSDHash + "\n" +
		"Debt asset hash is: https://testnet.steexp.com/tx/" + debtPaybackHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendPaybackNotifToInvestor sends a notification email to the investor when the recipient
// pays back towards a particular order
func SendPaybackNotifToInvestor(projIndex int, to string, stableUSDHash string, debtPaybackHash string) error {
	// this is sent to the investor on payback from an investor
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that the recipient has paid back towards project number: " + projIndexString + "\n\n" +
		"The recipient's proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableUSDHash + "\n" +
		"Debt asset hash is: https://testnet.steexp.com/tx/" + debtPaybackHash + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendUnlockNotifToRecipient sends a notification email to the investor when the recipient
// pays back towards a particular order
func SendUnlockNotifToRecipient(projIndex int, to string) error {
	// this is sent to the investor on payback from an investor
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that project number: " + projIndexString + " has been invested in\n\n" +
		"You are required to logon to the platform within a period of 3(THREE) days in order to accept the investment\n\n" +
		"If you choose to not accept the given investment in your project, please be warned that your reputation score " +
		"will be adjusted accordingly and this may affect any future proposal that you seek funding for on the platform\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendUnlockNotifToRecipientOZ sends an unlock notification as part of the opzones platform
func SendUnlockNotifToRecipientOZ(projIndex int, to string) error {
	// this is sent to the investor on payback from an investor
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opzones platform! \n\n" +
		"We're writing to let you know that project number: " + projIndexString + " has been invested in\n\n" +
		"You are required to logon to the platform within a period of 3(THREE) days in order to accept the investment\n\n" +
		"If you choose to not accept the given investment in your project, please be warned that your reputation score " +
		"will be adjusted accordingly and this may affect any future proposal that you seek funding for on the platform\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendEmail is a hlper for the rpc to send an email to an entity
func SendEmail(message string, to string, name string) error {
	// we can't send emails directly since we would need their gmail usernames and password for that
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that " + name + " has sent you a message. The message contents follow: \n\n"
	body := startString + message + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendAlertEmail sends an alert email
func SendAlertEmail(message string, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that you have received a message from the platform: \n\n\n" + message
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendPaybackAlertEmail sends a payback alert email. We don't know if the user has paid and send
// this even if the user has paid / received a donation towards this month
func SendPaybackAlertEmail(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"This is a kind reminder to let you know that your payment is due this period for project numbered: " + projIndexString +
		"\n\n If you have already paid or have received a donation towards this month, please ignore this alert."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendNicePaybackAlertEmail sends an email when the amount for 2 payment cycles is due
func SendNicePaybackAlertEmail(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"This is a kind reminder to let you know that your payment is due this period for project numbered: " + projIndexString +
		"\n\n Please payback at the earliest."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendSternPaybackAlertEmail sends an email when the amount for 4 payment cycles is due.
func SendSternPaybackAlertEmail(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that your payment is due this period for project numbered: " + projIndexString +
		"\n\n Please payback within two payback cycles to avoid re-routing of power services."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendDisconnectionEmail sends an email when the amount for 6 payment cycles is due
func SendDisconnectionEmail(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + projIndexString +
		"\n\nHas been redirected towards the main power grid. Please contact your guarantor to resume services"
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendDisconnectionEmailI sends an email to the investor when the amount for 6 payment cycles is due on the recipient's end
func SendDisconnectionEmailI(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + projIndexString +
		"\n\nHas been redirected towards the main power grid due to irregular payments by the recipient involved.\n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendSternPaybackAlertEmailI sends a stern payback email notification to the investor
func SendSternPaybackAlertEmailI(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that we are aware that payments towards the project: " + projIndexString +
		"\n\n haven't been made and we have reached out to the project recipient on the same. If this situation continues for " +
		"two more payment periods, we will be redirecting power towards the general grid and you would receive payments " +
		"for all periods where they were due. \n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendSternPaybackAlertEmailG sends a stern payback email notification to the guarantor
func SendSternPaybackAlertEmailG(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that we are aware that payments towards the project: " + projIndexString +
		"\n\n haven't been made and have reached out to the project recipient on the same. If this situation continues for " +
		"two more payment periods, we will be redirecting power towards the general grid and contact you for further" +
		"information on how the guarantee towards the project would be realized to investors.\n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendDisconnectionEmailG sends a disconnection email notification to the guarantor
func SendDisconnectionEmailG(projIndex int, to string) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + projIndexString +
		"\n\nHas been redirected towards the main power grid due to irregular payments by the recipient involved.\n\n" +
		"We will be reaching out to you in the coming days on how to proceed with realizing the guarantee towards this " +
		"project in order to safeguard investors. We will also be contacting the recipient involved to update them on the" +
		"situation and will make efforts to alleviate this problem as soon as possible." +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return email.SendMail(body, to)
}

// SendContractNotification sends a notification after an entity signs a contract
func SendContractNotification(Hash1 string, Hash2 string, Hash3 string, Hash4 string, Hash5 string, to string) error {
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that you have signed a contract\n\n" +
		"Your proofs of signing are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your first hash is: https://testnet.steexp.com/tx/" + Hash1 + "\n" +
		"Your second hash is: https://testnet.steexp.com/tx/" + Hash2 + "\n" +
		"Your third hash is: https://testnet.steexp.com/tx/" + Hash3 + "\n" +
		"Your fourth hash is: https://testnet.steexp.com/tx/" + Hash4 + "\n" +
		"Your fifth hash is: https://testnet.steexp.com/tx/" + Hash5 + "\n\n\n" +
		footerString
	return email.SendMail(body, to)
}

// SendTellerShutdownEmail sends the platform an email notifying that the teller has shut down
func SendTellerShutdownEmail(from string, projIndex string, deviceId string, tx1 string, tx2 string) error {
	body := "Greetings from the remote teller " + deviceId + " installed for: " + from + " on behalf of project: " + projIndex + "\n\n" +
		"We're writing to let you know that the teller has shut down and requires your immediate action. The proof of shutdown transactions " +
		"are atached below:" + "\n\n" +
		"Tx1: https://testnet.steexp.com/tx/" + tx1 + "\n\n" +
		"Tx2: https://testnet.steexp.com/tx/" + tx2 + "\n\n" +
		"Please tend to this situation at the earliest." + "\n\n\n" +
		footerString
	return email.SendMail(body, consts.PlatformEmail)
}

// SendTellerPaymentFailedEmail is a notification ot the platform that the teller's payback routine has been disturbed
func SendTellerPaymentFailedEmail(from string, projIndex string, deviceId string) error {
	body := "Greetings from the remote teller " + deviceId + " installed for: " + from + " on behalf of project: " + projIndex + "\n\n" +
		"We're writing to let you know that the teller encountered an error, didn't result in automatic payback and requires your immediate action. " +
		"Please tend to this situation at the earliest." + "\n\n\n" +
		footerString
	return email.SendMail(body, consts.PlatformEmail)
}

// SendTellerDownEmail is an email to the platform notifying that the teller for a particular project is down.
func SendTellerDownEmail(projIndex int, recpIndex int) error {
	projIndexString, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}

	recpIndexString, err := utils.ToString(recpIndex)
	if err != nil {
		return err
	}
	body := "Greetings from the opensolar platform! \n\nWe're writing to let you know that remote teller " + projIndexString +
		" installed on behalf of recipient with index: " + recpIndexString + " has not been responding to pings for a while. Please take action at " +
		"the earliest," + "\n\n\n" +
		footerString
	return email.SendMail(body, consts.PlatformEmail)
}

// SendSecretsEmail is an email to trusted social contacts notifying that a user has shared a secret with them
func SendSecretsEmail(userEmail string, email1 string, email2 string, email3 string, secret1 string, secret2 string, secret3 string) error {
	bodyBase := "Greetings from the opensolar platform! \n\nWe're writing to let you know that user with email: " + userEmail +
		" has designated you as a trusted entity. Towards this, we request that you keep the attached secret in a safe and secure place and provide " +
		"it to the above user in case they request for it. \n\n" + "SECRET:\n\n"
	body1 := bodyBase + secret1 + "\n\n\n" + footerString
	err := email.SendMail(body1, email1)
	if err != nil {
		return err
	}

	body2 := bodyBase + secret2 + "\n\n\n" + footerString
	err = email.SendMail(body2, email2)
	if err != nil {
		return err
	}

	body3 := bodyBase + secret3 + "\n\n\n" + footerString
	err = email.SendMail(body3, email3)
	if err != nil {
		return err
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email to the email address of the user
func SendPasswordResetEmail(to string, vCode string) error {
	body := "Greetings from the opensolar platform! \n\nWe're writing to let you know that you requested a password reset recently\n\n" +
		"Please use this given code along with the link attached in order to reset your password\n\n" +
		"VERIFICATION CODE: " + vCode + "\n\n\n" + footerString

	return email.SendMail(body, to)
}
