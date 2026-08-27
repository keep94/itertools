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
	assert.Equal(t, 4, MustFirst(it))
}

func TestFilter(t *testing.T) {
	x := []int{3, 4, 5, 6}
	f := func(i int) bool { return i%2 == 1 }
	it := Filter(f, slices.Values(x))
	assert.Equal(t, []int{3, 5}, slices.Collect(it))
	assert.Equal(t, 3, MustFirst(it))
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
	capture := NewCapture(Count(10, 1))
	assert.Empty(t, slices.Collect(Take(-1, capture.Values())))
	assert.Equal(t, 0, capture.Count())
	assert.Empty(t, slices.Collect(Take(0, capture.Values())))
	assert.Equal(t, 0, capture.Count())
	assert.Equal(t, []int{10}, slices.Collect(Take(1, capture.Values())))
	assert.Equal(t, 1, capture.Count())
	takeSeq := Take(3, capture.Values())
	assert.Equal(t, []int{10, 11, 12}, slices.Collect(takeSeq))
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, 10, MustFirst(takeSeq))
	assert.Equal(t, 1, capture.Count())
}

func TestTakeFinite(t *testing.T) {
	capture := NewCapture(slices.Values([]string{"abc", "123", "foo"}))
	takeSeq := Take(4, capture.Values())
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, "abc", MustFirst(takeSeq))
	assert.Equal(t, 1, capture.Count())
}

func TestTakeWhile(t *testing.T) {
	capture := NewCapture(
		slices.Values([]int{10, 11, 12, 13, 14, 15, 1, 2, 3, 4}))
	f := func(x int) bool { return x < 15 }
	g := func(x int) bool { return x < 10 }
	assert.Empty(t, slices.Collect(TakeWhile(g, capture.Values())))
	assert.Equal(t, 1, capture.Count())
	takeSeq := TakeWhile(f, capture.Values())
	assert.Equal(t, []int{10, 11, 12, 13, 14}, slices.Collect(takeSeq))
	assert.Equal(t, 6, capture.Count())
	assert.Equal(t, 10, MustFirst(takeSeq))
	assert.Equal(t, 1, capture.Count())
}

func TestTakeWhileFinite(t *testing.T) {
	capture := NewCapture(slices.Values([]string{"abc", "123", "foo"}))
	f := func(s string) bool { return len(s) < 4 }
	takeSeq := TakeWhile(f, capture.Values())
	assert.Equal(t, []string{"abc", "123", "foo"}, slices.Collect(takeSeq))
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, "abc", MustFirst(takeSeq))
	assert.Equal(t, 1, capture.Count())
}

func TestDropWhile(t *testing.T) {
	seq := slices.Values([]int{10, 13, 16, 1, 2, 3})
	f := func(x int) bool { return x < 10 }
	g := func(x int) bool { return x < 13 }
	h := func(x int) bool { return x < 16 }
	dropSeq := DropWhile(f, seq)
	assert.Equal(t, []int{10, 13, 16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 10, MustFirst(dropSeq))
	dropSeq = DropWhile(g, seq)
	assert.Equal(t, []int{13, 16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 13, MustFirst(dropSeq))
	dropSeq = DropWhile(h, seq)
	assert.Equal(t, []int{16, 1, 2, 3}, slices.Collect(dropSeq))
	assert.Equal(t, 16, MustFirst(dropSeq))
}

func TestDropWhileCallFilterOnlyAtBeginning(t *testing.T) {
	filterCallCount := 0
	filter := func(x int) bool {
		filterCallCount++
		return x == 1
	}
	seq := slices.Values([]int{1, 2, 1, 1, 1})
	dropSeq := DropWhile(filter, seq)
	assert.Equal(t, []int{2, 1, 1, 1}, slices.Collect(dropSeq))
	assert.Equal(t, 2, filterCallCount)
}

func TestDropWhileFinite(t *testing.T) {
	seq := slices.Values([]int{10, 13, 16})
	f := func(x int) bool { return x < 19 }
	dropSeq := DropWhile(f, seq)
	assert.Empty(t, slices.Collect(dropSeq))
	assert.Empty(t, slices.Collect(dropSeq))
}

func TestAt(t *testing.T) {
	capture := NewCapture(slices.Values([]string{"a", "b", "c"}))
	val, ok := At(-1, capture.Values())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 0, capture.Count())
	val, ok = At(0, capture.Values())
	assert.Equal(t, "a", val)
	assert.True(t, ok)
	assert.Equal(t, 1, capture.Count())
	val, ok = At(1, capture.Values())
	assert.Equal(t, "b", val)
	assert.True(t, ok)
	assert.Equal(t, 2, capture.Count())
	val, ok = At(2, capture.Values())
	assert.Equal(t, "c", val)
	assert.True(t, ok)
	assert.Equal(t, 3, capture.Count())
	val, ok = At(3, capture.Values())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 3, capture.Count())
	val, ok = At(4, capture.Values())
	assert.Equal(t, "", val)
	assert.False(t, ok)
	assert.Equal(t, 3, capture.Count())
}

func TestMustAt(t *testing.T) {
	seq := slices.Values([]string{"a", "b", "c"})
	assert.Equal(t, "a", MustAt(0, seq))
	assert.Equal(t, "b", MustAt(1, seq))
	assert.Equal(t, "c", MustAt(2, seq))
	assert.Panics(t, func() { MustAt(-1, seq) })
	assert.Panics(t, func() { MustAt(3, seq) })
}

func TestFirst(t *testing.T) {
	seq := slices.Values([]string{"a"})
	val, ok := First(seq)
	assert.Equal(t, "a", val)
	assert.True(t, ok)
	emptySeq := Chain[string]()
	_, ok = First(emptySeq)
	assert.False(t, ok)
}

func TestMustFirst(t *testing.T) {
	seq := slices.Values([]string{"a"})
	assert.Equal(t, "a", MustFirst(seq))
	emptySeq := Chain[string]()
	assert.Panics(t, func() { MustFirst(emptySeq) })
}

func TestCapture(t *testing.T) {
	seq := slices.Values([]string{"a", "b", "c"})
	capture := NewCapture(seq)
	result := slices.Collect(capture.Values())
	assert.Equal(t, []string{"a", "b", "c"}, result)
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, "c", capture.LastValue())

	emptyResult := slices.Collect(Take(0, capture.Values()))
	assert.Empty(t, emptyResult)
	assert.Zero(t, capture.Count())
	assert.Panics(t, func() { capture.FirstValue() })
	assert.Panics(t, func() { capture.LastValue() })

	singleResult := slices.Collect(Take(1, capture.Values()))
	assert.Equal(t, []string{"a"}, singleResult)
	assert.Equal(t, 1, capture.Count())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, "a", capture.LastValue())
}

func TestCapture2Values(t *testing.T) {
	seq := slices.All([]string{"a", "b", "c"})
	capture := NewCapture2(seq)
	result := slices.Collect(capture.Values())
	assert.Equal(t, []string{"a", "b", "c"}, result)
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, 0, capture.FirstKey())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, 2, capture.LastKey())
	assert.Equal(t, "c", capture.LastValue())

	emptyResult := slices.Collect(Take(0, capture.Values()))
	assert.Empty(t, emptyResult)
	assert.Zero(t, capture.Count())
	assert.Panics(t, func() { capture.FirstKey() })
	assert.Panics(t, func() { capture.FirstValue() })
	assert.Panics(t, func() { capture.LastKey() })
	assert.Panics(t, func() { capture.LastValue() })

	singleResult := slices.Collect(Take(1, capture.Values()))
	assert.Equal(t, []string{"a"}, singleResult)
	assert.Equal(t, 1, capture.Count())
	assert.Equal(t, 0, capture.FirstKey())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, 0, capture.LastKey())
	assert.Equal(t, "a", capture.LastValue())
}

func TestCapture2Keys(t *testing.T) {
	seq := slices.All([]string{"a", "b", "c"})
	capture := NewCapture2(seq)
	result := slices.Collect(capture.Keys())
	assert.Equal(t, []int{0, 1, 2}, result)
	assert.Equal(t, 3, capture.Count())
	assert.Equal(t, 0, capture.FirstKey())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, 2, capture.LastKey())
	assert.Equal(t, "c", capture.LastValue())

	emptyResult := slices.Collect(Take(0, capture.Keys()))
	assert.Empty(t, emptyResult)
	assert.Zero(t, capture.Count())
	assert.Panics(t, func() { capture.FirstKey() })
	assert.Panics(t, func() { capture.FirstValue() })
	assert.Panics(t, func() { capture.LastKey() })
	assert.Panics(t, func() { capture.LastValue() })

	singleResult := slices.Collect(Take(1, capture.Keys()))
	assert.Equal(t, []int{0}, singleResult)
	assert.Equal(t, 1, capture.Count())
	assert.Equal(t, 0, capture.FirstKey())
	assert.Equal(t, "a", capture.FirstValue())
	assert.Equal(t, 0, capture.LastKey())
	assert.Equal(t, "a", capture.LastValue())
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
