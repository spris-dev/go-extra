# go-extra

A collection of commonly used Go utilities.

## Installation

```sh
go get go.spris.dev
```

## Usage

```go
import (
  "go.spris.dev/x"
)

func main() {
  result := x.Iter([]int{1, 2, 3, 4}).
    Map(func(v int) int { return v * 2 }).
    Fold(0, func(sum, v int) int { return sum + v })

  fmt.Println(result)
}
```
