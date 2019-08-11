package oracle

import ()

// MonthlyBill returns the power tariffs for a month charged by the utility companies
func MonthlyBill() float64 {
	priceOfElectricity := 0.2
	averageConsumption := float64(600)
	return priceOfElectricity * averageConsumption
}
