package alita

import (
	"bufio"
	"fmt"
	"io"
)

type Aligner struct {
	w         io.Writer
	Space     *Space
	Margin    *Margin
	Delimiter *Delimiter
	Padding   *Padding
	lines     [][]string
}

func NewAligner(w io.Writer) *Aligner {
	return &Aligner{
		w:         w,
		Space:     NewSpace(),
		Margin:    NewMargin(),
		Delimiter: NewDelimiter(),
		Padding:   NewPadding(),
	}
}

func (a *Aligner) SetOutput(w io.Writer) {
	a.w = w
}

func (a *Aligner) AppendLine(s string) {
	sp := a.Delimiter.Split(a.Space.Strip(s))
	a.lines = append(a.lines, sp)

	if len(sp) > 1 {
		a.Space.UpdateHeadWidth(s)
		a.Padding.UpdateWidth(sp)
	}
}

func (a *Aligner) ReadAll(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.AppendLine(s.Text())
	}
	return s.Err()
}

func (a *Aligner) format(sp []string) string {
	if len(sp) == 1 {
		return sp[0]
	}
	return a.Space.Adjust(a.Margin.Join(a.Padding.Format(sp)))
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
