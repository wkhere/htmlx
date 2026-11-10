package main

import (
	"fmt"
	"io"

	"codeberg.org/wkhere/argp"
	pc "codeberg.org/wkhere/argp/parsecommon"
)

func defaults() config {
	return config{
		compactSpaces: true,
		trimAttr:      true,
	}
}

func usage(w io.Writer) {
	d := defaults()

	fmt.Fprintf(w,
		`Usage:
	parsehtml [flags] http(s)://url
	parsehtml [flags] file://path
	parsehtml [flags] path
	parsehtml [flags] - <file
	parsehtml [flags] one_input another_input
Flags:
	--compact-spaces   print SPC/LF etc if element data contains only spaces (default: %[1]v)
	--trim-attr        don't print empty attributes (default: %[1]v)
	-h, --help         show this help and exit
`,
		d.compactSpaces,
		d.trimAttr,
	)

}

func parseArgs(args []string) (c config, err error) {

	c = defaults()

	var p = &argp.PState{Args: args}
	var rest []string

loop:
	for ; len(p.Args) > 0 && p.Err == nil; p.Next() {
		switch {

		case p.MatchBoolFlag("", "--compact-spaces"):
			pc.ParseBool(&c.compactSpaces, p)

		case p.MatchBoolFlag("", "--trim-attr"):
			pc.ParseBool(&c.trimAttr, p)

		case p.MatchJustFlag("-h", "--help"):
			c.help = usage
			return c, nil

		case p.Arg() == "--":
			rest = append(rest, p.Args[1:]...)
			break loop

		case p.MatchAnyFlag():
			p.Errorf("unknown flag %s", p.Flag)

		default:
			rest = append(rest, p.Arg())
		}
	}

	if p.Err != nil {
		return c, p.Err
	}
	if len(rest) == 0 {
		return c, fmt.Errorf("nothing to parse; supply url, file or - as arguments")
	}

	c.args = rest
	return c, nil
}
