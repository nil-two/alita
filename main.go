package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
Usage: alita [OPTION]... [FILE]...
Align FILE(s), or standard input.

Delimiter control:
  -d, --delimiter=DELIM    delimit line by DELIM
  -r, --regexp             DELIM is a regular expression
  -c, --count=COUNT        delimit line COUNT times

Output control:
  -m, --margin=FORMAT      join cells by FORMAT
  -j, --justify=SEQUENCE   justify cells by SEQUENCE

Miscellaneous:
  -h, --help               show this help message
      --version            print the version
`[1:])
}

func printVersion() {
	os.Stderr.WriteString(`
0.7.1
`[1:])
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "alita:", err)
}

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'alita --help' for more information.
`[1:])
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
		usage()
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
