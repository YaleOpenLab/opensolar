package core

import (
	"log"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	consts "github.com/YaleOpenLab/opensolar/consts"
	openx "github.com/YaleOpenLab/openx/database"
)

// Recipient defines the recipient structure
type Recipient struct {

	// U imports the base User class from openx
	U *openx.User

	// C is a structure containing all details of the company the investor is part of
	C Company

	// Company denotes whether the given investor is acting on behalf of a company
	Company bool

	// ReceivedSolarProjects stores the projects that the recipient is receiver of
	ReceivedSolarProjects []string

	// ReceivedSolarProjectIndices stores the indices of the projects the recipient is part of
	ReceivedSolarProjectIndices []int

	// DeviceID is the device ID of the associated solar hub / IoT device
	DeviceID string

	// DeviceStarts contains the start time of the above IoT devices.
	DeviceStarts []string

	// DeviceLocation stores the physical location of the device powered by Google APIs.
	DeviceLocation string

	// StateHashes stores the list of state updates (ipfs hashes) of the teller
	StateHashes []string

	// TellerEnergy contains the net energy consumed during a given period
	TellerEnergy uint32

	// PastTellerEnergy contains a list of the energy values accumulated by the project
	PastTellerEnergy []uint32

	// NextPaymentInterval stores the date of the next payment interval
	NextPaymentInterval string

	// Autoreload is a bool to denote whether the recipient wants to reload balance from their secondary account
	Autoreload bool
}

// NewRecipient creates and returns a new recipient
func NewRecipient(uname string, pwd string, seedpwd string, Name string) (Recipient, error) {
	var a Recipient
	var err error
	user, err := NewUser(uname, utils.SHA3hash(pwd), seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "failed to retrieve new user")
	}
	a.U = &user
	err = a.Save()
	if err != nil {
		log.Println(err)
		return a, err
	}
	if !consts.Mainnet {
		// automatically get funds if on testnet
		err = xlm.GetXLM(user.StellarWallet.PublicKey)
		if err != nil {
			log.Println("couildn't get xlm: ", err)
		}
	}
	return a, err
}

// SetOneTimeUnlock sets a one time seedpwd that can be used to
// automatically unlock the project once an investment comes in
func (a *Recipient) SetOneTimeUnlock(projIndex int, seedpwd string) error {
	log.Println("setting one time unlock for project with index: ", projIndex)
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve project")
	}

	recp, err := RetrieveRecipient(project.RecipientIndex)
	if err != nil {
		return errors.Wrap(err, "did not retrieve recipient belonging to project")
	}

	if recp.U.Index != a.U.Index {
		return errors.Wrap(err, "recipient index does not match with project recipient index")
	}

	project.OneTimeUnlock = seedpwd
	return project.Save()
}

// SetCompany sets the company bool to true
func (a *Recipient) SetCompany() error {
	a.Company = true
	return a.Save()
}

// SetCompanyDetails stores the company details in the recipient class
func (a *Recipient) SetCompanyDetails(companyType, name, legalName, adminEmail, phoneNumber, address,
	country, city, zipCode, taxIDNumber, role string) error {

	a.C.CompanyType = companyType
	a.C.Name = name
	a.C.LegalName = legalName
	a.C.AdminEmail = adminEmail
	a.C.PhoneNumber = phoneNumber
	a.C.Address = address
	a.C.Country = country
	a.C.City = city
	a.C.ZipCode = zipCode
	a.C.TaxIDNumber = taxIDNumber
	a.C.Role = role

	return a.Save()
}
