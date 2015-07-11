package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yuya-takeyama/argf"
)

func guideToHelp() {
	os.Stderr.WriteString(`
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
	opt, err := ParseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}

	switch {
	case opt.IsHelp:
		longUsage()
		return 0
	case opt.IsVersion:
		version()
		return 0
	}

	a, err := NewAligner(opt)
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	r, err := argf.From(opt.Files)
	if err != nil {
		printErr(err)
		return 2
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
