package vec

import (
	"errors"

	"github.com/ncw/gmp"
)

// ShareVec is a secret share of a vector
type ShareVec struct {
	Vec   *BigVec
	P     *gmp.Int
	Index int // share number
}

// SecretShare returns secret shares of the vector where p is a prime modulus
func SecretShare(a *BigVec, numShares int, p *gmp.Int) []*ShareVec {

	// run 20 tests of Rabin-Miller
	if !p.ProbablyPrime(20) {
		panic("trying to secret share in a non-prime order field")
	}

	dim := len(a.Coords)

	sum := NewBigZeroVec(dim)
	shares := make([]*ShareVec, numShares)

	for i := 0; i < numShares-1; i++ {
		shares[i] = &ShareVec{}
		shares[i].Vec = NewBigRandomVec(dim, gmp.NewInt(0), p)
		shares[i].Index = i
		shares[i].P = p
		sum.Sub(shares[i].Vec)
	}

	// add the secret to the last coordinate
	for i, coord := range a.Coords {
		sum.Coords[i].Add(sum.Coords[i], coord)
	}

	sum.Mod(p)

	shares[numShares-1] = &ShareVec{
		sum,
		p,
		numShares - 1,
	}

	return shares
}

// NewShareVec constructs a share of a vector
func NewShareVec(coords []*gmp.Int, p *gmp.Int) *ShareVec {
	return &ShareVec{
		Vec: NewBigVec(coords),
		P:   p,
	}
}

// GetCoords returns the big vector of coordinates
func (a *ShareVec) GetCoords() []*gmp.Int {
	return a.Vec.Coords
}

// SetCoords sets the coordinates to the big vector
func (a *ShareVec) SetCoords(v *BigVec) {
	v.Mod(a.P)
	a.Vec = v
}

// RecoverVector outputs the recovered BigVec from the set of secret shares
func RecoverVector(shares ...*ShareVec) (*BigVec, error) {

	dim := len(shares[0].Vec.Coords)
	p := shares[0].P
	res := NewBigZeroVec(dim)

	for _, share := range shares {
		if _, err := res.Add(share.Vec); err != nil {
			return nil, err
		}
	}

	res = res.Mod(p)
	res = res.DecodeSignedValues(p)

	return res, nil
}

// RecoverInt returns a an integer encoded in the shares
func RecoverInt(p *gmp.Int, shares ...*gmp.Int) *gmp.Int {
	res := new(gmp.Int)
	for _, share := range shares {
		res.Add(res, share)
	}

	// decode the sign
	negThresh := new(gmp.Int).Quo(p, gmp.NewInt(2))
	if res.Cmp(negThresh) > 0 {
		res = new(gmp.Int).Sub(res, p)
	}

	return res
}

// FieldModulus returns the field modulus p
func (a *ShareVec) FieldModulus() *gmp.Int {
	return a.P
}

// Add returns the component-wise addition of a and b
// throws an error if the vectors are of different size
func (a *ShareVec) Add(b *ShareVec) (*ShareVec, error) {
	c, err := a.Vec.Add(b.Vec)
	if err != nil {
		return nil, err
	}

	if a.Index != b.Index {
		return nil, errors.New("Index of share a != index of share b")
	}

	c = c.Mod(a.P)

	return &ShareVec{c, a.P, a.Index}, nil
}

// Sub returns the component-wise subtaction of a and b
// throws an error if the vectors are of different size
func (a *ShareVec) Sub(b *ShareVec) (*ShareVec, error) {

	c, err := a.Vec.Sub(b.Vec)
	if err != nil {
		return nil, err
	}

	if a.Index != b.Index {
		return nil, errors.New("Index of share a != index of share b")
	}

	c = c.Mod(a.P)

	return &ShareVec{c, a.P, a.Index}, nil
}

// Mul returns the component-wise multiplication of a and b
// throws an error if the vectors are of different size
func (a *ShareVec) Mul(b *BigVec) (*ShareVec, error) {

	c, err := a.Vec.Mul(b)
	if err != nil {
		return nil, err
	}

	c = c.Mod(a.P)

	return &ShareVec{c, a.P, a.Index}, nil
}

// Dot returns the (encrypted) dot product of the two vectors a and b
// using the homomorphic encryption property of the encrypted vector
func (a *ShareVec) Dot(b *BigVec) (*gmp.Int, error) {

	if len(a.Vec.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	res := gmp.NewInt(0)
	for i := 0; i < len(a.Vec.Coords); i++ {
		prod := new(gmp.Int).Mul(a.Vec.Coords[i], b.Coords[i])
		res.Add(res, prod)
	}

	res.Mod(res, a.P)

	return res, nil
}
