package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wkhere/htmlx/pp"
	"golang.org/x/net/html"
)

type config struct {
	compactSpaces bool
	trimAttr      bool

	args []string
	help func(io.Writer)
}

func process(url string, p *pp.Printer) (err error) {

	if !strings.Contains(url, "://") {
		url = "file://" + url
	}

	var r io.ReadCloser

	switch tokens := strings.Split(url, "://"); tokens[0] {
	case "file":
		if tokens[1] == "-" {
			r = os.Stdin
			break
		}
		r, err = os.Open(tokens[1])
		if err != nil {
			return err
		}

	case "http", "https":
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		r = resp.Body

	default:
		return fmt.Errorf("unknown proto: %s", tokens[0])
	}
	defer r.Close()

	root, err := html.Parse(r)
	if err != nil {
		return err
	}

	p.Print(os.Stdout, root)

	return nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help(os.Stdout)
		os.Exit(0)
	}

	p := pp.Printer{
		CompactSpaces: conf.compactSpaces,
		TrimEmptyAttr: conf.trimAttr,
	}

	for _, arg := range conf.args {
		err = process(arg, &p)
		if err != nil {
			die(1, err)
		}
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "parsehtml:", err)
	os.Exit(code)
}
