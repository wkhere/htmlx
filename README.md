## htmlx

[![Build Status](https://github.com/wkhere/htmlx/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/wkhere/htmlx/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/wkhere/htmlx/badge.svg?branch=master)](https://coveralls.io/github/wkhere/htmlx?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/wkhere/htmlx.svg)](https://pkg.go.dev/github.com/wkhere/htmlx)

HTML parsing eXtensions.

`Finder` wraps a parsed `html.Node` and locates elements through chained `Find*` calls, each taking a predicate built with the `pred` package.

The chaining works much like Haskell's `Maybe` monad. Each `Find*` call searches relative to the current element (its subtree or its siblings) and wraps the match in a new `Finder`, so the next call in the chain continues from there. When nothing matches, the call returns an empty `Finder`, which plays the role of `Nothing`: it is still a valid receiver, and any further calls on it simply return empty results. This way a whole chain can be written without checking for a miss after each step; the result is checked once, at the end. Methods yielding multiple elements, like `FindAll` or `FindSiblings`, return a `FinderStream`, which acts like the list monad, with `Join` serving as its bind.

The project also includes a handy `parsehtml` command to visualize a parsed HTML file or URL content.

### Example

```Go
import (
        "fmt"
        "net/http"
        "strings"

        "github.com/wkhere/htmlx"
        "github.com/wkhere/htmlx/pred"
        "golang.org/x/net/html/atom"
)

func example() error {
        resp, err := http.Get("https://go.dev/doc/devel/release")
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        top, err := htmlx.FinderFromData(resp.Body)
        if err != nil {
                return err
        }

        elements := top.
                Find(pred.Element(atom.Main)).
                Find(pred.ID("go1.27.minor")).
                FindSiblings(pred.Element(atom.P, pred.IDCond(func(id string) bool {
                        return strings.HasPrefix(id, "go1.27.")
                })))

        for _, f := range elements.Collect() {
                fmt.Println(f.String())
        }
        return nil
}
```
