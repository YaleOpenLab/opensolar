package core

// AddFirstLossGuarantee adds the given entity as a first loss guarantor
func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}
