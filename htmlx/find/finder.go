package find

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Finder struct {
	*html.Node
}

func NewFinder(h *html.Node) Finder {
	return Finder{h}
}

func (f Finder) IsEmpty() bool {
	return f.Node == nil
}

func (f Finder) FirstChild() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.FirstChild}
}

func (f Finder) NextSibling() Finder {
	if f.Node == nil {
		return f
	}
	return Finder{f.Node.NextSibling}
}

type FoundPredicate func(*html.Node) bool
type walkerFunc func(*html.Node) bool

// Find is universal finder.
// Includes current node in search.
// Note that adding closures raises execution time
// from 210 ns/op to 231 ns/op. Can live with it.
// (Measurements done on `metal` machine, perf mode.)
// Btw measures above are obsolete, they predate use
// of Finder struct and methods.
func (f Finder) Find(pred FoundPredicate) (r Finder) {
	if f.Node == nil {
		return
	}
	var walker walkerFunc

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

// FindSibling performs flat find among node's siblings.
// No recursion. Omits current node, starts from a first sibling.
func (f Finder) FindSibling(pred FoundPredicate) (r Finder) {
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

func elementP(element atom.Atom) FoundPredicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.DataAtom == element
	}
}

func attrP(attr, val string) FoundPredicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && HasAttrVal(h.Attr, attr, val)
	}
}

func (f Finder) FindElement(element atom.Atom) Finder {
	return f.Find(elementP(element))
}

func (f Finder) FindSiblingElement(element atom.Atom) Finder {
	return f.FindSibling(elementP(element))
}

func (f Finder) FindByAttr(attr, val string) Finder {
	return f.Find(attrP(attr, val))
}

func (f Finder) FindSiblingByAttr(attr, val string) Finder {
	return f.FindSibling(attrP(attr, val))
}

func (f Finder) FindById(id string) Finder {
	return f.Find(attrP("id", id))
}

func (f Finder) FindSiblingById(id string) Finder {
	return f.FindSibling(attrP("id", id))
}

func AttrVal(attr []html.Attribute, key string) (val string, ok bool) {
	for _, a := range attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return
}

func HasAttrVal(attr []html.Attribute, key, val string) bool {
	foundVal, ok := AttrVal(attr, key)
	if !ok {
		return false
	}
	return val == foundVal
}
