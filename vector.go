package vec

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// Vec is a vector of float64 coordinates
type Vec struct {
	Coords []float64
}

// NewVec returns a vector of the coordinates
func NewVec(coords []float64) *Vec {

	return &Vec{
		Coords: coords,
	}
}

// NewRandomVec generates a random vector
func NewRandomVec(dim int, valueRangeMin float64, valueRangeMax float64) *Vec {

	if valueRangeMin >= valueRangeMax || valueRangeMax <= 0 {
		panic("incorrect range parameters provided: min should be less than max")
	}

	vector := make([]float64, 0)
	for j := 0; j < dim; j++ {
		vector = append(vector, float64(rand.Intn(int(valueRangeMax-valueRangeMin+1)))+valueRangeMin)
	}

	return NewVec(vector)
}

// Equal returns true if a = b on every coordinate
func (a *Vec) Equal(b *Vec) bool {

	if len(a.Coords) != len(b.Coords) {
		return false
	}

	equal := true
	for i := 0; i < len(a.Coords); i++ {
		if a.Coords[i] != b.Coords[i] {
			equal = false
			break
		}
	}

	return equal
}

// ContainedIn returns true if the array of vectors arr contains the vector a
func (a *Vec) ContainedIn(arr []*Vec) bool {

	for i := 0; i < len(arr); i++ {
		if arr[i].Equal(a) {
			return true
		}
	}

	return false
}

// Dot returns the dot product of the two vectors a and b
func (a *Vec) Dot(b *Vec) (float64, error) {

	if len(a.Coords) != len(b.Coords) {
		return 0.0, errors.New("cannot take dot product of different sized vectors")
	}

	res := 0.0
	for i := 0; i < len(a.Coords); i++ {
		res += a.Coords[i] * b.Coords[i]
	}

	return res, nil
}

// Coord returns the value of the ith coordinate in a
func (a *Vec) Coord(i int) float64 {
	return a.Coords[i]
}

// SetValueToCoord assigns value of the ith coordinate in a
func (a *Vec) SetValueToCoord(value float64, i int) {
	a.Coords[i] = value
}

// AddToCoord add values to the ith coordinate of a
func (a *Vec) AddToCoord(value float64, i int) {
	a.Coords[i] += value
}

// Size returns the dimentionality of a
func (a *Vec) Size() int {
	return len(a.Coords)
}

// IsBinary returns true if all coordinates are 0 or 1
func (a *Vec) IsBinary() bool {
	for _, c := range a.Coords {
		if c != 0 && c != 1 {
			return false
		}
	}
	return true
}

// Copy returns a copy of the array a
func (a *Vec) Copy() *Vec {

	cpy := make([]float64, len(a.Coords))
	copy(cpy, a.Coords)
	return &Vec{
		Coords: cpy,
	}
}

// EuclideanDistance returns euclidean distance between two vectors p and q (must have the same dimentions)
func EuclideanDistance(p, q *Vec) float64 {

	if p.Size() != q.Size() {
		panic("points must have the same dimentions")
	}

	distance := 0.0
	for i := 0; i < len(p.Coords); i++ {
		distance += math.Pow(p.Coords[i]-q.Coords[i], 2)
	}

	return math.Sqrt(distance)
}

// HammingDistance returns the hamming distance between the two vectors (number of differing bits)
func HammingDistance(p, q *Vec) float64 {

	if p.Size() != q.Size() {
		panic("points must have the same dimentions")
	}

	if !p.IsBinary() || !q.IsBinary() {
		panic(fmt.Sprintf("non binary vectors when computing hamming distance %v %v", p.Coords, q.Coords))

	}

	distance := 0.0
	for i := 0; i < len(p.Coords); i++ {
		distance += math.Abs(p.Coords[i] - q.Coords[i])
	}

	return distance
}

// CosineDistance returns cosine distance between two vectors p and q (must have the same dimentions)
func CosineDistance(p, q *Vec) float64 {

	if p.Size() != q.Size() {
		panic("points must have the same dimentions")
	}

	distanceNum := 0.0
	distanceDenomP := 0.0
	distanceDenomQ := 0.0
	for i := 0; i < len(p.Coords); i++ {
		distanceNum += p.Coords[i] * q.Coords[i]
		distanceDenomP += p.Coords[i] * p.Coords[i]
		distanceDenomQ += q.Coords[i] * q.Coords[i]
	}

	distance := distanceNum / (math.Sqrt(distanceDenomP) * math.Sqrt(distanceDenomQ))

	return distance
}

// AbsoluteDifference computes the component wise absolute difference between the two vectors
func AbsoluteDifference(p, q *Vec) float64 {

	if p.Size() != q.Size() {
		panic("points must have the same dimentions")
	}

	distance := 0.0
	for i := 0; i < len(p.Coords); i++ {
		distance += math.Abs(p.Coords[i] - q.Coords[i])
	}

	return distance
}
