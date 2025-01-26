itertools
=========

Package itertools contains iterator functions like Zip.

## Examples

```golang
package main

import (
    "fmt"
    "slices"

    "github.com/keep94/itertools"
)

func main() {
    lettersIter := slices.Values([]string{"a", "b", "c"})
    numbersIter := slices.Values([]int{1, 2, 3})

    // Prints:
    // a 1
    // b 2
    // c 3
    for letter, number := range itertools.Zip(lettersIter, numbersIter) {
        fmt.Println(letter, number)
    }
}
```

More documentation and examples can be found [here](https://pkg.go.dev/github.com/keep94/itertools).
