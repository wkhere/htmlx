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
			dataRepr = "D:SPC"
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
