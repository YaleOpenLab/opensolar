package opensolar

import (
	// "log"
	"github.com/pkg/errors"

	tickers "github.com/YaleOpenLab/openx/chains/exchangetickers"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	database "github.com/YaleOpenLab/openx/database"
)

// Investor defines the investor structure
type Investor struct {
	VotingBalance float64 // this will be equal to the amount of stablecoins that the investor possesses,
	// should update this every once in a while to ensure voting consistency.
	// These are votes to show opinions about bids done by contractors on the specific projects that investors invested in.
	AmountInvested float64
	// total amount, would be nice to track to contact them,
	// give them some kind of medals or something
	InvestedSolarProjects        []string
	InvestedSolarProjectsIndices []int
	// array of asset codes this user has invested in
	U           *database.User
	WeightedROI string
	// the weightedROI for all the projects under the investor's umbrella
	AllTimeReturns []float64
	// the all time returns accumulated by the investor during his time on the platform indexed by project index
	ReceivedRECs string
	// The renewable enrgy  certificated received by the investor as part o
	Prorata string
	// the pro rata in all the projects that the in vestor has invested in
}

// NewInvestor creates a new investor object when passed the username, password hash,
// name and an option to generate the seed and publicKey.
func NewInvestor(uname string, pwd string, seedpwd string, Name string) (Investor, error) {
	var a Investor
	var err error
	user, err := database.NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "error while creating a new user")
	}
	a.U = &user
	a.AmountInvested = -1.0
	err = a.Save()
	return a, err
}

// AddVotingBalance adds / subtracts voting balance
func (a *Investor) ChangeVotingBalance(votes float64) error {
	// this function is caled when we want to refund the user with the votes once
	// an order has been finalized.
	a.VotingBalance += votes
	if a.VotingBalance < 0 {
		a.VotingBalance = 0 // to ensure no one has negative votes or something
	}
	return a.Save()
}

// CanInvest checks whether an investor has the required balance to invest in a project
func (a *Investor) CanInvest(targetBalance float64) bool {
	usdBalance, err := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, "STABLEUSD")
	if err != nil {
		usdBalance = 0
	}

	xlmBalance, err := xlm.GetNativeBalance(a.U.StellarWallet.PublicKey)
	if err != nil {
		xlmBalance = 0
	}

	// need to fetch the oracle price here for the order
	oraclePrice := tickers.ExchangeXLMforUSD(xlmBalance)
	if usdBalance > targetBalance || oraclePrice > targetBalance {
		// return true since the user has enough USD balance to pay for the order
		return true
	}
	return false
}
