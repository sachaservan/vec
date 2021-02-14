package vec

import (
	"testing"
)

func TestVecEqual(t *testing.T) {

	vecA := [...]float64{0, 1, 2}
	vecB := [...]float64{0, 1, 2}

	a := NewVec(vecA[:])
	b := NewVec(vecB[:])

	if !a.Equal(b) {
		t.Fail()
	}

	if !a.Equal(a) {
		t.Fail()
	}
}

func TestContains(t *testing.T) {

	vecA := [...]float64{0, 1, 2}
	vecB := [...]float64{0, 1, 2}
	vecC := [...]float64{0, 0, 0}
	vecD := [...]float64{1, 1, 1}

	a := NewVec(vecA[:])
	b := NewVec(vecB[:])
	c := NewVec(vecC[:])
	d := NewVec(vecD[:])

	vecArr := make([]*Vec, 2)
	vecArr[0] = b
	vecArr[1] = c

	if !a.ContainedIn(vecArr) {
		t.Fail()
	}

	if d.ContainedIn(vecArr) {
		t.Fail()
	}
}

func TestDotProduct(t *testing.T) {

	// dot product of A and B should be 0
	vecA := [...]float64{0, 1, 0}
	vecB := [...]float64{0, 0, 0}

	// dot product of C and D should be 2
	vecC := [...]float64{1, 1, 1}
	vecD := [...]float64{1, 0, 1}

	a := NewVec(vecA[:])
	b := NewVec(vecB[:])
	c := NewVec(vecC[:])
	d := NewVec(vecD[:])

	res, err := a.Dot(b)
	if err != nil || res != 0 {
		t.Fatalf("Incorrest result. Expected 0 got %v", res)
	}

	res, err = c.Dot(d)
	if err != nil || res != 2 {
		t.Fatalf("Incorrest result. Expected 2 got %v", res)
	}
}

func TestCoord(t *testing.T) {
	vecA := [...]float64{0, 1, 2}
	a := NewVec(vecA[:])

	if a.Coord(0) != 0 {
		t.Fail()
	}

	if a.Coord(1) != 1 {
		t.Fail()
	}

	if a.Coord(2) != 2 {
		t.Fail()
	}
}

func TestAddToCoord(t *testing.T) {
	vecA := [...]float64{0, 1, 2}
	a := NewVec(vecA[:])

	a.AddToCoord(1, 0)
	a.AddToCoord(1, 1)
	a.AddToCoord(1, 2)

	if a.Coord(0) != 1 {
		t.Fail()
	}

	if a.Coord(1) != 2 {
		t.Fail()
	}

	if a.Coord(2) != 3 {
		t.Fail()
	}
}

func TestEuclideanDistance(t *testing.T) {
	// distance between A and B should be sqrt(4)=2
	vecA := [...]float64{1, 1, 1, 1}
	vecB := [...]float64{0, 0, 0, 0}

	// distance between C and D should be sqrt(0)=0
	vecC := [...]float64{1, 1, 1}
	vecD := [...]float64{1, 1, 1}

	a := NewVec(vecA[:])
	b := NewVec(vecB[:])
	c := NewVec(vecC[:])
	d := NewVec(vecD[:])

	if EuclideanDistance(a, b) != 2 {
		t.Fail()
	}

	if EuclideanDistance(c, d) != 0 {
		t.Fail()
	}

}
