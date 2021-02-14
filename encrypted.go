package vec

import (
	crand "crypto/rand"
	"errors"
	"log"
	"math/big"
	"math/rand"

	"github.com/ncw/gmp"
	"github.com/sachaservan/paillier"

	"time"
)

// EncryptedVec is a vector of encrypted float64 coordinates
type EncryptedVec struct {
	Pk     *paillier.PublicKey
	Coords []*paillier.Ciphertext
}

// NewEncryptVecWithCoords returns a new encrypted vector with the ciphertexts as coordinates
func NewEncryptVecWithCoords(pk *paillier.PublicKey, coords []*paillier.Ciphertext) *EncryptedVec {
	return &EncryptedVec{
		Pk:     pk,
		Coords: coords,
	}
}

// NewEncryptedVec constructs a share of a vector
func NewEncryptedVec(pk *paillier.PublicKey, coords []*paillier.Ciphertext) *EncryptedVec {
	return &EncryptedVec{
		Coords: coords,
		Pk:     pk,
	}
}

// Encrypt returns an encryption of the vector
func Encrypt(a *BigVec, pk *paillier.PublicKey) *EncryptedVec {
	encrypted := make([]*paillier.Ciphertext, len(a.Coords))

	for i, coord := range a.Coords {
		encrypted[i] = pk.Encrypt(coord)
	}

	return &EncryptedVec{
		Pk:     pk,
		Coords: encrypted,
	}
}

// GetCoords returns the big vector of coordinates
func (a *EncryptedVec) GetCoords() []*paillier.Ciphertext {
	return a.Coords
}

// SetCoords sets the coordinates to the big vector
func (a *EncryptedVec) SetCoords(v []*paillier.Ciphertext) {
	a.Coords = v
}

// GetCoord returns the coordinate at the index
func (a *EncryptedVec) GetCoord(i int) *paillier.Ciphertext {
	return a.Coords[i]
}

// Size returns the dimentionality of the vector
func (a *EncryptedVec) Size() int {
	return len(a.Coords)
}

// Add returns the component-wise addition of a and b
// throws an error if the vectors are of different size
func (a *EncryptedVec) Add(b *EncryptedVec) (*EncryptedVec, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot add vectors of different length")
	}

	pk := a.Pk
	res := make([]*paillier.Ciphertext, len(a.Coords))

	for i := range a.Coords {
		res[i] = pk.Add(a.Coords[i], b.Coords[i])
	}

	return &EncryptedVec{
		Pk:     pk,
		Coords: res,
	}, nil
}

// Sub returns the component-wise subtaction of a and b
// throws an error if the vectors are of different size
func (a *EncryptedVec) Sub(b *EncryptedVec) (*EncryptedVec, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot add vectors of different length")
	}

	pk := a.Pk
	res := make([]*paillier.Ciphertext, len(a.Coords))

	for i := range a.Coords {
		res[i] = pk.Sub(a.Coords[i], b.Coords[i])
	}

	return &EncryptedVec{
		Pk:     pk,
		Coords: res,
	}, nil
}

// Dot returns the (encrypted) dot product of the two vectors a and b
// using the homomorphic encryption property of the encrypted vector
func (a *EncryptedVec) Dot(b *BigVec, pk *paillier.PublicKey) (*paillier.Ciphertext, error) {

	if len(a.Coords) != len(b.Coords) {
		return nil, errors.New("cannot take dot product of different sized vectors")
	}

	res := pk.EncryptZero()
	for i := 0; i < len(a.Coords); i++ {
		res = pk.Add(res, pk.ConstMult(a.Coords[i], b.Coords[i]))
	}

	return res, nil
}

func shuffle(vals []*paillier.Ciphertext) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

// generates a new random number < max
func newCryptoRandom(max *gmp.Int) *gmp.Int {

	maxInt := new(big.Int).SetBytes(max.Bytes())
	rand, err := crand.Int(crand.Reader, maxInt)
	if err != nil {
		log.Println(err)
	}

	return new(gmp.Int).SetBytes(rand.Bytes())
}
