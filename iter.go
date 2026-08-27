// Package itertools contains iterator functions like Zip.
package itertools

import (
	"iter"
	"slices"
)

// Pair represents a pair of values.
type Pair[F, S any] struct {
	First  F
	Second S
}

// Filterer filters T values.
type Filterer[T any] func(T) bool

// Mapper maps T values to U values.
type Mapper[T, U any] func(T) U

// PairUp converts an iter.Seq2 into an iter.Seq of pairs.
func PairUp[F, S any](seq iter.Seq2[F, S]) iter.Seq[Pair[F, S]] {
	return func(yield func(Pair[F, S]) bool) {
		for first, second := range seq {
			p := Pair[F, S]{First: first, Second: second}
			if !yield(p) {
				return
			}
		}
	}
}

// Zip works like zip in python.
func Zip[F, S any](first iter.Seq[F], second iter.Seq[S]) iter.Seq2[F, S] {
	return func(yield func(F, S) bool) {
		firstPull, firstStop := iter.Pull(first)
		defer firstStop()
		secondPull, secondStop := iter.Pull(second)
		defer secondStop()
		firstValue, firstOk := firstPull()
		secondValue, secondOk := secondPull()
		for firstOk && secondOk {
			if !yield(firstValue, secondValue) {
				return
			}
			firstValue, firstOk = firstPull()
			secondValue, secondOk = secondPull()
		}
	}
}

// CycleValues returns values repeating in an infinite cycle.
func CycleValues[T any](values ...T) iter.Seq[T] {
	if len(values) == 0 {
		return empty[T]
	}
	valuesCopy := slices.Clone(values)
	return func(yield func(T) bool) {
		yieldCycleValues(valuesCopy, yield)
	}
}

// Cycle returns the values in seq repeating in an infinite cycle.
func Cycle[T any](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var saved []T
		for value := range seq {
			if !yield(value) {
				return
			}
			saved = append(saved, value)
		}
		if len(saved) > 0 {
			yieldCycleValues(saved, yield)
		}
	}
}

// Chain returns all the elements in the first sequence followed by all the
// elements in the second sequence etc.
func Chain[T any](sequences ...iter.Seq[T]) iter.Seq[T] {
	if len(sequences) == 0 {
		return empty[T]
	}
	if len(sequences) == 1 {
		return sequences[0]
	}
	sequencesCopy := slices.Clone(sequences)
	return func(yield func(T) bool) {
		for _, sequence := range sequencesCopy {
			for v := range sequence {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Filter returns an iter.Seq[T] that contains all the T values in seq for
// which f returns true.
func Filter[T any](f Filterer[T], seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x := range seq {
			if f(x) && !yield(x) {
				return
			}
		}
	}
}

// Map returns an iter.Seq[U] which is m applied to each element in seq.
func Map[T, U any](m Mapper[T, U], seq iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for x := range seq {
			if !yield(m(x)) {
				return
			}
		}
	}
}

// Count returns start, start + step, start + 2*step, ...
func Count(start, step int) iter.Seq[int] {
	if start == 0 && step == 1 {
		return simpleCount
	}
	return func(yield func(int) bool) {
		for i := start; ; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

// Enumerate works like enumerate in python. It is equivalent to
// Zip(Count(0, 1), seq)
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return Zip(Count(0, 1), seq)
}

// At returns the 0-based indexth element of seq. At returns false if seq
// has index or fewer elements or if index is negative.
func At[T any](index int, seq iter.Seq[T]) (element T, ok bool) {
	if index < 0 {
		return
	}
	return First(Drop(index, seq))
}

// MustAt works like At except that it panics if index is out of range
// instead of returning false.
func MustAt[T any](index int, seq iter.Seq[T]) T {
	return must(At(index, seq))
}

// First returns the first element of seq or false if seq is empty.
func First[T any](seq iter.Seq[T]) (first T, ok bool) {
	for x := range seq {
		first = x
		ok = true
		break
	}
	return
}

// MustFirst works like First except that it panics if seq is empty
// instead of returning false.
func MustFirst[T any](seq iter.Seq[T]) T {
	return must(First(seq))
}

// Drop returns everything but the first n elements of seq.
func Drop[T any](n int, seq iter.Seq[T]) iter.Seq[T] {
	if n <= 0 {
		return seq
	}
	return func(yield func(T) bool) {
		count := 0
		for x := range seq {
			if count >= n && !yield(x) {
				return
			}
			count++
		}
	}
}

// DropWhile returns everything but the first elements of seq for which f
// returns true.
func DropWhile[T any](f Filterer[T], seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		emit := false
		for x := range seq {
			if !emit && !f(x) {
				emit = true
			}
			if emit && !yield(x) {
				return
			}
		}
	}
}

// Take returns the first n elements of seq.
func Take[T any](n int, seq iter.Seq[T]) iter.Seq[T] {
	if n <= 0 {
		return empty[T]
	}
	return func(yield func(T) bool) {
		count := 0
		for x := range seq {
			count++
			if !yield(x) || count == n {
				return
			}
		}
	}
}

// TakeWhile returns the first elements of seq for which f returns true.
func TakeWhile[T any](f Filterer[T], seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x := range seq {
			if !f(x) || !yield(x) {
				return
			}
		}
	}
}

// Capture wraps an iter.Seq in order to capture statistics on it.
type Capture[V any] struct {
	seq        iter.Seq[V]
	firstValue V
	lastValue  V
	count      int
}

// NewCapture returns a Capture that wraps seq.
func NewCapture[V any](seq iter.Seq[V]) *Capture[V] {
	return &Capture[V]{seq: seq}
}

// Values resets the statistics and returns the values of the wrapped
// iterator.
func (c *Capture[V]) Values() iter.Seq[V] {
	c.reset()
	return func(yield func(value V) bool) {
		c.reset()
		for v := range c.seq {
			c.update(v)
			if !yield(v) {
				return
			}
		}
	}
}

// Count returns the number of values taken from the wrapped iterator.
func (c *Capture[V]) Count() int {
	return c.count
}

// FirstValue returns the first value taken from wrapped iterator.
// FirstValue panics if Count() returns 0.
func (c *Capture[V]) FirstValue() V {
	c.checkNotEmpty()
	return c.firstValue
}

// LastValue returns the last value taken from wrapped iterator.
// LastValue panics if Count() returns 0.
func (c *Capture[V]) LastValue() V {
	c.checkNotEmpty()
	return c.lastValue
}

func (c *Capture[V]) checkNotEmpty() {
	if c.count == 0 {
		panic("Capture empty")
	}
}

func (c *Capture[V]) reset() {
	c.count = 0
}

func (c *Capture[V]) update(v V) {
	if c.count == 0 {
		c.firstValue = v
	}
	c.lastValue = v
	c.count++
}

// Capture2 wraps an iter.Seq2 in order to capture statistics on it.
type Capture2[K, V any] struct {
	seq        iter.Seq2[K, V]
	firstValue V
	lastValue  V
	firstKey   K
	lastKey    K
	count      int
}

// NewCapture2 returns a Capture2 that wraps seq.
func NewCapture2[K, V any](seq iter.Seq2[K, V]) *Capture2[K, V] {
	return &Capture2[K, V]{seq: seq}
}

// Keys resets the statistics and returns just the keys of the wrapped
// iterator.
func (c *Capture2[K, V]) Keys() iter.Seq[K] {
	c.reset()
	return func(yield func(key K) bool) {
		c.reset()
		for k, v := range c.seq {
			c.update(k, v)
			if !yield(k) {
				return
			}
		}
	}
}

// Values resets the statistics and returns just the values of the wrapped
// iterator.
func (c *Capture2[K, V]) Values() iter.Seq[V] {
	c.reset()
	return func(yield func(value V) bool) {
		c.reset()
		for k, v := range c.seq {
			c.update(k, v)
			if !yield(v) {
				return
			}
		}
	}
}

// Count returns the number of key-value pairs taken from the wrapped
// iterator.
func (c *Capture2[K, V]) Count() int {
	return c.count
}

// FirstKey returns the first key taken from wrapped iterator.
// FirstKey panics if Count() returns 0.
func (c *Capture2[K, V]) FirstKey() K {
	c.checkNotEmpty()
	return c.firstKey
}

// LastKey returns the last key taken from wrapped iterator.
// LastKey panics if Count() returns 0.
func (c *Capture2[K, V]) LastKey() K {
	c.checkNotEmpty()
	return c.lastKey
}

// FirstValue returns the first value taken from wrapped iterator.
// FirstValue panics if Count() returns 0.
func (c *Capture2[K, V]) FirstValue() V {
	c.checkNotEmpty()
	return c.firstValue
}

// LastValue returns the last value taken from wrapped iterator.
// LastValue panics if Count() returns 0.
func (c *Capture2[K, V]) LastValue() V {
	c.checkNotEmpty()
	return c.lastValue
}

func (c *Capture2[K, V]) checkNotEmpty() {
	if c.count == 0 {
		panic("Capture empty")
	}
}

func (c *Capture2[K, V]) reset() {
	c.count = 0
}

func (c *Capture2[K, V]) update(k K, v V) {
	if c.count == 0 {
		c.firstKey = k
		c.firstValue = v
	}
	c.lastKey = k
	c.lastValue = v
	c.count++
}

func simpleCount(yield func(int) bool) {
	for i := 0; ; i++ {
		if !yield(i) {
			return
		}
	}
}

func empty[T any](yield func(T) bool) {
}

func yieldCycleValues[T any](values []T, yield func(T) bool) {
	for {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

func must[T any](result T, ok bool) T {
	if !ok {
		panic("No such element")
	}
	return result
}
