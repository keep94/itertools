package itertools

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCycleValuesEmpty(t *testing.T) {
	seq := CycleValues[int]()
	for range seq {
		assert.FailNow(t, "CycleValues should return empty sequence.")
	}
}

func TestChainCycleValues(t *testing.T) {
	startSeq := slices.Values([]int{1, 3, 5})
	seq := Chain(startSeq, CycleValues(2, 4))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2, 4}, firstNOf(7, seq))
	assert.Equal(t, []int{1, 3}, firstNOf(2, seq))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2, 4, 2}, firstNOf(8, seq))
}

func TestCycleEmpty(t *testing.T) {
	seq := Cycle(slices.Values(([]int)(nil)))
	for range seq {
		assert.FailNow(t, "Cycle should return empty sequence.")
	}
}

func TestChainCycle(t *testing.T) {
	startSeq := slices.Values([]int{1, 3, 5})
	seq := Chain(startSeq, Cycle(slices.Values([]int{2, 4})))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2, 4, 2, 4}, firstNOf(9, seq))
	assert.Equal(t, []int{1, 3}, firstNOf(2, seq))
	assert.Equal(t, []int{1, 3, 5}, firstNOf(3, seq))
	assert.Equal(t, []int{1, 3, 5, 2}, firstNOf(4, seq))
	assert.Equal(t, []int{1, 3, 5, 2, 4}, firstNOf(5, seq))
	assert.Equal(t, []int{1, 3, 5, 2, 4, 2}, firstNOf(6, seq))
}

func TestChainEmpty(t *testing.T) {
	seq := Chain[int]()
	for range seq {
		assert.FailNow(t, "Chain should return empty sequence.")
	}
}

func TestChainSingle(t *testing.T) {
	seq := Chain(slices.Values([]int{3, 5}))
	assert.Equal(t, []int{3, 5}, firstNOf(3, seq))
	assert.Equal(t, []int{3, 5}, firstNOf(2, seq))
	assert.Equal(t, []int{3}, firstNOf(1, seq))
	assert.Equal(t, []int{3, 5}, firstNOf(0, seq))
}

func TestMap(t *testing.T) {
	y := []string{"four", "fives", "sixsix"}
	m := func(s string) int { return len(s) }
	it := Map(m, slices.Values(y))
	assert.Equal(t, []int{4, 5, 6}, slices.Collect(it))
	assert.Equal(t, 4, firstOf(it))
}

func TestFilter(t *testing.T) {
	x := []int{3, 4, 5, 6}
	f := func(i int) bool { return i%2 == 1 }
	it := Filter(f, slices.Values(x))
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
	assert.Equal(t, []int{0, 1, 2}, firstNOf(3, seq))
	assert.Equal(t, []int{0, 1, 2, 3}, firstNOf(4, seq))
}

func TestCount1(t *testing.T) {
	seq := Count(3, 5)
	assert.Equal(t, []int{3, 8, 13}, firstNOf(3, seq))
	assert.Equal(t, []int{3, 8, 13, 18}, firstNOf(4, seq))
}

func TestTake(t *testing.T) {
	cc := newCountConsumed(Count(10, 1))
	assert.Empty(t, slices.Collect(Take(-1, cc.All())))
	assert.Equal(t, 0, cc.Count())
	assert.Empty(t, slices.Collect(Take(0, cc.All())))
	assert.Equal(t, 0, cc.Count())
	assert.Equal(t, []int{10}, slices.Collect(Take(1, cc.All())))
	assert.Equal(t, 1, cc.Count())
	takeSeq := Take(3, cc.All())
	assert.Equal(t, []int{10, 11, 12}, slices.Collect(takeSeq))
	assert.Equal(t, 3, cc.Count())
	assert.Equal(t, 10, firstOf(takeSeq))
	assert.Equal(t, 1, cc.Count())
}

func TestTakeFinite(t *testing.T) {
	cc := newCountConsumed(slices.Values([]string{"abc", "123", "foo"}))
	takeSeq := Take(4, cc.All())
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, 3, cc.Count())
	assert.Equal(t, "abc", firstOf(takeSeq))
	assert.Equal(t, 1, cc.Count())
}

func TestTakeWhile(t *testing.T) {
	cc := newCountConsumed(
		slices.Values([]int{10, 11, 12, 13, 14, 15, 1, 2, 3, 4}))
	f := func(x int) bool { return x < 15 }
	g := func(x int) bool { return x < 10 }
	assert.Empty(t, slices.Collect(TakeWhile(g, cc.All())))
	assert.Equal(t, 1, cc.Count())
	takeSeq := TakeWhile(f, cc.All())
	assert.Equal(t, []int{10, 11, 12, 13, 14}, slices.Collect(takeSeq))
	assert.Equal(t, 6, cc.Count())
	assert.Equal(t, 10, firstOf(takeSeq))
	assert.Equal(t, 1, cc.Count())
}

func TestTakeWhileFinite(t *testing.T) {
	cc := newCountConsumed(slices.Values([]string{"abc", "123", "foo"}))
	f := func(s string) bool { return len(s) < 4 }
	takeSeq := TakeWhile(f, cc.All())
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, 3, cc.Count())
	assert.Equal(t, "abc", firstOf(takeSeq))
	assert.Equal(t, 1, cc.Count())
}

func TestDropWhile(t *testing.T) {
	seq := slices.Values([]int{10, 13, 16, 1, 2, 3})
	f := func(x int) bool { return x < 10 }
	g := func(x int) bool { return x < 13 }
	h := func(x int) bool { return x < 16 }
	dropSeq := DropWhile(f, seq)
	assert.Equal(t, []int{10, 13, 16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 10, firstOf(dropSeq))
	dropSeq = DropWhile(g, seq)
	assert.Equal(t, []int{13, 16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 13, firstOf(dropSeq))
	dropSeq = DropWhile(h, seq)
	assert.Equal(t, []int{16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 16, firstOf(dropSeq))
}

func TestDropWhileFinite(t *testing.T) {
	seq := slices.Values([]int{10, 13, 16})
	f := func(x int) bool { return x < 19 }
	dropSeq := DropWhile(f, seq)
	assert.Empty(t, slices.Collect(dropSeq))
	assert.Empty(t, slices.Collect(dropSeq))
}

func TestAt(t *testing.T) {
	cc := newCountConsumed(slices.Values([]string{"a", "b", "c"}))
	val, ok := At(-1, cc.All())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 0, cc.Count())
	val, ok = At(0, cc.All())
	assert.Equal(t, "a", val)
	assert.True(t, ok)
	assert.Equal(t, 1, cc.Count())
	val, ok = At(1, cc.All())
	assert.Equal(t, "b", val)
	assert.True(t, ok)
	assert.Equal(t, 2, cc.Count())
	val, ok = At(2, cc.All())
	assert.Equal(t, "c", val)
	assert.True(t, ok)
	assert.Equal(t, 3, cc.Count())
	val, ok = At(3, cc.All())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 3, cc.Count())
	val, ok = At(4, cc.All())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 3, cc.Count())
}

func firstOf[T any](seq iter.Seq[T]) T {
	result, _ := First(seq)
	return result
}

func firstNOf[T any](n int, seq iter.Seq[T]) []T {
	if n <= 0 {
		return slices.Collect(seq)
	}
	return slices.Collect(Take(n, seq))
}

// values returns an iter.Seq over s that also sets s elements to zero until
// yield return false.
func values(s []int) iter.Seq[int] {
	return func(yield func(x int) bool) {
		for i := range s {
			if !yield(s[i]) {
				return
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

type countConsumed[T any] struct {
	delegate iter.Seq[T]
	count    int
}

func newCountConsumed[T any](delegate iter.Seq[T]) *countConsumed[T] {
	return &countConsumed[T]{delegate: delegate}
}

func (c *countConsumed[T]) Count() int { return c.count }

func (c *countConsumed[T]) All() iter.Seq[T] {
	c.count = 0
	return func(yield func(T) bool) {
		c.count = 0
		for e := range c.delegate {
			c.count++
			if !yield(e) {
				return
			}
		}
	}
}
