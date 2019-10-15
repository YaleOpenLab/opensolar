package main

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"time"

	erpc "github.com/Varunram/essentials/rpc"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
)

// StartTeller starts the teller
func StartTeller() error {
	var err error

	client = erpc.SetupLocalHttpsClient(os.Getenv("HOME")+"/go/src/github.com/YaleOpenLab/opensolar/server.crt", 60*time.Second)

	err = login(loginUsername, loginPwhash)
	if err != nil {
		return errors.Wrap(err, "Error while logging on to the platform")
	}

	LocalProject, err = getLocalProjectDetails(loginProjIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't get local project details")
	}

	if LocalProject.Index == 0 {
		return errors.New("couldn't retrieve project from the database, quitting")
	}

	projIndex, err := getProjectIndex(AssetName)
	if err != nil {
		return errors.Wrap(err, "couldn't get project index")
	}

	if projIndex != LocalProject.Index {
		log.Println("Project indices don't match, quitting!")
		return errors.New("Project indices don't match, quitting!")
	}

	go refreshLogin(loginUsername, loginPwhash) // update local copy of the recipient every 5 minutes

	seed, err := wallet.DecryptSeed(LocalRecipient.U.StellarWallet.EncryptedSeed, LocalSeedPwd)
	if err != nil {
		return errors.Wrap(err, "Error while decrypting seed")
	}

	pubkey, err := wallet.ReturnPubkey(seed)
	if err != nil {
		return errors.Wrap(err, "Error while returning publickey")
	}

	if pubkey != LocalRecipient.U.StellarWallet.PublicKey {
		log.Println("PUBLIC KEYS DON'T MATCH, QUITTING!")
		return errors.New("PUBLIC KEYS DON'T MATCH, QUITTING!")
	}

	if LocalProject.Stage < 4 {
		log.Println("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
		return errors.New("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
	}

	// check for device id and set if none is set
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
