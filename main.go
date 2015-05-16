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
  -s, --separator=STRING   use STRING to separate line (default: spaces)
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
	isDefault bool
}

func NewSeparator() *Separator {
	return &Separator{
		isDefault: true,
	}
}

func (s *Separator) String() string {
	return fmt.Sprint(*s)
}

func (s *Separator) Set(expr string) error {
	re, err := regexp.Compile(regexp.QuoteMeta(expr))
	if err != nil {
		return err
	}
	s.isDefault = false
	s.re = re
	return nil
}

func (s *Separator) Split(t string) []string {
	if s.isDefault {
		return SPACES.Split(t, -1)
	}

	ils := s.re.FindAllStringIndex(t, -1)
	if len(ils) == 0 {
		return []string{t}
	}

	nls := make([]string, 0, len(ils)*2+1)
	nls = append(nls, strings.TrimSpace(t[:ils[0][0]]))
	for i := 1; i < len(ils); i++ {
		nls = append(nls, strings.TrimSpace(t[ils[i-1][0]:ils[i-1][1]]))
		nls = append(nls, strings.TrimSpace(t[ils[i-1][1]:ils[i][0]]))
	}
	nls = append(nls, strings.TrimSpace(t[ils[len(ils)-1][0]:ils[len(ils)-1][1]]))
	nls = append(nls, strings.TrimSpace(t[ils[len(ils)-1][1]:]))
	return nls
}

type Aligner struct {
	sep   *Separator
	lines [][]string
	width []int
}

func NewAligner() *Aligner {
	return &Aligner{
		sep: NewSeparator(),
	}
}

func (a *Aligner) appendLine(s string) {
	sp := a.sep.Split(s)
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
	a.Flush(os.Stdout)
	return nil
}

func _main() error {
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")

	a := NewAligner()
	flag.Var(a.sep, "s", "")
	flag.Var(a.sep, "separator", "")

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
