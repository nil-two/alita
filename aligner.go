package alita

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Aligner struct {
	w       io.Writer
	Margin  *Margin
	Delim   *Delimiter
	Padding *Padding
	lines   [][]string
}

func NewAligner(w io.Writer) *Aligner {
	return &Aligner{
		w:       w,
		Margin:  NewMargin(),
		Delim:   NewDelimiter(),
		Padding: NewPadding(),
	}
}

func (a *Aligner) appendLine(s string) {
	sp := a.Delim.Split(s)
	a.lines = append(a.lines, sp)

	a.Padding.UpdateWidth(sp)
}

func (a *Aligner) ReadAll(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.appendLine(s.Text())
	}
	return s.Err()
}

func (a *Aligner) Flush() error {
	for _, sp := range a.lines {
		s := strings.TrimSpace(a.Margin.Join(a.Padding.Format(sp)))
		_, err := fmt.Fprintln(a.w, s)
		if err != nil {
			return err
		}
	}
	return nil
}
