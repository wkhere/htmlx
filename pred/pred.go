package pred

import (
	"github.com/wkhere/htmlx/attr"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Predicate func(*html.Node) bool

func Element(element atom.Atom) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.DataAtom == element
	}
}

func Attr(a, val string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasVal(a, val)
	}
}

func AttrWord(a, word string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasWord(a, word)
	}
}

func ID(id string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasID(id)
	}
}

func Class(class string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasClass(class)
	}
}
