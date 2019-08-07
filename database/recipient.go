package opensolar

import (
	"github.com/pkg/errors"

	database "github.com/YaleOpenLab/openx/database"
)

// Recipient defines the recipient structure
type Recipient struct {
	U *database.User
	// user related functions are called as an instance directly
	ReceivedSolarProjects       []string
	ReceivedSolarProjectIndices []int
	ReceivedConstructionBonds   []string
	// ReceivedProjects denotes the projects that have been received by the recipient
	// instead of storing the PaybackAssets and the DebtAssets, we store this
	DeviceId string
	// the device ID of the associated solar hub. We don't do much with it here,
	// but we need it on the IoT Hub side to check login stuff
	DeviceStarts []string
	// the start time of the devices recorded for reference. We could monitor unscheduled
	// closes on the platform level as well and send email notifications or similar
	DeviceLocation string
	// the location of the device. Teller gets location using google's geolocation
	// API. Accuracy is of the order ~1km radius. Not great, but enough to detect
	// theft or something
	StateHashes []string
	// StateHashes provides the list of state updates (ipfs hashes) that the teller associated with this
	// particular recipient has communicated.
	TotalEnergyCP float64
	// the total energy produced by the recipient's assets in the current period
	TotalEnergy float64
	// the total energy produced by the recipient's assets over all billed periods
	Autoreload bool
	// a bool to denote whether the recipient wants to reload balance from his secondary account to pay any dues that are remaining
}

// NewRecipient returns a new recipient provided with the function parameters
func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	user, err := database.NewUser(uname, pwd, seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "failed to retrieve new user")
	}
	a.U = &user
	err = a.Save()
	return a, err
}
