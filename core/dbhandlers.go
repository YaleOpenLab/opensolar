package core

import (
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

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

// NewContractor creates a new contractor
func NewContractor(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "contractor")
}

// Save saves a Project's details
func (a *Project) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, ProjectsBucket, a, a.Index)
}

// Save saves an Investor's details
func (a *Investor) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, InvestorBucket, a, a.U.Index)
}

// Save saves a Recipient's details
func (a *Recipient) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, RecipientBucket, a, a.U.Index)
}

// Save saves an Entity's details
func (a *Entity) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, ContractorBucket, a, a.U.Index)
}

// RetrieveInvestor retrieves an investor from the database
func RetrieveInvestor(key int) (Investor, error) {
	var inv Investor
	user, err := RetrieveUser(key)
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

// RetrieveRecipient retrieves a recipient from the database
func RetrieveRecipient(key int) (Recipient, error) {
	var recp Recipient
	user, err := RetrieveUser(key)
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

// RetrieveAllInvestors gets a list of all investors in the database
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

// TopReputationInvestors gets a list of all the investors sorted by descending reputation
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

// TopReputationRecipients returns a list of recipients sorted by descending reputation
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

// ValidateInvestor validates an investor's token and username
func ValidateInvestor(name string, token string) (Investor, error) {
	var rec Investor
	user, err := ValidateUser(name, token)
	if err != nil {
		return rec, errors.Wrap(err, "failed to validate user")
	}
	if user.Index == 0 {
		return rec, errors.New("Error while validating user")
	}
	return RetrieveInvestor(user.Index)
}

// ValidateRecipient validates a recipient's token and username
func ValidateRecipient(name string, token string) (Recipient, error) {
	var rec Recipient
	user, err := ValidateUser(name, token)
	if err != nil {
		return rec, errors.Wrap(err, "Error while validating user")
	}
	if user.Index == 0 {
		return rec, errors.New("Error while validating user")
	}
	return RetrieveRecipient(user.Index)
}

// RetrieveProject retrieves a project from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, ProjectsBucket, key)
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &inv)
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var projects []Project
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, ProjectsBucket)
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

// RetrieveProjectsAtStage retrieves projects at a specific stage from the database
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

// RetrieveContractorProjects retrieves projects that are associated with a specific contractor from the db
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

// RetrieveOriginatorProjects retrieves projects that are associated with a specific originator from the db
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

// RetrieveRecipientProjects retrieves projects that are associated with a specific recipient from the db
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

// RetrieveLockedProjects retrieves all the projects that are locked and are waiting for the recipient to unlock them
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

// SaveOriginatorMoU saves the MoU's hash in the database
func SaveOriginatorMoU(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveContractHash saves a contract's hash in the database
func SaveContractHash(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveInvPlatformContract saves the investor-platform contract's hash in the database
func SaveInvPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// SaveRecPlatformContract saves the recipient-platform contract's hash in the database
func SaveRecPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}
	a.StageData = append(a.StageData, hash)
	return a.Save()
}

// MarkFlagged is used by an admin to mark the project as flagged
func MarkFlagged(projIndex int, adminIndex int) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	if a.Reports > consts.ProjectReportThreshold {
		a.AdminFlagged = true
		a.FlaggedBy = adminIndex
	} else {
		return errors.New("project hasn't reached report threshold yet")
	}
	return a.Save()
}

// UserMarkFlagged is used by users to mark the project as flagged
func UserMarkFlagged(projIndex int, userIndex int) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	a.UserFlaggedBy = append(a.UserFlaggedBy, userIndex)
	a.Reports += 1
	return a.Save()
}

func AddTellerDetails(projIndex int, url string, brokerurl string, topic string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	a.TellerUrl = url
	a.BrokerUrl = brokerurl
	a.TellerPublishTopic = topic

	return a.Save()
}
