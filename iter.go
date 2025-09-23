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

// First returns the first element of seq or false if seq is empty.
func First[T any](seq iter.Seq[T]) (first T, ok bool) {
	for x := range seq {
		first = x
		ok = true
		break
	}
	return
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
			if !f(x) {
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
			if count == n || !yield(x) {
				return
			}
			count++
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

func simpleCount(yield func(int) bool) {
	for i := 0; ; i++ {
		if !yield(i) {
			return
		}
	}
}

func empty[T any](yield func(T) bool) {
	return
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
