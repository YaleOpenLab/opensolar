package opensolar

// this should contain the future guarantor related functions once we define them concretely

// AddFirstLossGuarantee adds the given entity as a first loss guarantor
func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}
