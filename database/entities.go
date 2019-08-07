package opensolar

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"

	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	notif "github.com/YaleOpenLab/openx/notif"
)

// Entity defines a common structure for contractors, developers and originators. Will be split
// into their respective roles once they are defined in a better way.
type Entity struct {
	U *database.User
	// inherit the base user class
	Contractor bool
	// the name of the contractor / company that is contracting
	// A contractor is party who proposes a specific some of money towards a
	// particular project. This is the actual amount that the investors invest in.
	// This ideally must include the developer fee within it, so that investors
	// don't have to invest in two things. It would also make sense because the contractors
	// sometimes would hire developers themselves.
	Developer bool
	// A developer is someone who installs the required equipment (Raspberry Pi,
	// network adapters, anti tamper installations and similar) In the initial
	// projects, this will be us, since we'd be installing the pi ourselves, but in
	// the future, we expect third party developers / companies to do this for us
	// and act in a decentralized fashion. This money can either be paid out of chain
	// in fiat or can be a portion of the funds the investors chooses to invest in.
	// a contractor may also employ developers by himself, so this entity is not
	// strictly necessary.
	Originator bool
	// An Originator is an entity that will start a project and get a fixed fee for
	// rendering its service. An Originator's role is not restricted, the originator
	// can also be the developer, contractor or guarantor. The originator should take
	// the responsibility of auditing the requirements of the project - panel size,
	// location, number of panels needed, etc. He then should ideally be able to fill
	// out some kind of form on the website so that the originator's proposal is live
	// and shown to potential investors. The originators get paid only when the project
	// is live, else they can just spam, without any actual investment
	Guarantor bool
	// A guarantor is somebody who can assure investors that the school will get paid
	// on time. This authority should be trusted and either should be vetted by the law
	// or have a multisig paying out to the investors beyond a certain timeline if they
	// don't get paid by the school. This way, the guarantor can be anonymous, like the
	// nice Pineapple Fund guy. This can also be an insurance company, who is willing to
	// guarantee for specific school and the school can pay him out of chain / have
	// that as fee within the contract the originator
	PastContracts []Project
	// list of all the contracts that the contractor has won in the past
	ProposedContracts []Project
	// the Originator proposes a contract which will then be taken up
	// by a contractor, who publishes his own copy of the proposed contract
	// which will be the set of contracts that will be sent to auction
	PresentContracts []Project
	// list of all contracts that the contractor is presently undertaking
	PastFeedback []Feedback
	// feedback received on the contractor from parties involved in the past. This is an automated
	// feedback system which falls backto a manual one in the event of disputes
	Collateral float64
	// the amount of collateral that the entity is willing to hold in case it reneges
	// on a specific contract's details. This is an optional parameter but having collateral
	// would increase investor confidence that a particular entity will keep its word
	// regarding a particular contract.
	CollateralData []string
	// the specific thing(s) which the contractor wants to hold as collateral described
	// as a string (for eg, if a cash bond worth 5000 USD is held as collateral,
	// collateral would be set to 5000 USD and CollateralData would be "Cash Bond")
	FirstLossGuarantee string
	// FirstLossGuarantee is the seedpwd of the guarantor's account. This will be used
	// only when the recipient defaults and we need to cover losses of investors
	FirstLossGuaranteeAmt float64
	// FirstLossGuaranteeAmt is the amoutn that the guarantor is expected to cover in the case
	// of first loss
}

func (a *Entity) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, database.ContractorBucket, a, a.U.Index)
}

// RetrieveAllEntitiesWithoutRole gets all the entities in the opensolar platform
func RetrieveAllEntitiesWithoutRole() ([]Entity, error) {
	var users []Entity
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, database.ContractorBucket)
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

// RetrieveAllEntities gets all the proposed contracts for a particular recipient
func RetrieveAllEntities(role string) ([]Entity, error) {
	var entities []Entity

	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, database.ContractorBucket)
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

// RetrieveEntityHelper is a helper associated with the RetrieveEntity function
func RetrieveEntityHelper(key int) (Entity, error) {
	var entity Entity
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, database.ContractorBucket, key)
	if err != nil {
		return entity, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &entity)
	return entity, err
}

// RetrieveEntity retrieves a specific entity from the database
func RetrieveEntity(key int) (Entity, error) {
	var entity Entity
	user, err := database.RetrieveUser(key)
	if err != nil {
		return entity, err
	}

	entity, err = RetrieveEntityHelper(key)
	if err != nil {
		return entity, err
	}

	entity.U = &user
	return entity, entity.Save()
}

// newEntity creates a new entity based on the role passed
func newEntity(uname string, pwd string, seedpwd string, Name string, Address string, Description string, role string) (Entity, error) {
	var a Entity
	var err error
	user, err := database.NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "couldn't retrieve new user from db")
	}

	user.Address = Address
	user.Description = Description

	err = user.Save()
	if err != nil {
		return a, err
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
		return a, errors.New("invalid entity type passed!")
	}

	a.U = &user
	err = a.Save()
	return a, err
}

// TopReputationEntitiesWithoutRole returns the list of all the top reputed entities in descending order
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

// TopReputationEntities returns the list of all the top reputed entities with the specific role in descending order
func TopReputationEntities(role string) ([]Entity, error) {
	// caller knows what role he needs this list for, so directly retrieve and do stuff here
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

// ValidateEntity validates the entity with the specific name and pwhash and returns true if everything matches the thing on record
func ValidateEntity(name string, pwhash string) (Entity, error) {
	var rec Entity
	user, err := database.ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "couldn't validate user")
	}
	return RetrieveEntity(user.Index)
}

// AgreeToContractConditions agrees to some specific contract conditions
func AgreeToContractConditions(contractHash string, projIndex string,
	debtAssetCode string, entityIndex int, seedpwd string) error {
	// we need to display this on the frontend and once the user presses agree, commit
	// a tx to the blockchain with the outcome
	message := "I agree to the terms and conditions specified in contract " + contractHash +
		"and by signing this message to the blockchain agree that I accept the investment in project " + projIndex +
		"whose debt asset is: " + debtAssetCode
	// hash the message and transmit the message in 5 parts
	// eg.
	// CONTRACTHASH9a768ace36ff3d17
	// 71d5c145a544de3d68343b2e7609
	// 3cb7b2a8ea89ac7f1a20c852e6fc
	// 1d71275b43abffefac381c5b906f
	// 55c3bcff4225353d02f1d3498758

	user, err := database.RetrieveUser(entityIndex)
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
