package main

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	"github.com/spf13/viper"
)

// StartTeller starts the teller
func StartTeller() error {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Error while reading email values from config file")
	}

	if viper.Get("platformPublicKey") == nil || viper.Get("seedpwd") == nil || viper.Get("username") == nil ||
		viper.Get("password") == nil || viper.Get("apiurl") == nil || viper.Get("mapskey") == nil ||
		viper.Get("projIndex") == nil || viper.Get("assetName") == nil {
		return errors.New("Required parameters to be present in the config file: platformPublicKey, " +
			"seedpwd, username, password, apiurl, mapskey, projIndex, assetName (case-sensitive)")
	}

	PlatformPublicKey = viper.Get("platformPublicKey").(string)
	LocalSeedPwd = viper.Get("seedpwd").(string)                       // seed password used to unlock the seed of the recipient on the platform
	username := viper.Get("username").(string)                         // username of the recipient on the platform
	password := utils.SHA3hash(viper.Get("password").(string))         // password of the recipient on the platform
	ApiUrl = viper.Get("apiurl").(string)                              // ApiUrl of the remote / local openx node
	mapskey := viper.Get("mapskey").(string)                           // google maps API key
	LocalProjIndex, err = utils.ToString(viper.Get("projIndex").(int)) // get the project index which should be in the config file
	if err != nil {
		return err
	}
	AssetName = viper.Get("assetName").(string) // used to double check before starting the teller
	SwytchUsername = viper.Get("susername").(string)
	SwytchPassword = viper.Get("spassword").(string)
	SwytchClientid = viper.Get("sclientid").(string)
	SwytchClientSecret = viper.Get("sclientsecret").(string)

	projIndex, err := GetProjectIndex(AssetName)
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

	// don't allow login before this since that becomes an attack vector where a person can guess
	// multiple passwords
	err = LoginToPlatform(username, password)
	if err != nil {
		return errors.Wrap(err, "Error while logging on to the platform")
	}

	go RefreshLogin(username, password) // update local copy of the recipient every 5 minutes

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

	LocalProject, err = GetLocalProjectDetails(LocalProjIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't get local project details")
	}

	if LocalProject.Stage < 4 {
		log.Println("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
		return errors.New("TRYING TO INSTALL A PROJECT THAT HASN'T BEEN FUNDED YET, QUITTING!")
	}

	// check for device id and set it if none is set
	err = CheckDeviceID()
	if err != nil {
		return errors.Wrap(err, "could not check device id")
	}

	DeviceId, err = GetDeviceID() // Stores DeviceId
	if err != nil {
		return errors.Wrap(err, "could not get device id from local storage")
	}

	err = StoreStartTime()
	if err != nil {
		return errors.Wrap(err, "could not store start time locally")
	}

	// store location at the start because if a person changes location, it is likely that the
	// teller goes offline and we get notified
	err = StoreLocation(mapskey) // stores DeviceLocation
	if err != nil {
		return errors.Wrap(err, "could not store location of teller")
	}

	log.Println("STORED LOCATION SUCCESSFULLY")
	err = GetPlatformEmail()
	if err != nil {
		return errors.Wrap(err, "could not store platform email")
	}

	DeviceInfo = "Raspberry Pi3 Model B+"
	return nil
}
