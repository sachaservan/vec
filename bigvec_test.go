package vec

import (
	"testing"

	"github.com/ncw/gmp"
)

const dim = 100

func TestToBigVecBasic(t *testing.T) {
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {
		a := NewRandomVec(10, -10, 10)
		aBig := a.ToBigVec(scale)

		for i := 0; i < a.Size(); i++ {
			if aBig.Coords[i].Int64() != int64(a.Coords[i]) {
				t.Fatalf("Expected %v, got %v\n", int64(a.Coords[i]), aBig.Coords[i].Int64())
			}
		}
	}
}

func TestNewZeroVec(t *testing.T) {
	aBig := NewBigZeroVec(dim)
	for i := 0; i < len(aBig.Coords); i++ {
		if aBig.Coords[i].Int64() != 0 {
			t.Fatalf("Expected zero vector, got %v\n", aBig.Coords)
		}
	}
}

func TestNewRandomVector(t *testing.T) {

	for trial := 0; trial < 100; trial++ {
		numSameCoodinates := 0
		numNegativeCoordinates := 0

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		bBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		for i := 0; i < aBig.Size(); i++ {
			if aBig.Coords[i].Int64() == bBig.Coords[i].Int64() {
				numSameCoodinates++
			}

			if aBig.Coords[i].Int64() < 0 {
				numNegativeCoordinates++
			}
		}

		if numSameCoodinates > 5 {
			t.Fatalf(
				"Two random vectors share %v coordinates (very unlikely)\n a = %v\n b = %v",
				numSameCoodinates, aBig.Coords, bBig.Coords)
		}

		if numNegativeCoordinates > 2*dim/3 || numNegativeCoordinates < dim/3 {
			t.Fatalf("Random vector has %v negative coordinates (very unlikely)\n", numNegativeCoordinates)
		}
	}
}

func TestGetCoord(t *testing.T) {
	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))

		for i := 0; i < aBig.Size(); i++ {
			if aBig.Coords[i] != aBig.Coord(i) {
				t.Fatalf("Expected %v got %v\n", aBig.Coords[i], aBig.Coord(i))
			}
		}
	}
}

func TestClone(t *testing.T) {
	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		aBigCpy := aBig.Clone()

		for i := 0; i < aBig.Size(); i++ {
			if aBig.Coords[i].Int64() != aBigCpy.Coords[i].Int64() {
				t.Fatalf("Expected %v got %v\n", aBig.Coords[i], aBigCpy.Coords[i])
			}
		}
	}
}

func TestAdd(t *testing.T) {

	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		bBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		aBigCpy := aBig.Clone()

		sum, err := aBig.Add(bBig)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < aBig.Size(); i++ {
			expected := aBigCpy.Coords[i].Int64() + bBig.Coords[i].Int64()
			got := sum.Coords[i].Int64()
			if got != expected {
				t.Fatalf("Expected %v, got %v\n", expected, got)
			}
		}
	}
}

func TestSub(t *testing.T) {

	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		bBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		aBigCpy := aBig.Clone()

		sum, err := aBig.Sub(bBig)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < aBig.Size(); i++ {
			expected := aBigCpy.Coords[i].Int64() - bBig.Coords[i].Int64()
			got := sum.Coords[i].Int64()
			if got != expected {
				t.Fatalf("Expected %v, got %v\n", expected, got)
			}
		}
	}
}

func TestMul(t *testing.T) {

	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		bBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		aBigCpy := aBig.Clone()

		sum, err := aBig.Mul(bBig)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < aBig.Size(); i++ {
			expected := aBigCpy.Coords[i].Int64() * bBig.Coords[i].Int64()
			got := sum.Coords[i].Int64()
			if got != expected {
				t.Fatalf("Expected %v, got %v\n", expected, got)
			}
		}
	}
}

func TestDecodeSigned(t *testing.T) {

	field := randomPrime(100)

	for trial := 0; trial < 100; trial++ {
		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		aEncoded := aBig.Clone().Mod(field)
		aDecoded := aEncoded.DecodeSignedValues(field)

		if !aBig.Equal(aDecoded) {
			t.Fatalf("Expected %v, got %v\n", aBig, aDecoded)
		}
	}
}

func TestMod(t *testing.T) {

	for trial := 0; trial < 100; trial++ {

		field := randomPrime(10)

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		mod := aBig.Mod(field)

		for i := 0; i < aBig.Size(); i++ {
			expected := aBig.Coords[i].Int64() % field.Int64()
			if expected < 0 {
				expected = field.Int64() - expected
			}

			got := mod.Coords[i].Int64()
			if got != expected {
				t.Fatalf("Expected %v, got %v\n", expected, got)
			}
		}
	}
}

func TestBigVecEqual(t *testing.T) {

	for trial := 0; trial < 100; trial++ {

		aBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))
		bBig := NewBigRandomVec(dim, gmp.NewInt(-1000), gmp.NewInt(1000))

		if aBig.Equal(bBig) {
			t.Fatalf("Two random vectors are equal (very unlikely)\n a = %v\n b = %v\n", aBig, bBig)
		}

		if !aBig.Equal(aBig) {
			t.Fatalf("Vector does not equal itself!")
		}
	}
}
