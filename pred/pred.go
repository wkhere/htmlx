package pred

import (
	"github.com/wkhere/htmlx/attr"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Predicate func(*html.Node) bool

func True() Predicate {
	return func(*html.Node) bool { return true }
}

func AnyElement() Predicate {
	return func(h *html.Node) bool { return h.Type == html.ElementNode }
}

func Element(element atom.Atom, pp ...Predicate) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.ElementNode && h.DataAtom == element &&
			all(pp, h)
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

func AnyText() Predicate {
	return func(h *html.Node) bool { return h.Type == html.TextNode }
}

func Text(s string, pp ...Predicate) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.TextNode && h.Data == s && all(pp, h)
	}
}

func TextCond(p func(string) bool) Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.TextNode && p(h.Data)
	}
}

func IsText() Predicate {
	return func(h *html.Node) bool {
		return h.Type == html.TextNode
	}
}

func all(pp []Predicate, h *html.Node) bool {
	for _, p := range pp {
		if !p(h) {
			return false
		}
	}
	return true
}
