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
Usage: ali [OPTION]... [FILE]...
Align FILE(s), or standard input.

Options:
  -r, --regexp             PATTERN is a regular expression
  -s, --separator=PATTERN  use PATTERN to separate line (default: /\s+/)
  --help                   show this help message
  --version                print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.0
`[1:])
}

var SPACES = regexp.MustCompile(`\s+`)

type Separator struct {
	re        *regexp.Regexp
	UseRegexp bool
}

func NewSeparator() *Separator {
	return &Separator{}
}

func (s *Separator) String() string {
	return fmt.Sprint(*s)
}

func (s *Separator) Set(expr string) error {
	if !s.UseRegexp {
		expr = regexp.QuoteMeta(expr)
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	s.re = re
	return nil
}

func (s *Separator) Split(t string) []string {
	if s.re == nil {
		return SPACES.Split(t, -1)
	}

	matches := s.re.FindAllStringIndex(t, -1)
	if len(matches) == 0 {
		return []string{t}
	}

	sls := make([]string, 0, len(matches)*2+1)
	beg, end := 0, 0
	for _, match := range matches {
		end = match[0]
		sls = append(sls, t[beg:end])
		beg, end = match[0], match[1]
		sls = append(sls, t[beg:end])
		beg = match[1]
	}
	sls = append(sls, t[beg:])
	for i := 0; i < len(sls); i++ {
		sls[i] = strings.TrimSpace(sls[i])
	}
	return sls
}

type Aligner struct {
	Sep   *Separator
	lines [][]string
	width []int
}

func NewAligner() *Aligner {
	return &Aligner{
		Sep: NewSeparator(),
	}
}

func (a *Aligner) appendLine(s string) {
	sp := a.Sep.Split(s)
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
		a.appendLine(s.Text())
	}
	return s.Err()
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
	a.Flush(os.Stdout)
	return nil
}

func _main() error {
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")

	a := NewAligner()
	flag.Var(a.Sep, "s", "")
	flag.Var(a.Sep, "separator", "")
	flag.BoolVar(&a.Sep.UseRegexp, "r", false, "")
	flag.BoolVar(&a.Sep.UseRegexp, "regexp", false, "")

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
