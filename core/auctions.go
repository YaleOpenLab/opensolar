package core

import (
	"github.com/pkg/errors"
)

// SelectContractBlind selects the winning bid based on blind auction rules. In a
// blind auction, the bid with the highest price wins
func SelectContractBlind(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, errors.New("empty array passed")
	}
	// array is not empty, min 1 elem
	a = arr[0]
	for _, elem := range arr {
		if elem.TotalValue < a.TotalValue {
			a = elem
			continue
		}
	}
	return a, nil
}

// SelectContractVickrey selects the winning bid based on vickrey auction rules. In a
// vickrey auction, the bid with the second highest price is paid by the highest bidder.
func SelectContractVickrey(arr []Project) (Project, error) {
	var winningContract Project
	if len(arr) == 0 {
		return winningContract, errors.New("empty array passed")
	}
	// array is not empty, min 1 elem
	winningContract = arr[0]
	var pos int
	for i, elem := range arr {
		if elem.TotalValue < winningContract.TotalValue {
			winningContract = elem
			pos = i
			continue
		}
	}
	// we have the highest bidder. Now we need to delete this guy from the array
	// and get the second highest bidder
	arr = append(arr[:pos], arr[pos+1:]...)
	if len(arr) == 0 {
		// means only one contract was proposed for this project, so fall back to blind auction
		return winningContract, nil
	}
	vickreyPrice := arr[0].TotalValue
	for _, elem := range arr {
		if elem.TotalValue < vickreyPrice {
			vickreyPrice = elem.TotalValue
		}
	}
	// we have the winner, who's elem and we have the price which is vickreyPrice
	// overwrite the winning contractor's contract
	winningContract.TotalValue = vickreyPrice
	return winningContract, winningContract.Save()
}

// SelectContractTime selects the winning contract based on the least time
// proposed for completion of installation.
func SelectContractTime(arr []Project) (Project, error) {
	var a Project
	if len(arr) == 0 {
		return a, errors.New("empty array passed")
	}

	a = arr[0]
	for _, elem := range arr {
		if elem.EstimatedAcquisition < a.EstimatedAcquisition {
			a = elem
			continue
		}
	}
	return a, nil
}

// SetAuctionType sets the auction type of a project. There are four options
// blind, civkrey, english and dutch.
func (project *Project) SetAuctionType(auctionType string) error {
	switch auctionType {
	case "blind":
		project.AuctionType = "blind"
	case "vickrey":
		project.AuctionType = "vickrey"
	case "english":
		project.AuctionType = "english"
	case "dutch":
		project.AuctionType = "dutch"
	default:
		project.AuctionType = "blind"
	}
	return project.Save()
}
