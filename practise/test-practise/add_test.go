package test

import "testing"

func TestAdd(t *testing.T) {
	x, y := 1, 2
	z := Add(x, y)
	if z != 3 {
		t.Errorf("Add(%d,%d)=%d,real=%d", x, y, z, 3)
	}
}
