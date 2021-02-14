package vec

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/ncw/gmp"
)

// BigVec is a vector of float64 coordinates
type BigVec struct {
	Coords []*gmp.Int
}

// NewBigVec returns a vector of the coordinates
func NewBigVec(coords []*gmp.Int) *BigVec {
	return &BigVec{
		Coords: coords,
	}
}

// ToBigVec converts a vector to a BigVec with fixedPoint encoding
func (v *Vec) ToBigVec(fpScaleFactor *gmp.Int) *BigVec {

	fpScaleFloat := new(big.Float).SetInt(new(big.Int).SetBytes(fpScaleFactor.Bytes()))
	vector := make([]*gmp.Int, 0)
	for j := 0; j < len(v.Coords); j++ {
		e := big.NewFloat(v.Coords[j])
		e.Mul(e, fpScaleFloat)
		eInt, _ := e.Int(big.NewInt(0))
		eGmp := new(gmp.Int).SetBytes(eInt.Bytes())
		eGmp.Mul(eGmp, gmp.NewInt(int64(eInt.Sign())))
		vector = append(vector, eGmp)
	}

	return NewBigVec(vector)
}

// NewBigZeroVec generates a new all-zero vector
func NewBigZeroVec(dim int) *BigVec {

	vector := make([]*gmp.Int, 0)
	for j := 0; j < dim; j++ {
		vector = append(vector, gmp.NewInt(0))
	}

	return NewBigVec(vector)
}

// NewBigRandomVec generates a random vector of big integers
func NewBigRandomVec(dim int, valueRangeMin *gmp.Int, valueRangeMax *gmp.Int) *BigVec {

	if valueRangeMax.Cmp(valueRangeMin) <= 0 || valueRangeMax.Cmp(gmp.NewInt(0)) <= 0 {
		panic("incorrect range parameters provided: min should be less than max")
	}

	valMinInt := new(big.Int).SetBytes(valueRangeMin.Bytes())
	valMinInt.Mul(valMinInt, big.NewInt(int64(valueRangeMin.Sign())))

	bound := new(big.Int).SetBytes(valueRangeMax.Bytes())
	bound.Sub(bound, valMinInt)
	bound.Add(bound, big.NewInt(1))

	vector := make([]*gmp.Int, 0)
	for j := 0; j < dim; j++ {
		rand, _ := rand.Int(rand.Reader, bound)
		randGmp := new(gmp.Int).SetBytes(rand.Bytes())
		randGmp.Add(randGmp, valueRangeMin)

		vector = append(vector, randGmp)
	}

	return NewBigVec(vector)
}

// Dot returns the dot product of the two vectors a and b
func (a *BigVec) Dot(b *BigVec) (*gmp.Int, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	res := gmp.NewInt(0)
	for i := 0; i < len(a.Coords); i++ {
		prod := gmp.NewInt(0).Mul(a.Coords[i], b.Coords[i])
		res.Add(res, prod)
	}

	return res, nil
}

// Coord returns the coordinate at the index
func (a *BigVec) Coord(i int) *gmp.Int {
	return a.Coords[i]
}

// Size returns the dimentionality of the vector
func (a *BigVec) Size() int {
	return len(a.Coords)
}

// Add returns the coordinate-wise sum of a and b
func (a *BigVec) Add(b *BigVec) (*BigVec, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	for i := 0; i < len(a.Coords); i++ {
		a.Coords[i].Add(a.Coords[i], b.Coords[i])
	}

	return a, nil
}

// Sub returns the coordinate-wise difference of a and b
func (a *BigVec) Sub(b *BigVec) (*BigVec, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	for i := 0; i < len(a.Coords); i++ {
		a.Coords[i].Sub(a.Coords[i], b.Coords[i])
	}

	return a, nil
}

// Mul returns the coordinate-wise multiplication of a and b
func (a *BigVec) Mul(b *BigVec) (*BigVec, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	for i := 0; i < len(a.Coords); i++ {
		a.Coords[i].Mul(a.Coords[i], b.Coords[i])
	}

	return a, nil
}

// Mod reduces each coordinate modulo n
func (a *BigVec) Mod(n *gmp.Int) *BigVec {

	zero := gmp.NewInt(0)
	for i := 0; i < len(a.Coords); i++ {
		a.Coords[i].Mod(a.Coords[i], n)
		if a.Coords[i].Cmp(zero) < 0 {
			a.Coords[i].Add(n, a.Coords[i])
		}
	}

	return a
}

// DecodeSignedValues returns signed coordinates from an encoding in Z_n
// all values  > n/2 treated as a negative value
func (a *BigVec) DecodeSignedValues(n *gmp.Int) *BigVec {

	signedVec := NewBigZeroVec(a.Size())
	negThresh := new(gmp.Int).Quo(n, gmp.NewInt(2))
	for i := 0; i < len(a.Coords); i++ {
		a.Coords[i].Mod(a.Coords[i], n)
		if a.Coords[i].Cmp(negThresh) > 0 {
			signedVec.Coords[i] = new(gmp.Int).Sub(a.Coords[i], n)
		} else {
			signedVec.Coords[i] = new(gmp.Int).Set(a.Coords[i])
		}
	}

	return signedVec
}

// Equal returns true if a = b on every coordinate
func (a *BigVec) Equal(b *BigVec) bool {

	if len(a.Coords) != len(b.Coords) {
		return false
	}

	equal := true
	for i := 0; i < len(a.Coords); i++ {
		if a.Coords[i].Cmp(b.Coords[i]) != 0 {
			equal = false
			break
		}
	}

	return equal
}

// Clone returns a copy of a
func (a *BigVec) Clone() *BigVec {

	newCoords := make([]*gmp.Int, a.Size())
	for i := 0; i < a.Size(); i++ {
		newCoords[i] = new(gmp.Int).Set(a.Coords[i])
	}

	return NewBigVec(newCoords)
}
