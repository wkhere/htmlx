package pp

import (
	"io"
	"strconv"
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

func (p Printer) Print(w io.Writer, top *html.Node) {

	var f func(*html.Node, int)

	f = func(node *html.Node, i int) {
		var dataRepr, attrRepr string
		if p.CompactSpaces && len(strings.TrimSpace(node.Data)) == 0 {
			dataRepr = ppSpaces(node.Data)
		} else {
			dataRepr = "`" + node.Data + "`"
		}
		attrRepr = ppAttr(node.Attr)

		io.WriteString(w, strings.Repeat(" ", i*2))
		io.WriteString(w, "T:")
		io.WriteString(w, nodeTypes[node.Type])
		io.WriteString(w, " D:")
		io.WriteString(w, dataRepr)
		if attrRepr != "" || !p.TrimEmptyAttr {
			io.WriteString(w, " A:")
			io.WriteString(w, attrRepr)
		}
		w.Write([]byte{'\n'})

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c, i+1)
		}
	}

	f(top, 0)
}

func ppAttr(aa []html.Attribute) string {
	if len(aa) == 0 {
		return ""
	}
	b := new(strings.Builder)
	b.WriteByte('[')
	bppAttr1(b, aa[0])
	for _, a := range aa[1:] {
		b.WriteByte(' ')
		bppAttr1(b, a)
	}
	b.WriteByte(']')
	return b.String()
}

func bppAttr1(b *strings.Builder, a html.Attribute) {
	if a.Namespace != "" {
		b.WriteString(a.Namespace)
		b.WriteByte(':')
	}
	b.WriteString(a.Key)
	if a.Val != "" {
		b.WriteString(`="`)
		b.WriteString(a.Val)
		b.WriteByte('"')
	}
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

	bpp := func(b *strings.Builder, tok token) {
		if tok.cnt == 1 {
			b.WriteString(tok.val)
			return
		}
		b.WriteString(tok.val)
		b.WriteByte('x')
		b.WriteString(strconv.Itoa(tok.cnt))
	}

	b := new(strings.Builder)
	bpp(b, r[0])
	for _, tok := range r[1:] {
		b.WriteByte(',')
		bpp(b, tok)
	}
	return b.String()
}
