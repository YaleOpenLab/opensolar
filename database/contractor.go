package opensolar

import (
	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
)

// When contractors are proposing a contract towards something,
// we need to ensure they are not following the price (eg bidding down) and are giving
// their best quote. In this scenario, a blind auction method is the best option.

// Propose proposes a new stage 2 contract
func (contractor *Entity) Propose(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, projectIndex int, auctionType string) (Project, error) {
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
	if err != nil {
		return pc, errors.Wrap(err, "couldn't retrieve recipient from db")
	}
	pc.RecipientIndex = iRecipient.U.Index
	pc.Stage = 2 // 2 since we need to filter this out while retrieving the propsoed contracts
	pc.AuctionType = auctionType
	pc.ContractorIndex = contractor.U.Index
	err = pc.Save()
	return pc, err
}

// AddCollateral adds a collateral that can be used as guarantee in case the contractor reneges
// on a particular contract
func (contractor *Entity) AddCollateral(amount float64, data string) error {
	contractor.Collateral += amount
	contractor.CollateralData = append(contractor.CollateralData, data)
	return contractor.Save()
}

// Slash slashes the contractor's reputation in the event of bad behaviour.
func (contractor *Entity) Slash(contractValue float64) error {
	// slash an entity's reputation score if it reneges on an agreed contract
	contractor.U.Reputation -= contractValue * 0.1
	return contractor.Save()
}

// RepInstalledProject adds reputatuon to the contractor on completion of installation of a project
// by default, we add reputation to the entity. In case the recipient wants to dispute this
// claim, we cna review and  change the reputation accordingly.
func RepInstalledProject(contrIndex int, projIndex int) error {
	contractor, err := RetrieveEntity(contrIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve all entities from db")
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project from db")
	}

	err = project.SetStage(5)
	if err != nil {
		return errors.Wrap(err, "couldn't set installed project's stage")
	}

	contractor.U.Reputation += project.TotalValue * ContractorWeight
	return contractor.Save()
}
