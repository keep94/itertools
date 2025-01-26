package itertools

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCyclePanic(t *testing.T) {
	assert.Panics(t, func() { Cycle[int]() })
}

func TestMap(t *testing.T) {
	y := []string{"four", "fives", "sixsix"}
	m := func(s string) int { return len(s) }
	it := Map(slices.Values(y), m)
	assert.Equal(t, []int{4, 5, 6}, slices.Collect(it))
	assert.Equal(t, 4, firstOf(it))
}

func TestFilter(t *testing.T) {
	x := []int{3, 4, 5, 6}
	f := func(i int) bool { return i%2 == 1 }
	it := Filter(slices.Values(x), f)
	assert.Equal(t, []int{3, 5}, slices.Collect(it))
	assert.Equal(t, 3, firstOf(it))
}

func TestPairUp(t *testing.T) {
	x := []int{1, 2, 3}
	y := []string{"four", "fives", "sixsix", "sevenen", "eightate"}
	pairIter := PairUp(Zip(slices.Values(x), slices.Values(y)))
	var prod int
	for p := range pairIter {
		prod = p.First * len(p.Second)
		break
	}
	assert.Equal(t, 4, prod)
	var z []int
	for p := range pairIter {
		z = append(z, p.First*len(p.Second))
	}
	assert.Equal(t, []int{4, 10, 18}, z)
}

func TestZipNormal(t *testing.T) {
	x := []int{1, 2, 3}
	y := []string{"four", "fives", "sixsix", "sevenen", "eightate"}
	zipIter := Zip(slices.Values(x), slices.Values(y))
	var z []int
	for i, j := range zipIter {
		z = append(z, i*len(j))
	}
	assert.Equal(t, []int{4, 10, 18}, z)
	z = nil
	for i, j := range zipIter {
		z = append(z, i*len(j))
	}
	assert.Equal(t, []int{4, 10, 18}, z)
}

func TestZipExitEarly(t *testing.T) {
	x := []int{1, 2, 3}
	y := []string{"four", "fives", "sixsix", "sevenen", "eightate"}
	zipIter := Zip(slices.Values(x), slices.Values(y))
	var prod int
	for i, j := range zipIter {
		prod = i * len(j)
		break
	}
	assert.Equal(t, 4, prod)
	var z []int
	for i, j := range zipIter {
		z = append(z, i*len(j))
	}
	assert.Equal(t, []int{4, 10, 18}, z)
}

func TestZipLazy(t *testing.T) {
	x := []int{1, 2, 3}
	y := []int{4, 5, 6, 7, 8}
	zipIter := Zip(values(x), values(y))
	var z []int
	for i, j := range zipIter {
		z = append(z, i*j)
	}
	assert.Equal(t, []int{4, 10, 18}, z)
	assert.Equal(t, []int{0, 0, 0}, x)

	// yield for y starts returning false at 7 because the x's are
	// exhausted
	assert.Equal(t, []int{0, 0, 0, 7, 8}, y)
}

func TestZipLazyExitEarly(t *testing.T) {
	x := []int{1, 2, 3}
	y := []int{4, 5, 6, 7, 8}
	zipIter := Zip(values(x), values(y))
	var prod int
	for i, j := range zipIter {
		prod = i * j
		break
	}
	assert.Equal(t, 4, prod)

	// yield for x returns false at 1.
	assert.Equal(t, []int{1, 2, 3}, x)

	// yield for y returns false at 4.
	assert.Equal(t, []int{4, 5, 6, 7, 8}, y)
}

func TestZipMisbehaved(t *testing.T) {
	x := []int{1, 2, 3}
	y := []int{4, 5, 6, 7, 8}
	zipIter := Zip(misbehaved(x), misbehaved(y))
	var prod int
	for i, j := range zipIter {
		prod = i * j
		break
	}
	assert.Equal(t, 4, prod)
}

func firstOf[T any](seq iter.Seq[T]) T {
	var result T
	for x := range seq {
		result = x
		break
	}
	return result
}

// values returns an iter.Seq over s that also sets s elements to zero until
// yield return false.
func values(s []int) iter.Seq[int] {
	return func(yield func(x int) bool) {
		for i := range s {
			if !yield(s[i]) {
				break
			}
			s[i] = 0
		}
	}
}

// misbehaved returns an iter.Seq over s that never checks to see if yield
// returns false.
func misbehaved(s []int) iter.Seq[int] {
	return func(yield func(x int) bool) {
		for i := range s {
			yield(s[i])
		}
	}
}
