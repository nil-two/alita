package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

func usage() {
	os.Stderr.WriteString(`
Usage: ali [OPTION]... [FILE]...
Align FILE(s), or standard input.

Options:
  -m, --margin=FORMAT      join line by FORMAT (default: 1:1)
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

var (
	DIGIT_ONLY             = regexp.MustCompile(`^\d+$`)
	COLON_SEPARATED_DIGITS = regexp.MustCompile(`^(\d+):(\d+)$`)
)

type Margin struct {
	left  int
	right int
}

func NewMargin() *Margin {
	return &Margin{
		left:  1,
		right: 1,
	}
}

func (m *Margin) String() string {
	return fmt.Sprint(*m)
}

func (m *Margin) Set(format string) error {
	switch {
	case DIGIT_ONLY.MatchString(format):
		n, err := strconv.Atoi(format)
		if err != nil {
			return err
		}
		m.left, m.right = n, n
	case COLON_SEPARATED_DIGITS.MatchString(format):
		a := COLON_SEPARATED_DIGITS.FindAllStringSubmatch(format, -1)
		left, err := strconv.Atoi(a[0][1])
		if err != nil {
			return err
		}
		right, err := strconv.Atoi(a[0][2])
		if err != nil {
			return err
		}
		m.left, m.right = left, right
	default:
		return fmt.Errorf("margin:", "invalid format:", format)
	}
	return nil
}

func (m *Margin) Join(a []string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := (m.left + m.right) * (len(a) / 2)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	lm, rm := strings.Repeat(" ", m.left), strings.Repeat(" ", m.right)

	b := make([]byte, n)
	bp := copy(b, a[0])
	for i := 2; i <= len(a); i += 2 {
		bp += copy(b[bp:], lm)
		bp += copy(b[bp:], a[i-1])
		if i != len(a) {
			bp += copy(b[bp:], rm)
			bp += copy(b[bp:], a[i])
		}
	}
	return string(b)
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
	w      io.Writer
	Margin *Margin
	Sep    *Separator
	lines  [][]string
	width  []int
}

func NewAligner(w io.Writer) *Aligner {
	return &Aligner{
		w:      w,
		Margin: NewMargin(),
		Sep:    NewSeparator(),
	}
}

func (a *Aligner) appendLine(s string) {
	sp := a.Sep.Split(s)
	a.lines = append(a.lines, sp)
	if len(sp) == 1 {
		return
	}
	for i, cell := range sp {
		if i == len(a.width) {
			a.width = append(a.width, 0)
		}

		w := runewidth.StringWidth(cell)
		if w > a.width[i] {
			a.width[i] = w
		}
	}
}

func (a *Aligner) ReadAll(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.appendLine(s.Text())
	}
	return s.Err()
}

func (a *Aligner) format(l []string) string {
	if len(l) == 1 {
		return l[0]
	}
	for i := 0; i < len(l); i++ {
		l[i] = l[i] + strings.Repeat(" ", a.width[i]-runewidth.StringWidth(l[i]))
	}
	return strings.TrimSpace(a.Margin.Join(l))
}

func (a *Aligner) Flush() error {
	for _, l := range a.lines {
		_, err := fmt.Fprintln(a.w, a.format(l))
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
	return a.Flush()
}

func _main() error {
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")

	a := NewAligner(os.Stdout)
	flag.Var(a.Margin, "m", "")
	flag.Var(a.Margin, "margin", "")
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
