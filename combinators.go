package htmlx

import (
	"strings"

	"github.com/wkhere/htmlx/pred"
)

type (
	SplitFunc  func(Finder) FinderStream
	MapFunc    func(Finder)
	ReduceFunc func(prev, y Finder) (Finder, bool)
)

// Join performs depth-first join of streams generated by applying
// the splitFunc on every item from the current stream.
func (ff FinderStream) Join(split SplitFunc) FinderStream {
	ff2 := make(chan Finder)
	go func() {
		for f := range ff {
			for g := range split(f) {
				ff2 <- g
			}
		}
		close(ff2)
	}()
	return ff2
}

// Map maps the stream, changing each node wrapped by the finder
// according to the given function.
// Note that this function as a signature: `func (Finder)` and is supposed
// to change the internally copied html.Node struct whose pointer is wrapped
// by the finder.
// Also note that such an implicit copy is not made in the Reduce function,
// where we may return an original Finder.
func (ff FinderStream) Map(m MapFunc) FinderStream {
	ff2 := make(chan Finder)
	go func() {
		for x := range ff {
			x = x.Copy()
			m(x)
			ff2 <- x
		}
		close(ff2)
	}()
	return ff2
}

// Reduce reduces the stream according to the given combine function.
// It copies the first item of the stream, then for each item y,
// uses combine(prev, y) func to decide what should be copied to
// the output stream. There are three possibilities:
// (_, false)    no item is copied in this iteration
// (y, true)     item x is copied
// (newy, true)  combined item newx is copied.
// Note: for the last case, combine func should do explicit x.Copy()
// to create new html.Node struct.
func (ff FinderStream) Reduce(combine ReduceFunc) FinderStream {
	ff2 := make(chan Finder)
	go func() {
		x := <-ff
		ff2 <- x
		for {
			select {
			case y, ok := <-ff:
				if !ok {
					goto end
				}
				if y, ok := combine(x, y); ok {
					ff2 <- y
				}
				x = y
			}
		}
	end:
		close(ff2)
	}()
	return ff2
}

// AllText is a splitFunc returning a stream of all text nodes
// found descending from the given node, depth-first.
// The nodes containing only whitespaces are not included.
func AllText(f Finder) FinderStream {
	return f.FindAll(
		pred.TextCond(
			func(data string) bool {
				return strings.TrimSpace(data) != ""
			},
		),
	)
}

var _ SplitFunc = AllText