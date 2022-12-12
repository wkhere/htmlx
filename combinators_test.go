package htmlx

import (
	"strings"
	"testing"

	p "github.com/wkhere/htmlx/pred"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestStreamCombinators(t *testing.T) {
	f := testdata("simple2.html")
	top, _ := FinderFromData(f)
	f.Close()

	t.Run("Join", func(t *testing.T) {
		t.Parallel()

		ff := top.FindWithSiblings(p.Element(atom.Li)).Join(AllText)
		ee := ff.Collect()

		if len(ee) != 5 {
			t.Errorf("got %d, exp %d", len(ee), 5)
		}
		for i, e := range ee {
			if nt := e.Type; nt != html.TextNode {
				t.Errorf("ee[%d] type: got %v, exp %v", i, nt, html.TextNode)
			}
		}
		tab := []string{`1st`, `2nd`, `inner`, `3th`, `4th`}
		for i, data := range tab {
			if res := strings.TrimSpace(ee[i].Data); res != data {
				t.Errorf("ee[%d] data: got `%s`, exp `%s`", i, res, data)
			}
		}
	})

	t.Run("JoinMap", func(t *testing.T) {
		t.Parallel()

		ff := top.FindWithSiblings(p.Element(atom.Li)).Join(AllText).
			Map(
				func(x Finder) {
					x.Data = "foo"
				},
			)
		ee := ff.Collect()

		if len(ee) != 5 {
			t.Errorf("got %d, exp %d", len(ee), 5)
		}
		for i, e := range ee {
			if nt := e.Type; nt != html.TextNode {
				t.Errorf("ee[%d] type: got %v, exp %v", i, nt, html.TextNode)
			}
		}
		for i, e := range ee {
			if s, res := "foo", strings.TrimSpace(e.Data); res != s {
				t.Errorf("ee[%d] data: got `%s`, exp `%s`", i, res, s)
			}
		}
	})

	t.Run("JoinReduce", func(t *testing.T) {
		t.Parallel()

		ff := top.FindWithSiblings(p.Element(atom.Li)).Join(
			func(li Finder) FinderStream {
				return AllText(li).Reduce(
					func(prev, x Finder) (Finder, bool) {
						if strings.TrimSpace(x.Data) == "inner" {
							return x, false
						}
						return x, true
					},
				)
			},
		)
		ee := ff.Collect()

		if len(ee) != 4 {
			t.Errorf("got %d, exp %d", len(ee), 4)
		}
		for i, e := range ee {
			if nt := e.Type; nt != html.TextNode {
				t.Errorf("ee[%d] type: got %v, exp %v", i, nt, html.TextNode)
			}
		}
		tab := []string{`1st`, `2nd`, `3th`, `4th`}
		for i, data := range tab {
			if res := strings.TrimSpace(ee[i].Data); res != data {
				t.Errorf("ee[%d] data: got `%s`, exp `%s`", i, res, data)
			}
		}
	})
}
