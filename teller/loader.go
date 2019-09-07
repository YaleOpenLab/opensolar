package main

import (
	"github.com/pkg/errors"
	"log"
	"os"

	utils "github.com/Varunram/essentials/utils"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	erpc "github.com/Varunram/essentials/rpc"
)

// StartTeller starts the teller
func StartTeller() error {
	var err error

	// don't allow login before this since that becomes an attack vector where a person can guess
	// multiple passwords
	client = erpc.SetupLocalHttpsClient(os.Getenv("HOME") + "/go/src/github.com/YaleOpenLab/opensolar/server.crt")
	err = login(Username, Pwhash)
	if err != nil {
		return errors.Wrap(err, "Error while logging on to the platform")
	}

	projIndex, err := getProjectIndex(AssetName)
	if err != nil {
		return errors.Wrap(err, "couldn't get project index")
	}

	projIndexS, err := utils.ToString(projIndex)
	if err != nil {
		return err
	}
	if projIndexS != LocalProjIndex {
		log.Println("Project indices don't match, quitting!")
		return errors.New("Project indices don't match, quitting!")
	}

	go refreshLogin(Username, Pwhash) // update local copy of the recipient every 5 minutes

	RecpSeed, err = wallet.DecryptSeed(LocalRecipient.U.StellarWallet.EncryptedSeed, LocalSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Error while decrypting seed")
	}

	RecpPublicKey, err = wallet.ReturnPubkey(RecpSeed)
	if err != nil {
		return errors.Wrap(err, "Error while returning publickey")
	}

	if RecpPublicKey != LocalRecipient.U.StellarWallet.PublicKey {
		log.Println("PUBLIC KEYS DON'T MATCH, QUITTING!")
		return errors.New("PUBLIC KEYS DON'T MATCH, QUITTING!")
	}

	LocalProject, err = getLocalProjectDetails(LocalProjIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't get local project details")
	}

	if LocalProject.Stage < 4 {
		log.Println("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
		return errors.New("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
	}

	// check for device id and set it if none is set
	err = checkDeviceID()
	if err != nil {
		return errors.Wrap(err, "could not check device id")
	}

	DeviceId, err = getDeviceID() // Stores DeviceId
	if err != nil {
		return errors.Wrap(err, "could not get device id from local storage")
	}

	err = storeStartTime()
	if err != nil {
		return errors.Wrap(err, "could not store start time locally")
	}

	// store location at the start because if a person changes location, it is likely that the
	// teller goes offline and we get notified
	err = storeLocation(Mapskey) // stores DeviceLocation
	if err != nil {
		return errors.Wrap(err, "could not store location of teller")
	}

	log.Println("STORED LOCATION SUCCESSFULLY")
	err = getPlatformEmail()
	if err != nil {
		return errors.Wrap(err, "could not store platform email")
	}

	DeviceInfo = "Raspberry Pi3 Model B+"
	return nil
}
