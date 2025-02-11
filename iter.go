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
func PairUp[F, S any](s iter.Seq2[F, S]) iter.Seq[Pair[F, S]] {
	return func(yield func(Pair[F, S]) bool) {
		for first, second := range s {
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

// Cycle returns values repeating in an infinite cycle.
func Cycle[T any](values ...T) iter.Seq[T] {
	valueCopy := slices.Clone(values)
	return func(yield func(T) bool) {
		if len(valueCopy) == 0 {
			return
		}
		for {
			for _, v := range valueCopy {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Flatten flattens the passed in iterators into a single iterator.
func Flatten[T any](iterators ...iter.Seq[T]) iter.Seq[T] {
	iteratorCopy := slices.Clone(iterators)
	return func(yield func(T) bool) {
		for _, iterator := range iteratorCopy {
			for v := range iterator {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Filter returns an iter.Seq[T] that contains all the T values in seq for
// which f returns true.
func Filter[T any](seq iter.Seq[T], f Filterer[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x := range seq {
			if !f(x) {
				continue
			}
			if !yield(x) {
				return
			}
		}
	}
}

// Map returns an iter.Seq[U] which is m applied to each element in seq.
func Map[T, U any](seq iter.Seq[T], m Mapper[T, U]) iter.Seq[U] {
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

func simpleCount(yield func(int) bool) {
	for i := 0; ; i++ {
		if !yield(i) {
			return
		}
	}
}
