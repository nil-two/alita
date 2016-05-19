package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yuya-takeyama/argf"
)

var (
	name    = "alita"
	version = "0.7.1"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [OPTION]... [FILE]...
Align FILE(s), or standard input.

Delimiter control:
  -d, --delimiter=DELIM    separate lines by DELIM
  -r, --regexp             DELIM is a regular expression
  -c, --count=COUNT        separate lines only COUNT times

Output control:
  -m, --margin=FORMAT      join cells by FORMAT
  -j, --justify=SEQUENCE   justify cells by SEQUENCE

Miscellaneous:
  -h, --help               display this help and exit
      --version            display version information and exit
`[1:], name)
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
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
		printUsage()
		return 0
	case opt.IsVersion:
		printVersion()
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
		guideToHelp()
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
