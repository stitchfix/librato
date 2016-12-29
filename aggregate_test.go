package librato

import "testing"

func TestSumSquares(t *testing.T) {
	a := Aggregate{}
	a.Add(12)
	a.Add(6)
	a.Add(5)
	a.Add(4)
	a.Add(5)
	a.Add(10)
	a.Add(3)

	if a.SumSquares()-65.714286 > 0.000001 {
		t.Errorf("Expected 65.714286, got %f", a.SumSquares())
	}
}
