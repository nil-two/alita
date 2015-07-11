package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yuya-takeyama/argf"
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

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "alita:", err)
}

func do(a *Aligner, r io.Reader) error {
	if err := a.ReadAll(r); err != nil {
		return err
	}
	return a.Flush(os.Stdout)
}

func _main() int {
	a := NewAlignerDefault()
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
		return 0
	case isVersion:
		version()
		return 0
	}

	r, err := argf.From(flag.Args())
	if err != nil {
		printErr(err)
		return 1
	}
	if err = do(a, r); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
