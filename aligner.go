package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Aligner struct {
	w      io.Writer
	Margin *Margin
	Delim  *Delimiter
	lines  [][]string
	width  []int
}

func NewAligner(w io.Writer) *Aligner {
	return &Aligner{
		w:      w,
		Margin: NewMargin(),
		Delim:  NewDelimiter(),
	}
}

func (a *Aligner) appendLine(s string) {
	sp := a.Delim.Split(s)
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

func (a *Aligner) format(sp []string) string {
	if len(sp) == 1 {
		return sp[0]
	}
	for i := 0; i < len(sp); i++ {
		sp[i] = sp[i] + strings.Repeat(" ", a.width[i]-runewidth.StringWidth(sp[i]))
	}
	return strings.TrimSpace(a.Margin.Join(sp))
}

func (a *Aligner) Flush() error {
	for _, sp := range a.lines {
		_, err := fmt.Fprintln(a.w, a.format(sp))
		if err != nil {
			return err
		}
	}
	return nil
}
