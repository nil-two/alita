package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kusabashira/alita"
)

func shortUsage() {
	os.Stderr.WriteString(`
Usage: alita [OPTION]... [FILE]...
Try 'alita --help' for more information.
`[1:])
}

func longUsage() {
	os.Stderr.WriteString(`
Usage: alita [OPTION]... [FILE]...
Align FILE(s), or standard input.

Delimiter control:
  -c, --count=COUNT        delimit line COUNT times
  -r, --regexp             DELIM is a regular expression
  -d, --delimiter=DELIM    delimit line by DELIM

Output control:
  -m, --margin=FORMAT      join cells by FORMAT
  -j, --justfy=SEQUENCE    justfy cells by SEQUENCE

Miscellaneous:
  -h, --help               show this help message
      --version            print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.6.0
`[1:])
}

func do(a *alita.Aligner, r io.Reader) error {
	if err := a.ReadAll(r); err != nil {
		return err
	}
	return a.Flush()
}

func _main() error {
	a := alita.NewAligner(os.Stdout)
	flag.IntVar(&a.Delimiter.Count, "c", -1, "")
	flag.IntVar(&a.Delimiter.Count, "count", -1, "")
	flag.BoolVar(&a.Delimiter.UseRegexp, "r", false, "")
	flag.BoolVar(&a.Delimiter.UseRegexp, "regexp", false, "")
	flag.Var(a.Delimiter, "d", "")
	flag.Var(a.Delimiter, "delimiter", "")
	flag.Var(a.Margin, "m", "")
	flag.Var(a.Margin, "margin", "")
	flag.Var(a.Padding, "j", "")
	flag.Var(a.Padding, "justfy", "")

	var isHelp, isVersion bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isVersion, "version", false, "")
	flag.Usage = shortUsage
	flag.Parse()
	switch {
	case isHelp:
		longUsage()
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
