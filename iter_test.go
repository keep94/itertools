package itertools

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCycleEmpty(t *testing.T) {
	seq := Cycle[int]()
	for range seq {
		assert.FailNow(t, "Cycle should return empty sequence.")
	}
}

func TestChainCycle(t *testing.T) {
	startSeq := slices.Values([]int{1, 3, 5})
	seq := Chain(startSeq, Cycle(2, 4))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2, 4}, firstNOf(seq, 7))
	assert.Equal(t, []int{1, 3}, firstNOf(seq, 2))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2, 4, 2}, firstNOf(seq, 8))
}

func TestChainEmpty(t *testing.T) {
	seq := Chain[int]()
	for range seq {
		assert.FailNow(t, "Chain should return empty sequence.")
	}
}

func TestChainSingle(t *testing.T) {
	seq := Chain(slices.Values([]int{3, 5}))
	assert.Equal(t, []int{3, 5}, firstNOf(seq, 3))
	assert.Equal(t, []int{3, 5}, firstNOf(seq, 2))
	assert.Equal(t, []int{3}, firstNOf(seq, 1))
	assert.Equal(t, []int{3, 5}, firstNOf(seq, 0))
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

func TestCount(t *testing.T) {
	seq := Count(0, 1)
	assert.Equal(t, []int{0, 1, 2}, firstNOf(seq, 3))
	assert.Equal(t, []int{0, 1, 2, 3}, firstNOf(seq, 4))
}

func TestCount1(t *testing.T) {
	seq := Count(3, 5)
	assert.Equal(t, []int{3, 8, 13}, firstNOf(seq, 3))
	assert.Equal(t, []int{3, 8, 13, 18}, firstNOf(seq, 4))
}

func TestTake(t *testing.T) {
	seq := Count(10, 1)
	assert.Empty(t, slices.Collect(Take(seq, -1)))
	assert.Empty(t, slices.Collect(Take(seq, 0)))
	assert.Equal(t, []int{10}, slices.Collect(Take(seq, 1)))
	takeSeq := Take(seq, 3)
	assert.Equal(t, []int{10, 11, 12}, slices.Collect(takeSeq))
	assert.Equal(t, 10, firstOf(takeSeq))
}

func TestTakeFinite(t *testing.T) {
	seq := slices.Values([]string{"abc", "123", "foo"})
	takeSeq := Take(seq, 4)
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, "abc", firstOf(takeSeq))
}

func TestTakeWhile(t *testing.T) {
	seq := Count(10, 1)
	f := func(x int) bool { return x < 15 }
	g := func(x int) bool { return x < 10 }
	assert.Empty(t, slices.Collect(TakeWhile(seq, g)))
	takeSeq := TakeWhile(seq, f)
	assert.Equal(t, []int{10, 11, 12, 13, 14}, slices.Collect(takeSeq))
	assert.Equal(t, 10, firstOf(takeSeq))
}

func TestTakeWhileFinite(t *testing.T) {
	seq := slices.Values([]string{"abc", "123", "foo"})
	f := func(s string) bool { return len(s) < 4 }
	takeSeq := TakeWhile(seq, f)
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, "abc", firstOf(takeSeq))
}

func firstOf[T any](seq iter.Seq[T]) T {
	var result T
	for x := range seq {
		result = x
		break
	}
	return result
}

func firstNOf[T any](seq iter.Seq[T], n int) []T {
	var result []T
	for x := range seq {
		result = append(result, x)
		if len(result) == n {
			break
		}
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
