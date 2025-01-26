// Package itertools contains iterator functions like Zip.
package itertools

import (
	"iter"
	"slices"
	"sync"
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
		z := &zipper[F, S]{
			fc:   make(chan F),
			sc:   make(chan S),
			fadv: make(chan bool),
			sadv: make(chan bool),
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			first(z.firstYield)
			close(z.fc)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			second(z.secondYield)
			close(z.sc)
			wg.Done()
		}()
		z.iterate(yield)
		wg.Wait()
	}
}

// Cycle returns values repeating in an infinite cycle.
func Cycle[T any](values ...T) iter.Seq[T] {
	if len(values) == 0 {
		panic("values must be non-empty")
	}
	valueCopy := slices.Clone(values)
	return func(yield func(T) bool) {
		for {
			for _, v := range valueCopy {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Flatten flattens the passed in iter.Seq[T] into a single iter.Seq[T]
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

type zipper[F, S any] struct {
	fc   chan F
	sc   chan S
	fadv chan bool
	sadv chan bool
}

func (z *zipper[F, S]) firstYield(first F) bool {
	select {
	case <-z.fadv:
		return false
	default:
		z.fc <- first
		return <-z.fadv
	}
}

func (z *zipper[F, S]) secondYield(second S) bool {
	select {
	case <-z.sadv:
		return false
	default:
		z.sc <- second
		return <-z.sadv
	}
}

func (z *zipper[F, S]) iterate(yield func(F, S) bool) {
	for {
		first, fok := <-z.fc
		second, sok := <-z.sc
		if !fok || !sok {
			break
		}
		if !yield(first, second) {
			break
		}
		z.fadv <- true
		z.sadv <- true
	}
	close(z.fadv)
	close(z.sadv)
}
