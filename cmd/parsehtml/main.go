package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wkhere/htmlx/debug"
	"golang.org/x/net/html"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s url-or-file\n", os.Args[0])
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	os.Exit(1)
}

func dieIf(err error) {
	if err != nil {
		die(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}
	url := os.Args[1]
	if !strings.Contains(url, "://") {
		url = "file://" + url
	}

	var r io.ReadCloser
	var err error

	switch tokens := strings.Split(url, "://"); tokens[0] {
	case "file":
		if tokens[1] == "-" {
			r = os.Stdin
			break
		}
		r, err = os.Open(tokens[1])
		dieIf(err)

	case "http", "https":
		resp, err := http.Get(url)
		dieIf(err)
		r = resp.Body

	default:
		die(fmt.Errorf("unknown proto: %s", tokens[0]))
	}
	defer r.Close()

	root, err := html.Parse(r)
	dieIf(err)

	debug.PrintHTML(root)
}
