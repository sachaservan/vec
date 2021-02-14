package vec

import (
	"crypto/rand"
	"testing"

	"github.com/ncw/gmp"
)

func TestSecretShare(t *testing.T) {

	field := randomPrime(100)
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {

		vecA := [...]float64{0, 1, 0}
		a := NewVec(vecA[:])
		aBig := NewBigVec(a.ToBigVec(scale).Coords)
		t.Logf("aBig = %v\n", aBig.Coords)
		sharesA := SecretShare(aBig, 2, field)
		t.Logf("sharesA[0] = %v\n", sharesA[0].GetCoords())
		t.Logf("sharesA[1] = %v\n", sharesA[1].GetCoords())

		recovered, err := RecoverVector(sharesA...)
		if err != nil || !recovered.Equal(aBig) {
			t.Fatalf("Incorrest result. Expected %v got %v", a.Coords, recovered.Coords)
		}
	}
}

func TestSecretShareAdd(t *testing.T) {

	field := randomPrime(100)
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {

		a := NewRandomVec(dim, -100, 100)
		b := NewRandomVec(dim, -100, 100)
		aBig := a.ToBigVec(scale)
		bBig := b.ToBigVec(scale)
		sharesA := SecretShare(aBig, 2, field)
		sharesB := SecretShare(bBig, 2, field)

		res0, err := sharesA[0].Add(sharesB[0])
		res1, err := sharesA[1].Add(sharesB[1])

		res, err := RecoverVector(res0, res1)
		if err != nil {
			t.Fatal(err)
		}

		expected, _ := aBig.Add(bBig)
		if !res.Equal(expected) {
			t.Fatalf("Incorrest result. \nExpected %v \nGot %v", expected, res)
		}
	}
}

func TestSecretShareSub(t *testing.T) {

	field := randomPrime(100)
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {

		a := NewRandomVec(dim, -100, 100)
		b := NewRandomVec(dim, -100, 100)
		aBig := a.ToBigVec(scale)
		bBig := b.ToBigVec(scale)
		sharesA := SecretShare(aBig, 2, field)
		sharesB := SecretShare(bBig, 2, field)

		res0, err := sharesA[0].Sub(sharesB[0])
		res1, err := sharesA[1].Sub(sharesB[1])

		res, err := RecoverVector(res0, res1)
		if err != nil {
			t.Fatal(err)
		}

		expected, _ := aBig.Sub(bBig)
		if !res.Equal(expected) {
			t.Fatalf("Incorrest result. \nExpected %v \nGot %v", expected, res)
		}
	}
}

func TestSecretShareMul(t *testing.T) {

	field := randomPrime(100)
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {

		a := NewRandomVec(dim, -100, 100)
		b := NewRandomVec(dim, -100, 100)
		aBig := a.ToBigVec(scale)
		bBig := b.ToBigVec(scale)
		sharesA := SecretShare(aBig, 2, field)

		res0, err := sharesA[0].Mul(bBig)
		res1, err := sharesA[1].Mul(bBig)

		res, err := RecoverVector(res0, res1)
		if err != nil {
			t.Fatal(err)
		}

		expected, _ := aBig.Mul(bBig)
		if !res.Equal(expected) {
			t.Fatalf("Incorrest result. \nExpected %v \nGot %v", expected, res)
		}
	}
}

func TestSecretSharedDotProduct(t *testing.T) {

	field := randomPrime(100)
	scale := gmp.NewInt(1)

	for trial := 0; trial < 100; trial++ {

		a := NewRandomVec(dim, -100, 100)
		b := NewRandomVec(dim, -100, 100)

		sharesA := SecretShare(a.ToBigVec(scale), 2, field)

		res0, err := sharesA[0].Dot(b.ToBigVec(scale))
		res1, err := sharesA[1].Dot(b.ToBigVec(scale))
		got := RecoverInt(field, res0, res1)
		expected, _ := a.Dot(b)

		if err != nil || float64(got.Int64()) != expected {
			t.Fatalf("Incorrest result. Expected %v, got %v", expected, got)
		}
	}
}

func randomPrime(bits int) *gmp.Int {
	for {
		p, err := rand.Prime(rand.Reader, bits)
		if err != nil {
			continue
		} else {
			return new(gmp.Int).SetBytes(p.Bytes())
		}
	}
}
