package main

import (
	"fmt"
	"io"

	"github.com/spf13/pflag"
)

func parseArgs(args []string) (c config, err error) {
	var help bool

	fs := pflag.NewFlagSet("parsehtml", pflag.ContinueOnError)
	fs.SortFlags = false

	fs.BoolVar(&c.compactSpaces, "compact-spaces", true,
		"print SPC/LF etc if element data contains only spaces")

	fs.BoolVar(&c.trimAttr, "trim-attr", true,
		"don't print empty attributes")

	fs.BoolVarP(&help, "help", "h", false, "show this help and exit")

	err = fs.Parse(args)
	if err != nil {
		return c, err
	}

	if help {
		c.help = func(w io.Writer) {
			fs.SetOutput(w)
			p := func(s string) { fmt.Fprintln(fs.Output(), s) }
			p("Usage:")
			p("\tparsehtml [flags] http(s)://url")
			p("\tparsehtml [flags] file://path")
			p("\tparsehtml [flags] path")
			p("\tparsehtml [flags] - <file")
			p("\tparsehtml [flags] one_input another_input")
			p("Flags:")
			fs.PrintDefaults()
		}
		return c, nil
	}

	c.args = fs.Args()

	if len(c.args) == 0 {
		return c, fmt.Errorf("nothing to parse; supply url, file or - as arguments")
	}

	return c, nil
}
