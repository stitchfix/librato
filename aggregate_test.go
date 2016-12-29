package librato

import "testing"

func TestSumSquares(t *testing.T) {
	a := Aggregate{}
	a.Add(12).Add(6).Add(5).Add(4).Add(5).Add(10).Add(3)

	if a.SumSquares()-65.714286 > 0.000001 {
		t.Errorf("Expected 65.714286, got %f", a.SumSquares())
	}
}
