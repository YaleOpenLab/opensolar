package core

import (
	"github.com/pkg/errors"
)

// Originate creates and saves a new origin contract
func (contractor *Entity) Originate(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, auctionType string) (Project, error) {

	var pc Project
	var err error

	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, errors.New("projects could not be retrieved")
	}
	pc.Index = len(indexCheck) + 1
	pc.TotalValue = totalValue
	pc.State = location
	pc.EstimatedAcquisition = years
	pc.Metadata = metadata
	// pc.DateInitiated = utils.Timestamp()
	iRecipient, err := RetrieveRecipient(recIndex)
	if err != nil {
		return pc, errors.Wrap(err, "couldn't retrieve recipient from db")
	}
	pc.RecipientIndex = iRecipient.U.Index
	pc.Stage = 0
	pc.AuctionType = auctionType
	pc.OriginatorIndex = contractor.U.Index
	pc.Reputation = totalValue
	err = pc.Save()
	return pc, err
}

// RepOriginatedProject adds reputation to an originator on successful origination of a contract
func RepOriginatedProject(origIndex int, projIndex int) error {
	originator, err := RetrieveEntity(origIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve entity from db")
	}
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project from db")
	}
	return originator.U.ChangeReputation(project.TotalValue * OriginatorWeight)
}
