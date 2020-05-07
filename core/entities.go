package core

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	wallet "github.com/Varunram/essentials/xlm/wallet"
	openx "github.com/YaleOpenLab/openx/database"

	consts "github.com/YaleOpenLab/opensolar/consts"
	notif "github.com/YaleOpenLab/opensolar/notif"
)

// Entity defines a common structure for contractors, developers and originators
type Entity struct {
	// U is the base User class inherited from openx
	U *openx.User

	// Contractor is a bool that is set if the entity is a contractor
	Contractor bool

	// Developer is a bool that is set if the entity is a developer
	Developer bool

	// Originator is a bool that is set if the entity is a originator
	Originator bool

	// Guarantor is a bool that is set if the entity is a guarantor
	Guarantor bool

	// PastContracts contains a list of all past contracts associated with the entity
	PastContracts []Project

	// ProposedContracts contains a list of all proposed contracts associated with the entity
	ProposedContract []Project

	// PresentContracts contains a list of all present contracts associated with the entity
	PresentContract []Project

	// ProposedContractIndices contains the indices of all proposed projects
	ProposedContractIndices []int

	// PresentContractIndices contains the indices of all present projects
	PresentContractIndices []int

	// PastFeedback contains a list of all feedback on the given entity
	PastFeedback []Feedback

	// Collateral is the amount the entity is willing to put up as collateral to secure projects
	Collateral float64

	// CollateralData contains data on the collateral amount that the entity is willing to pledge
	CollateralData []string

	// FirstLossGuarantee is the seed that will be used to transfer funds to investors in case the recipient refuses to pay
	FirstLossGuarantee string

	// 	FirstLossGuaranteeAmt is the amount that the guarantor is expected to cover in the case of a breach
	FirstLossGuaranteeAmt float64
}

// RetrieveAllEntitiesWithoutRole retrieves all the entities (contractors, developers,
// originators, and guarantors) from the database
func RetrieveAllEntitiesWithoutRole() ([]Entity, error) {
	var users []Entity
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, ContractorBucket)
	if err != nil {
		return users, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Entity
		err = json.Unmarshal(value, &temp)
		if err != nil {
			return users, errors.New("could not unmarshal json")
		}
		users = append(users, temp)
	}

	return users, nil
}

// RetrieveAllEntities gets all the proposed contracts associated with a particular entity
func RetrieveAllEntities(role string) ([]Entity, error) {
	var entities []Entity

	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, ContractorBucket)
	if err != nil {
		return entities, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var entity Entity
		err = json.Unmarshal(value, &entity)
		if err != nil {
			return entities, errors.New("could not unmarshal entity")
		}
		if entity.Contractor && role == "contractor" ||
			entity.Originator && role == "originator" ||
			entity.Guarantor && role == "guarantor" ||
			entity.Developer && role == "developer" {
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

// RetrieveEntity retrieves an entity from the database
func RetrieveEntity(key int) (Entity, error) {
	var entity Entity
	user, err := RetrieveUser(key)
	if err != nil {
		return entity, err
	}

	x, err := edb.Retrieve(consts.DbDir+consts.DbName, ContractorBucket, key)
	if err != nil {
		return entity, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &entity)
	if err != nil {
		return entity, err
	}

	entity.U = &user
	return entity, entity.Save()
}

// newEntity is a helper function that creates a new entity based on the passed roles
func newEntity(uname string, pwd string, seedpwd string, name string, role string) (Entity, error) {
	var a Entity
	var err error
	user, err := NewUser(uname, utils.SHA3hash(pwd), seedpwd, name)
	if err != nil {
		return a, errors.Wrap(err, "failed to retrieve new user")
	}

	switch role {
	case "contractor":
		a.Contractor = true
	case "developer":
		a.Developer = true
	case "originator":
		a.Originator = true
	case "guarantor":
		a.Guarantor = true
	default:
		return a, errors.New("invalid entity type passed")
	}

	a.U = &user
	err = a.Save()
	if err != nil {
		log.Println(err)
		return a, err
	}
	if !consts.Mainnet {
		// automatically get funds if on testnet
		err = xlm.GetXLM(user.StellarWallet.PublicKey)
		if err != nil {
			log.Println("couildn't get xlm: ", err)
		}
	}
	return a, err
}

// TopReputationEntitiesWithoutRole returns the list of all the entities in
// descending order of reputation
func TopReputationEntitiesWithoutRole() ([]Entity, error) {
	allEntities, err := RetrieveAllEntitiesWithoutRole()
	if err != nil {
		return allEntities, errors.Wrap(err, "couldn't retrieve all entities without role")
	}
	for i := range allEntities {
		for j := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

// TopReputationEntities returns the list of all the entities belonging to a
// role in descending order of reputation
func TopReputationEntities(role string) ([]Entity, error) {
	allEntities, err := RetrieveAllEntities(role)
	if err != nil {
		return allEntities, errors.Wrap(err, "couldn't retrieve all entities")
	}
	for i := range allEntities {
		for j := range allEntities {
			if allEntities[i].U.Reputation < allEntities[j].U.Reputation {
				tmp := allEntities[i]
				allEntities[i] = allEntities[j]
				allEntities[j] = tmp
			}
		}
	}
	return allEntities, nil
}

// ValidateEntity validates the username and token of the entity
func ValidateEntity(name string, token string) (Entity, error) {
	var rec Entity
	user, err := ValidateUser(name, token)
	if err != nil {
		return rec, errors.Wrap(err, "couldn't validate user")
	}
	return RetrieveEntity(user.Index)
}

// AgreeToContractConditions agrees to specified contract conditions. This is a precursor to
// a legeal contract template
func AgreeToContractConditions(contractHash string, projIndex string,
	debtAssetCode string, entityIndex int, seedpwd string) error {
	// we need to display this on the frontend and once the user presses agree, commit
	// a tx to the blockchain with the outcome

	message := "I agree to the terms and conditions specified in contract " + contractHash +
		"and by signing this message to the blockchain agree that I accept the investment in project " + projIndex +
		"whose debt asset is: " + debtAssetCode

	// hash the message and transmit the message in 5 parts due to stellar's memo field limit
	// eg.
	// CONTRACTHASH9a768ace36ff3d17
	// 71d5c145a544de3d68343b2e7609
	// 3cb7b2a8ea89ac7f1a20c852e6fc
	// 1d71275b43abffefac381c5b906f
	// 55c3bcff4225353d02f1d3498758

	user, err := RetrieveUser(entityIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve user from db")
	}

	seed, err := wallet.DecryptSeed(user.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		return errors.Wrap(err, "couldn't decrypt seed")
	}

	messageHash := "CONTRACTHASH" + strings.ToUpper(utils.SHA3hash(message))
	firstPart := messageHash[:28] // higher limit is not included in the slice
	secondPart := messageHash[28:56]
	thirdPart := messageHash[56:84]
	fourthPart := messageHash[84:112]
	fifthPart := messageHash[112:140]

	timestamp := float64(utils.Unix())

	_, firstHash, err := xlm.SendXLM(user.StellarWallet.PublicKey, timestamp, seed, firstPart)
	if err != nil {
		return errors.Wrap(err, "couldn't send tx 1")
	}

	_, secondHash, err := xlm.SendXLM(user.StellarWallet.PublicKey, timestamp, seed, secondPart)
	if err != nil {
		return errors.Wrap(err, "couldn't send tx 2")
	}

	_, thirdHash, err := xlm.SendXLM(user.StellarWallet.PublicKey, timestamp, seed, thirdPart)
	if err != nil {
		return errors.Wrap(err, "couldn't send tx 3")
	}

	_, fourthHash, err := xlm.SendXLM(user.StellarWallet.PublicKey, timestamp, seed, fourthPart)
	if err != nil {
		return errors.Wrap(err, "couldn't send tx 4")
	}

	_, fifthHash, err := xlm.SendXLM(user.StellarWallet.PublicKey, timestamp, seed, fifthPart)
	if err != nil {
		return errors.Wrap(err, "couldn't send tx 5")
	}

	if user.Notification {
		notif.SendContractNotification(firstHash, secondHash, thirdHash, fourthHash, fifthHash, user.Email)
	}

	return nil
}
