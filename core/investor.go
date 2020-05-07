package core

import (
	"log"

	"github.com/pkg/errors"

	tickers "github.com/Varunram/essentials/exchangetickers"
	utils "github.com/Varunram/essentials/utils"
	xlm "github.com/Varunram/essentials/xlm"
	openxconsts "github.com/YaleOpenLab/openx/consts"
	openx "github.com/YaleOpenLab/openx/database"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

// Investor defines the investor structure
type Investor struct {
	// U is the base User class inherited from openx
	U *openx.User

	// C is a structure containing all details of the company the investor is part of
	C Company

	// Company denotes whether the given investor is acting on behalf of a company
	Company bool

	// VotingBalance is the balance associated with the particular investor (equal to the amount of USD he possesses)
	VotingBalance float64

	// AmountInvested is the total amount invested by the investor
	AmountInvested float64

	// InvestedSolarProjects is a list of the investor assets of the opensolar projects the investor has invested in
	InvestedSolarProjects []string

	// InvestedSolarProjectsIndices is an integer list of the projects the investor has invested in
	InvestedSolarProjectsIndices []int

	// InvestedSolarProjects is a list of the investor assets of the opensolar projects the investor has invested in
	SeedInvestedSolarProjects []string

	// InvestedSolarProjectsIndices is an integer list of the projects the investor has invested in
	SeedInvestedSolarProjectsIndices []int
}

// Company is a struct that is used if an investor/recipient is acting on behalf of their company
type Company struct {
	// CompanyType is the type of the company
	CompanyType string

	// Name is the name of the company
	Name string

	// LegalName is the legal name of the company
	LegalName string

	// AdminEmail is the email of the admin / contact point of the company
	AdminEmail string

	// PhoneNumber is the phone number of the main contact in the company
	PhoneNumber string

	// Address is the registered address of the company
	Address string

	// Country is the country where the company is registered in
	Country string

	// City is the city in which the company is registered at
	City string

	// ZipCode is the zipcode of the city where the company is at
	ZipCode string

	// TaxIDNumber is the tax id number associated with the company
	TaxIDNumber string

	// Role isthe role of the investor in the above company
	Role string
}

// NewInvestor creates a new investor
func NewInvestor(uname string, pwd string, seedpwd string, Name string) (Investor, error) {
	var a Investor
	var err error
	user, err := NewUser(uname, utils.SHA3hash(pwd), seedpwd, Name)
	if err != nil {
		return a, errors.Wrap(err, "error while creating a new user")
	}
	a.U = &user
	a.AmountInvested = -1.0
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

// ChangeVotingBalance changes the voting balance of a user
func (a *Investor) ChangeVotingBalance(votes float64) error {
	// this function is caled when we want to refund the user with their votes once
	// an order has been finalized.
	a.VotingBalance += votes
	if a.VotingBalance < 0 {
		a.VotingBalance = 0 // to ensure no one has negative votes
	}
	return a.Save()
}

// CanInvest checks whether an investor has the required funds to invest in a project
func (a *Investor) CanInvest(targetBalance float64) bool {

	// if !a.U.Legal {
	// 	log.Println("user has not accepted terms and conditions associated with the platform")
	// 	return false
	// }

	if !consts.Mainnet {
		usdBalance := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, "STABLEUSD")
		xlmBalance := xlm.GetNativeBalance(a.U.StellarWallet.PublicKey)
		// need to fetch the oracle price here for the order
		oraclePrice := tickers.ExchangeXLMforUSD(xlmBalance)
		return usdBalance > targetBalance+1 || oraclePrice > targetBalance
	}

	// mainnet
	usdBalance := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, openxconsts.AnchorUSDCode)
	return usdBalance > targetBalance+1
}

// SetCompany sets the company bool to true. This enables the creation of investors who
// can act on behalf of a company
func (a *Investor) SetCompany() error {
	a.Company = true
	return a.Save()
}

// SetCompanyDetails sets the company struct details of the investor class
func (a *Investor) SetCompanyDetails(companyType, name, legalName, adminEmail, phoneNumber, address,
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
