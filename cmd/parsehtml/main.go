package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wkhere/htmlx/debug"
	"golang.org/x/net/html"
)

var (
	compactSpaces = flag.Bool("compact-spaces", true,
		"print SPC/LF etc if element data contains only spaces")
	trimAttr = flag.Bool("trim-attr", true,
		"don't print empty attributes")
)

func perr(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func usage() {
	perr("Usage:")
	const pname = "parsehtml"
	perr("\tparsehtml [flags] http(s)://url")
	perr("\tparsehtml [flags] file://path")
	perr("\tparsehtml [flags] path")
	perr("\tparsehtml [flags] - <file")
	perr("\tparsehtml [flags] one_input another_input")
	perr("Flags:")
	flag.PrintDefaults()
}

func die(err error) {
	perr("parsehtml:", err)
	os.Exit(1)
}

func dieIf(err error) {
	if err != nil {
		die(err)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		perr("nothing to parse; supply url, file or - as arguments")
		os.Exit(2)
	}

	for _, arg := range args {
		process(arg)
	}
}

func process(url string) {

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

	debug.Printer{
		CompactSpaces: *compactSpaces,
		TrimEmptyAttr: *trimAttr,
	}.Print(root)
}
