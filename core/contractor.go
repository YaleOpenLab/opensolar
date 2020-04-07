package core

import (
	"github.com/pkg/errors"
)

// Propose is called by a contractor when they want to propose a new stage 2 contract based on
// an originated project.
func (contractor *Entity) Propose(panelSize string, totalValue float64, location string,
	years int, metadata string, recIndex int, projectIndex int, auctionType string) (Project, error) {
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
	pc.Stage = 2
	pc.AuctionType = auctionType
	pc.ContractorIndex = contractor.U.Index
	err = pc.Save()
	return pc, err
}

// AddCollateral adds a collateral that can be used as guarantee in case the
// contractor reneges on a particular contract or changes fees later on during installation. The
// document should be stored in IPFS and the platform should guarantee its security.
func (contractor *Entity) AddCollateral(amount float64, data string) error {
	contractor.Collateral += amount
	contractor.CollateralData = append(contractor.CollateralData, data)
	return contractor.Save()
}

// Slash slashes the contractor's reputation in the event of bad behaviour. Slashing is a user
// trigerred action and happens only when entities report against a contractor.
func (contractor *Entity) Slash(contractValue float64) error {
	// slash an entity's reputation score if it reneges on an agreed contract
	contractor.U.Reputation -= contractValue * 0.1
	return contractor.Save()
}

// RepInstalledProject automatically adds reputation to the contractor on installation
// of a project. If entities want to reduce or report against the contractor later on,
// the slashing function above is called.
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
