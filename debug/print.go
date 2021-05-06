package debug

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var nodeTypes = map[html.NodeType]string{
	html.ErrorNode:    "ERR",
	html.DocumentNode: "DOC",
	html.DoctypeNode:  "DOCTYPE",
	html.ElementNode:  "ELEM",
	html.TextNode:     "TEXT",
	html.CommentNode:  "COMMENT",
}

type Printer struct {
	CompactSpaces bool
	TrimEmptyAttr bool
}

func (p Printer) Print(top *html.Node) {

	var f func(*html.Node, int)

	f = func(node *html.Node, i int) {
		var dataRepr, attrRepr string
		if p.CompactSpaces && len(strings.TrimSpace(node.Data)) == 0 {
			dataRepr = "D:" + ppSpaces(node.Data)
		} else {
			dataRepr = fmt.Sprintf("D:`%s`", node.Data)
		}
		if !(p.TrimEmptyAttr && len(node.Attr) == 0) {
			attrRepr = fmt.Sprintf("A:%q", node.Attr)
		}

		fmt.Printf("%sT:%s %s %s\n", strings.Repeat(" ", i*2),
			nodeTypes[node.Type], dataRepr, attrRepr)

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c, i+1)
		}
	}

	f(top, 0)
}

// ppSpaces pretty prints various whitespaces, possibly with counters.
// Input string `s` must be already checked that it contains whitespaces only.
func ppSpaces(s string) string {
	type token struct {
		val string
		cnt int
	}
	a := make([]token, 0, len(s))

	for _, c := range s {
		var t string
		switch c {
		case '\n':
			t = "LF"
		case '\r':
			t = "CR"
		case ' ':
			t = "SPC"
		default:
			t = "WS"
		}
		a = append(a, token{val: t, cnt: 1})
	}

	if len(a) == 0 {
		return ""
	}

	r := a[:1]
	i := 0
	for _, tok := range a[1:] {
		if r[i].val == tok.val {
			r[i].cnt++
		} else {
			r = append(r, tok)
			i++
		}
	}
	a = r

	pp := func(tok token) string {
		if tok.cnt == 1 {
			return tok.val
		}
		return fmt.Sprintf("%sx%d", tok.val, tok.cnt)
	}

	b := new(strings.Builder)

	b.WriteString(pp(r[0]))
	for _, tok := range r[1:] {
		b.WriteByte(',')
		b.WriteString(pp(tok))
	}
	return b.String()
}
