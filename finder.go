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
	// todo: check if strings.Reader is faster -> bigger benchmarks
	return FinderFromData(bytes.NewBufferString(s))
}

func (f Finder) IsEmpty() bool {
	return f.Node == nil
}

func (f Finder) Write(w io.Writer) error {
	return html.Render(w, f.Node)
}

func (f Finder) String() string {
	if f.Node == nil {
		return ""
	}
	// todo: check if strings.Builder is faster -> bigger benchmarks
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

func (ff FinderStream) First() Finder {
	return <-ff
}

func (ff FinderStream) Last() (f Finder) {
	for f = range ff {
	}
	return
}

func (ff FinderStream) Select(p pred.Predicate) Finder {
	for f := range ff {
		if p(f.Node) {
			return f
		}
	}
	return Finder{}
}

func (ff FinderStream) Filter(p pred.Predicate) FinderStream {
	ff2 := make(chan Finder)
	go func() {
		for f := range ff {
			if p(f.Node) {
				ff2 <- f
			}
		}
		close(ff2)
	}()
	return ff2
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
	// fixme: element can have many text children
	f1 := f.Find(pred.IsText())
	if f1.IsEmpty() {
		return ""
	}
	return f1.Data
}

// Find performs depth-first traversal looking for the node satisfying
// given precicate.
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
// a stream of all the found nodes satisfying the predicate,
// each wrapped in a Finder.
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

// FindSibling performs flat find of the first node satistying the predicate
// among current node's next (right) siblings.
// No recursion. Omits current node, starts from a first sibling.
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

// FindSiblings performs flat find of all the nodes satistying the predicate
// among current node's next (right) siblings.
// No recursion. Omits current node, starts from a first sibling.
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

// FindPrevSibling performs flat find of the first node
// satistying the predicate among current node's previous (left) siblings.
// No recursion. Omits current node, starts from a first sibling.
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

// FindPrevSiblings performs flat find of all the nodes
// satistying the predicate among current node's previous (left) siblings.
// No recursion. Omits current node, starts from a first sibling.
func (f Finder) FindPrevSiblings(pred pred.Predicate) FinderStream {
	ch := make(chan Finder)

	if f.Node == nil {
		close(ch)
		return ch
	}

	go func() {
		for c := f.Node.PrevSibling; c != nil; c = c.PrevSibling {
			if pred(c) {
				ch <- Finder{c}
			}
		}
		close(ch)
	}()
	return ch
}

// FindWithSiblings performs depth-first travelsal looking for the first
// of nodes satisfying given predicate.
// Then it continues flat find of all the siblings satisfying the predicate.
// All results are pushed into the returted stream of Finders.
func (f Finder) FindWithSiblings(pred pred.Predicate) (FinderStream, bool) {
	ch := make(chan Finder, 1)

	f = f.Find(pred)

	if f.Node == nil {
		close(ch)
		return ch, false
	}

	ch <- f
	go func() {
		for c := f.Node.NextSibling; c != nil; c = c.NextSibling {
			if pred(c) {
				ch <- Finder{c}
			}
		}
		close(ch)
	}()
	return ch, true
}
