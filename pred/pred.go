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

func AttrCond(a string, p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasValCond(a, p)
	}
}

func AttrWord(a, word string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasWord(a, word)
	}
}

func AttrWordCond(a string, p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasWordCond(a, p)
	}
}

func ID(id string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasID(id)
	}
}

func IDCond(p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasIDCond(p)
	}
}

func Class(class string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasClass(class)
	}
}

func ClassCond(p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && attr.L(h.Attr).HasClassCond(p)
	}
}

func InnerText(s string) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.FirstChild != nil &&
			h.FirstChild.Type == html.TextNode &&
			h.FirstChild.Data == s
	}
}

func InnerTextCond(p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.FirstChild != nil &&
			h.FirstChild.Type == html.TextNode &&
			p(h.FirstChild.Data)
	}
}
