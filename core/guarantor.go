package core

import (
	"log"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	consts "github.com/YaleOpenLab/opensolar/consts"
)

// AddFirstLossGuarantee adds the given entity as a first loss guarantor
func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}

// RefillEscrowAsset refills the escrow with USD from the guarantor's account. Escrow
// should already have a trustline set with the stablecoin provider
func (a *Entity) RefillEscrowAsset(projIndex int, asset string,
	amount float64, seedpwd string) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	balancex := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, asset)
	balance, err := utils.ToFloat(balancex)
	if err != nil {
		log.Println(err)
		return err
	}

	if balance < amount {
		log.Println("guarantor does not required amount, refilling what amount they have")
		amount = balance - 1.0 // fees
	}

	seed, err := wallet.DecryptSeed(a.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	if !consts.Mainnet {
		_, txhash, err := assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey,
			project.EscrowPubkey, amount, seed, "guarantor refund")
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("txhash: ", txhash)
	} else {
		_, txhash, err := assets.SendAsset(consts.AnchorUSDCode, consts.AnchorUSDAddress,
			project.EscrowPubkey, amount, seed, "guarantor refund")
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("txhash: ", txhash)
	}
	return nil
}

// RefillEscrowXLM refills the escrow with XLM from the guarantor's account
func (a *Entity) RefillEscrowXLM(projIndex int, amount float64, seedpwd string) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	balancex := xlm.GetNativeBalance(a.U.StellarWallet.PublicKey)
	balance, err := utils.ToFloat(balancex)
	if err != nil {
		log.Println(err)
		return err
	}

	if balance < amount {
		log.Println("guarantor does not required amount, refilling what amount they have")
		amount = balance - 1.0 // fees
	}

	seed, err := wallet.DecryptSeed(a.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	_, txhash, err := xlm.SendXLM(project.EscrowPubkey, amount, seed, "guarantor refund")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("txhash: ", txhash)
	return nil
}
