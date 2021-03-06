package htmlx

import (
	"bytes"
	"io"

	"golang.org/x/net/html"

	"github.com/wkhere/htmlx/attr"
	"github.com/wkhere/htmlx/pred"
)

type Finder struct {
	*html.Node
}

func FinderFromNode(h *html.Node) Finder {
	return Finder{h}
}

func FinderFromData(r io.Reader) (Finder, error) {
	h, err := html.Parse(r)
	return Finder{h}, err
}

func FinderFromString(s string) (Finder, error) {
	return FinderFromData(bytes.NewBufferString(s))
}

func (f Finder) IsEmpty() bool {
	return f.Node == nil
}

func (f Finder) String() string {
	if f.Node == nil {
		return ""
	}
	var b bytes.Buffer
	html.Render(&b, f.Node)
	return b.String()
}

func (f Finder) Parent() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.Parent}
}

func (f Finder) FirstChild() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.FirstChild}
}

func (f Finder) LastChild() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.LastChild}
}

func (f Finder) PrevSibling() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.PrevSibling}
}

func (f Finder) NextSibling() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.NextSibling}
}

func (f Finder) Attr() attr.List {
	if f.Node == nil {
		return nil
	}
	return f.Node.Attr
}

// Find performs universal depth-first find.
// Includes current node in the search.
func (f Finder) Find(pred pred.Predicate) (r Finder) {
	if f.Node == nil {
		return
	}

	var walker func(node *html.Node) bool

	walker = func(node *html.Node) bool {
		if pred(node) {
			r.Node = node
			return true
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if walker(c) {
				return true
			}
		}
		return false
	}

	walker(f.Node)
	return
}

// FindSibling performs flat find among node's next (right) siblings.
// No recursion. Omits current node, starts from a first such sibling.
func (f Finder) FindSibling(pred pred.Predicate) (r Finder) {
	if f.Node == nil {
		return
	}

	for c := f.Node.NextSibling; c != nil; c = c.NextSibling {
		if pred(c) {
			return Finder{c}
		}
	}
	return
}

// FindPrevSibling performs flat find among node's previous (left) siblings.
// No recursion. Omits current node, starts from a first such sibling.
func (f Finder) FindPrevSibling(pred pred.Predicate) (r Finder) {
	if f.Node == nil {
		return
	}

	for c := f.Node.PrevSibling; c != nil; c = c.PrevSibling {
		if pred(c) {
			return Finder{c}
		}
	}
	return
}
