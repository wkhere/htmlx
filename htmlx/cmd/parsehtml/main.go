package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/wkhere/x/htmlx/debug"
	"golang.org/x/net/html"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s url\n", os.Args[0])
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}
	url := os.Args[1]
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	var err error

	resp, err := http.Get(url)
	dieIf(err)
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	dieIf(err)

	debug.PrintHTML(root)
}
