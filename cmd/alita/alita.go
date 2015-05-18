package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kusabashira/alita"
)

func usage() {
	os.Stderr.WriteString(`
Usage: alita [OPTION]... [FILE]...
Align FILE(s), or standard input.

Options:
  -m, --margin=FORMAT      join line by FORMAT (default: 1:1)
  -r, --regexp             DELIM is a regular expression
  -d, --delimiter=DELIM    use DELIM to separate line (default: /\s+/)
  -h, --help               show this help message
  -v, --version            print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.2.0
`[1:])
}

func do(a *alita.Aligner, r io.Reader) error {
	if err := a.ReadAll(r); err != nil {
		return err
	}
	return a.Flush()
}

func _main() error {
	var isHelp, isVersion bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isVersion, "v", false, "")
	flag.BoolVar(&isVersion, "version", false, "")

	a := alita.NewAligner(os.Stdout)
	flag.Var(a.Margin, "m", "")
	flag.Var(a.Margin, "margin", "")
	flag.Var(a.Delim, "d", "")
	flag.Var(a.Delim, "delimiter", "")
	flag.BoolVar(&a.Delim.UseRegexp, "r", false, "")
	flag.BoolVar(&a.Delim.UseRegexp, "regexp", false, "")

	flag.Usage = usage
	flag.Parse()
	switch {
	case isHelp:
		usage()
		return nil
	case isVersion:
		version()
		return nil
	}

	if flag.NArg() < 1 {
		return do(a, os.Stdin)
	}

	var input []io.Reader
	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer f.Close()
		input = append(input, f)
	}
	return do(a, io.MultiReader(input...))
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, "alita:", err)
		os.Exit(1)
	}
}
