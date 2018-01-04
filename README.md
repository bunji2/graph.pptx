# bunji2/graph.pptx

Go package that extracts the graph data of nodes and edges from directed graph edited in PPTX.

## Installation

```
go get github.com/bunji2/graph.pptx
```

## Usage

Followin sample extracts nodes and edges from PPTX file "sample.pptx".

```go
import (
    "fmt"
    "github.com/bunji2/graph.pptx"
)

func sample() {
    gp.Init()
    if err := gp.Parse("sample.pptx"); err == nil {
        gp.Dump()
    }
}
```

## Document
 * [godoc.org](https://godoc.org/github.com/bunji2/graph.pptx)

## License

under the MIT License

by Bunji2