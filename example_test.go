package itertools_test

import (
	"fmt"
	"slices"

	"github.com/keep94/itertools"
)

func ExampleEnumerate() {
	notesIter := slices.Values([]string{"do", "re", "mi", "fa", "so"})
	for i, n := range itertools.Enumerate(notesIter) {
		fmt.Println(i, n)
	}
	// Output:
	// 0 do
	// 1 re
	// 2 mi
	// 3 fa
	// 4 so
}

func ExampleZip() {
	notesIter := slices.Values([]string{"do", "re", "mi", "fa", "so"})
	ordinalsIter := slices.Values([]int{1, 2, 3})
	for n, o := range itertools.Zip(notesIter, ordinalsIter) {
		fmt.Println(n, o)
	}
	// Output:
	// do 1
	// re 2
	// mi 3
}

func ExampleChain() {
	notes := []string{"do", "re", "mi", "fa", "so"}
	ordinals := []int{1, 2, 3}
	notesIter := slices.Values(notes)
	ordinalsIter := itertools.Chain(
		slices.Values(ordinals), itertools.CycleValues(0))
	for n, o := range itertools.Zip(notesIter, ordinalsIter) {
		fmt.Println(n, o)
	}
	// Output:
	// do 1
	// re 2
	// mi 3
	// fa 0
	// so 0
}

func ExampleTakeWhile() {
	seq := itertools.CycleValues(1, 2, 3, 4, 5)
	f := func(x int) bool { return x < 4 }
	for x := range itertools.TakeWhile(f, seq) {
		fmt.Println(x)
	}
	// Output:
	// 1
	// 2
	// 3
}
