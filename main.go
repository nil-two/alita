package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

func usage() {
	os.Stderr.WriteString(`
Usage: gotran [OPTION]... [FILE]...
Align FILE(s), or standard input.

Options:
	--help       show this help message
	--version    print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.0
`[1:])
}

var SPACES = regexp.MustCompile(`\s+`)

type Aligner struct {
	lines [][]string
	width []int
}

func NewAligner() *Aligner {
	return &Aligner{}
}

func (a *Aligner) AppendLine(s string) {
	sp := SPACES.Split(s, -1)
	for i, cell := range sp {
		if i == len(a.width) {
			a.width = append(a.width, 0)
		}

		w := runewidth.StringWidth(cell)
		if w > a.width[i] {
			a.width[i] = w
		}
	}
	a.lines = append(a.lines, sp)
}

func (a *Aligner) ReadAll(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.AppendLine(s.Text())
	}
	return nil
}

func (a *Aligner) format(l []string) string {
	for i := 0; i < len(l); i++ {
		l[i] = l[i] + strings.Repeat(" ", a.width[i]-runewidth.StringWidth(l[i]))
	}
	return strings.TrimSpace(strings.Join(l, " "))
}

func (a *Aligner) Flush(w io.Writer) error {
	for _, l := range a.lines {
		_, err := fmt.Fprintln(w, a.format(l))
		if err != nil {
			return err
		}
	}
	return nil
}

func do(a *Aligner, r io.Reader) error {
	if err := a.ReadAll(r); err != nil {
		return err
	}
	return a.Flush(os.Stdout)
}

func _main() error {
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")
	flag.Usage = usage
	flag.Parse()
	switch {
	case *isHelp:
		usage()
		return nil
	case *isVersion:
		version()
		return nil
	}

	a := NewAligner()
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
		fmt.Fprintln(os.Stderr, "ali:", err)
		os.Exit(1)
	}
}
