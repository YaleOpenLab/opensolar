package opensolar

import (
	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
)

// An originator is someone who approaches the recipient in real life and proposes
// that he can start a contract on the opensolar platform that will be open ot investors
// he needs to make clear that he is an originator and if interested, he can volunteer
// to be the contractor as well, in which case there will be no auction and we can go
// straight ahead to the auction phase with investors investing in the contract. A MOU
// must also be sigend between the originator and the recipient defining terms of agreement
// as per legal standards

// Originate creates and saves a new origin contract
func (contractor *Entity) Originate(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, auctionType string) (Project, error) {

	var pc Project
	var err error

	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, errors.New("Projects could not be retrieved!")
	}
	pc.Index = len(indexCheck) + 1
	pc.PanelSize = panelSize
	pc.TotalValue = totalValue
	pc.State = location
	pc.EstimatedAcquisition = years
	pc.Metadata = metadata
	pc.DateInitiated = utils.Timestamp()
	iRecipient, err := RetrieveRecipient(recIndex)
	if err != nil { // recipient does not exist
		return pc, errors.Wrap(err, "couldn't retrieve recipient from db")
	}
	pc.RecipientIndex = iRecipient.U.Index
	pc.Stage = 0 // 0 since we need to filter this out while retrieving the propsoed contracts
	pc.AuctionType = auctionType
	pc.OriginatorIndex = contractor.U.Index
	pc.Reputation = totalValue // reputation is equal to the total value of the project
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
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
