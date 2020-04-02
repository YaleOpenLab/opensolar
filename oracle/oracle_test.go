// +build all travis

package oracle

// this test function actually does nothing, since the oracle itself is a placeholder
// until we arrive at a consensus on how it should be structured
import (
	"testing"
)

func TestOracle(t *testing.T) {
	billF := MonthlyBill()
	if billF != 120.0 {
		t.Fatalf("Oracle does not output constant value")
	}
}
