package debug

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func PrintHTML(top *html.Node) {
	var f func(*html.Node, int)

	f = func(node *html.Node, i int) {
		fmt.Printf("%sT:%d D:`%s` A:%q\n", strings.Repeat(" ", i*2),
			node.Type, node.Data, node.Attr)

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c, i+1)
		}
	}

	f(top, 0)
	return
}
