package opensolar

import (
	//"log"
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
)

var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")

// NewOriginator creates a new originator
func NewOriginator(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "originator")
}

// NewDeveloper creates a new developer
func NewDeveloper(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "developer")
}

// NewGuarantor returns a new guarantor
func NewGuarantor(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "guarantor")
}

// NewContractor creates a new contractor and inherits properties from Users
func NewContractor(uname string, pwd string, seedpwd string, Name string, Address string, Description string) (Entity, error) {
	// Create a new entity with the boolean of 'contractor' set to 'true.' This is
	// done just by passing the string "contractor"
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "contractor")
}

// Save or Insert inserts a specific Project into the database
func (a *Project) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, database.ProjectsBucket, a, a.Index)
}

// Save inserts a passed Investor object into the database
func (a *Investor) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, InvestorBucket, a, a.U.Index)
}

// Save saves a given recipient's details
func (a *Recipient) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, RecipientBucket, a, a.U.Index)
}

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor(key int) (Investor, error) {
	var inv Investor
	user, err := database.RetrieveUser(key)
	if err != nil {
		return inv, err
	}

	x, err := edb.Retrieve(consts.DbDir+consts.DbName, InvestorBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &inv)
	if err != nil {
		return inv, errors.Wrap(err, "could not unmarshal investor")
	}

	inv.U = &user
	return inv, inv.Save()
}

// RetrieveRecipient retrieves a specific recipient from the database
func RetrieveRecipient(key int) (Recipient, error) {
	var recp Recipient
	user, err := database.RetrieveUser(key)
	if err != nil {
		return recp, err
	}

	x, err := edb.Retrieve(consts.DbDir+consts.DbName, RecipientBucket, key)
	if err != nil {
		return recp, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &recp)
	if err != nil {
		return recp, errors.New("could not unmarshal recipient")
	}

	recp.U = &user
	return recp, recp.Save()
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllInvestors() ([]Investor, error) {
	var arr []Investor

	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, InvestorBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys lim")
	}

	for _, value := range x {
		var temp Investor
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.Wrap(err, "error while unmarshalling json, quitting")
		}
		if temp.U.Index != 0 {
			arr = append(arr, temp)
		}
	}

	return arr, nil
}

// RetrieveAllRecipients gets a list of all Recipients in the database
func RetrieveAllRecipients() ([]Recipient, error) {
	var arr []Recipient

	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, RecipientBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Recipient
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.Wrap(err, "error while unmarshalling json, quitting")
		}
		if temp.U.Index != 0 {
			arr = append(arr, temp)
		}
	}

	return arr, nil
}

// TopReputationInvestors gets a list of all the investors with top reputation
func TopReputationInvestors() ([]Investor, error) {
	arr, err := RetrieveAllInvestors()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].U.Reputation > arr[j].U.Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// TopReputationRecipient returns a list of recipients with the best reputation
func TopReputationRecipients() ([]Recipient, error) {
	arr, err := RetrieveAllRecipients()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].U.Reputation > arr[j].U.Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// ValidateInvestor is a function to validate the investors username and password to log them into the platform, and find the details related to the investor
// This is separate from the publicKey/seed pair (which are stored encrypted in the database); since we can help users change their password, but we can't help them retrieve their seed.
func ValidateInvestor(name string, pwhash string) (Investor, error) {
	var rec Investor
	user, err := database.ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "failed to validate user")
	}
	return RetrieveInvestor(user.Index)
}

// ValidateRecipient validates a particular recipient
func ValidateRecipient(name string, pwhash string) (Recipient, error) {
	var rec Recipient
	user, err := database.ValidateUser(name, pwhash)
	if err != nil {
		return rec, errors.Wrap(err, "Error while validating user")
	}
	return RetrieveRecipient(user.Index)
}

// RetrieveProject retrieves the project with the specified index from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, database.ProjectsBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &inv)
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var projects []Project
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, database.ProjectsBucket)
	if err != nil {
		return projects, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range x {
		var temp Project
		err = json.Unmarshal(value, &temp)
		if err != nil {
			return projects, errors.New("could not unmarshal json")
		}
		projects = append(projects, temp)
	}

	return projects, nil
}

// RetrieveProjectsAtStage retrieves projects at a specific stage
func RetrieveProjectsAtStage(stage int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveContractorProjects retrieves projects that are associated with a specific contractor
func RetrieveContractorProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.ContractorIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveOriginatorProjects retrieves projects that are associated with a specific originator
func RetrieveOriginatorProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.OriginatorIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveRecipientProjects retrieves projects that are associated with a specific recipient
func RetrieveRecipientProjects(stage int, index int) ([]Project, error) {
	var arr []Project
	if stage > 9 { // check for this and fail early instead of wasting compute time on this
		return arr, errors.Wrap(errors.New("stage can not be greater than 9, quitting!"), "stage can not be greater than 9, quitting!")
	}

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Stage == stage && project.RecipientIndex == index {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// RetrieveLockedProjects retrieves all the projects that are locked and are waiting
// for the recipient to unlock them
func RetrieveLockedProjects() ([]Project, error) {
	var arr []Project

	projects, err := RetrieveAllProjects()
	if err != nil {
		return arr, err
	}

	for _, project := range projects {
		if project.Lock {
			arr = append(arr, project)
		}
	}

	return arr, nil
}

// SaveOriginatorMoU saves the MoU's hash in the platform's database
func SaveOriginatorMoU(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveContractHash saves a contract's hash in the platform's database
func SaveContractHash(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveInvPlatformContract saves the investor-platform contract's hash in the platform's database
func SaveInvPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveRecPlatformContract saves the recipient-platform contract's hash in the platform's database
func SaveRecPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}
