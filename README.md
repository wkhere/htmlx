## htmlx

[![Build Status](https://github.com/wkhere/htmlx/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/wkhere/htmlx/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/wkhere/htmlx/badge.svg?branch=master&kill_cache=1)](https://coveralls.io/github/wkhere/htmlx?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/wkhere/htmlx.svg)](https://pkg.go.dev/github.com/wkhere/htmlx)

HTML parsing eXtensions, aka monadic Finder of elements.

The idea is to wrap parsed `html.Node` with `Finder` and chain the calls for finding HTML elements within subtrees.

Includes handy `parsehtml` command to visualize given HTML file or the URL content.

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

[...]

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
```
