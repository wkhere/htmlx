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

type FinderStream <-chan Finder

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

func (f Finder) Write(w io.Writer) {
	html.Render(w, f.Node)
}

func (f Finder) String() string {
	if f.Node == nil {
		return ""
	}
	var b bytes.Buffer
	f.Write(&b)
	return b.String()
}

func (ff FinderStream) Collect() (res []Finder) {
	for f := range ff {
		res = append(res, f)
	}
	return
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

func (f Finder) InnerText() string {
	f1 := f.Find(pred.IsText())
	if f1.IsEmpty() {
		return ""
	}
	return f1.Data
}

// Find performs depth-first traversal looking for a node satisfying `pred`.
// Includes current node in the search.
// Stops at the first found node and returns it, wrapped in a new Finder.
func (f Finder) Find(pred pred.Predicate) (r Finder) {
	if f.Node == nil {
		return
	}

	var walker func(*html.Node) bool

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

// FindAll performs depth-first traversal like Find, but returns
// a stream of all found nodes, each wrapped in a Finder.
// FinderStream is really a readonly channel.
// Current node is included in the search.
func (f Finder) FindAll(pred pred.Predicate) FinderStream {
	ch := make(chan Finder)

	if f.Node == nil {
		close(ch)
		return ch
	}

	var walker func(*html.Node)

	walker = func(node *html.Node) {
		if pred(node) {
			ch <- FinderFromNode(node)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}

	go func() {
		walker(f.Node)
		close(ch)
	}()
	return ch
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

func (f Finder) FindSiblings(pred pred.Predicate) FinderStream {
	ch := make(chan Finder)

	if f.Node == nil {
		close(ch)
		return ch
	}

	go func() {
		for c := f.Node.NextSibling; c != nil; c = c.NextSibling {
			if pred(c) {
				ch <- Finder{c}
			}
		}
		close(ch)
	}()
	return ch
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
