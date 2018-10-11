package htmlx

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Finder struct {
	*html.Node
}

func FinderFromNode(h *html.Node) Finder {
	return Finder{h}
}

func FinderFromData(r io.Reader) (f Finder, err error) {
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

func (f Finder) AttrVal(key string) (val string, ok bool) {
	if f.Node == nil {
		return val, false
	}
	return AttrVal(f.Node.Attr, key)
}

func (f Finder) Id() (id string, ok bool) {
	return f.AttrVal("id")
}

func (f Finder) ClassList() (classes []string, ok bool) {
	classStr, ok := f.AttrVal("class")
	if !ok {
		return classes, false
	}
	return strings.Fields(classStr), ok
}

func (f Finder) HasAttr(key string) bool {
	if f.Node == nil {
		return false
	}
	return HasAttr(f.Node.Attr, key)
}

func (f Finder) HasAttrVal(key, val string) bool {
	if f.Node == nil {
		return false
	}
	return HasAttrVal(f.Node.Attr, key, val)
}

func (f Finder) HasAttrWord(key, word string) bool {
	if f.Node == nil {
		return false
	}
	return HasAttrWord(f.Node.Attr, key, word)
}

func (f Finder) HasId(id string) bool {
	return f.HasAttrVal("id", id)
}

func (f Finder) HasClass(class string) bool {
	return f.HasAttrWord("class", class)
}

type FinderPredicate func(*html.Node) bool

// Find is universal finder.
// Includes current node in search.
func (f Finder) Find(pred FinderPredicate) (r Finder) {
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

// FindSibling performs flat find among node's siblings.
// No recursion. Omits current node, starts from a first sibling.
func (f Finder) FindSibling(pred FinderPredicate) (r Finder) {
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

func elementP(element atom.Atom) FinderPredicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.DataAtom == element
	}
}

func attrP(attr, val string) FinderPredicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && HasAttrVal(h.Attr, attr, val)
	}
}

func attrWordP(attr, word string) FinderPredicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && HasAttrWord(h.Attr, attr, word)
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

func (f Finder) FindByClass(class string) Finder {
	return f.Find(attrWordP("class", class))
}

func (f Finder) FindSiblingById(id string) Finder {
	return f.FindSibling(attrP("id", id))
}

func (f Finder) FindSiblingByClass(class string) Finder {
	return f.FindSibling(attrWordP("class", class))
}

func AttrVal(attr []html.Attribute, key string) (val string, ok bool) {
	for _, a := range attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return
}

func HasAttr(attr []html.Attribute, key string) (ok bool) {
	for _, a := range attr {
		if a.Key == key {
			return true
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

func HasAttrWord(attr []html.Attribute, key, word string) bool {
	foundVal, ok := AttrVal(attr, key)
	if !ok {
		return false
	}
	for _, w := range strings.Fields(foundVal) {
		if w == word {
			return true
		}
	}
	return false
}
