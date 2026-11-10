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
	"github.com/wkhere/htmlx"
	"github.com/wkhere/htmlx/pred"
	"golang.org/x/net/html/atom"
)

// ...
top, err := htmlx.FinderFromData(resp.Body)
// err handling...

myElement := top.
	Find(pred.Element(atom.Body)).
	Find(pred.Element(atom.Div, pred.ID("content"))).
	Find(pred.Element(atom.Table, pred.Class("foo"))).
	Find(pred.Element(atom.Tr)).FindSiblings(pred.Element(atom.Tr))
```
