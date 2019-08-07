// +build all

package oracle

// this test function actually does nothing, since the oracle itself is a placeholder
// until we arrive at consensus on how it should be structured and stuff
import (
	"testing"
)

func BenchmarkMonthlyBill(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonthlyBill()
	}
}

func BenchmarkMonthlyBillInFloat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonthlyBillInFloat()
	}
}

func BenchmarkExchange(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExchangeXLMforUSD("1")
	}
}
