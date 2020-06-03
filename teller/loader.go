package main

import (
	"log"

	"github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	wallet "github.com/Varunram/essentials/xlm/wallet"
)

// StartTeller starts the teller
func StartTeller() error {
	var err error

	erpc.SetConsts(60)
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
		return errors.New("project indices don't match, quitting")
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
		return errors.New("public keys don't match, quitting")
	}

	if LocalProject.Stage < 4 {
		log.Println("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
		return errors.New("trying to install a project that hasn't been funded yet, quitting")
	}

	// check for device id and set if none is set
	err = checkDeviceID()
	if err != nil {
		return errors.Wrap(err, "could not check device id")
	}

	DeviceID, err = getDeviceID() // Stores DeviceID
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
		colorOutput(RedColor, "could not store location of teller")
	}

	err = getPlatformEmail()
	if err != nil {
		return errors.Wrap(err, "could not store platform email")
	}

	DeviceInfo = "Raspberry Pi3 Model B+"
	return nil
}
