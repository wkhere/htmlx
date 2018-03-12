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
		"print SPC if element data contains only spaces")
	trimAttr = flag.Bool("trim-attr", true,
		"don't print empty attributes")
)

func perrf(format string, vs ...interface{}) {
	fmt.Fprintf(os.Stderr, format, vs...)
}

func usage() {
	perrf("Usage:\n")
	perrf("\t%s [flags] http(s)://url\n", os.Args[0])
	perrf("\t%s [flags] file://path\n", os.Args[0])
	perrf("\t%s [flags] path\n", os.Args[0])
	perrf("\t%s [flags] - <file \n", os.Args[0])
	perrf("\t%s [flags] one_input another_input \n", os.Args[0])
	perrf("Flags:\n")
	flag.PrintDefaults()
}

func die(err error) {
	perrf("%s: %v\n", os.Args[0], err)
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
		perrf("nothing to parse; supply url, file or - as arguments\n")
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

	debug.NewPrinter(
		debug.CompactSpaces(*compactSpaces),
		debug.TrimEmptyAttr(*trimAttr),
	).Print(root)
}
