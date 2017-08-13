package find

import (
	"bytes"

	"golang.org/x/net/html"
)

func f2s(f Finder) string {
	return h2s(f.Node)
}

func h2s(h *html.Node) string {
	if h == nil {
		return ""
	}
	var b bytes.Buffer
	err := html.Render(&b, h)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func s2f(s string) Finder {
	h, err := html.Parse(bytes.NewBufferString(s))
	if err != nil {
		panic(err)
	}
	return NewFinder(h)
}
