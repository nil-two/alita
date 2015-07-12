package main

import (
	"bufio"
	"fmt"
	"io"
)

type Aligner struct {
	delimiter *Delimiter
	margin    *Margin
	padding   *Padding
	space     *Space
	lines     [][]string
}

func NewAligner(opt *Option) (*Aligner, error) {
	if opt == nil {
		opt = &Option{}
	}

	d, err := NewDelimiter(opt.Delimiter, opt.UseRegexp, opt.Count)
	if err != nil {
		return nil, err
	}
	m, err := NewMargin(opt.Margin)
	if err != nil {
		return nil, err
	}
	p, err := NewPadding(opt.Justify)
	if err != nil {
		return nil, err
	}
	s := NewSpace()
	return &Aligner{
		delimiter: d,
		padding:   p,
		margin:    m,
		space:     s,
	}, nil
}

func (a *Aligner) AppendLine(s string) {
	sp := a.delimiter.Split(a.space.Trim(s))
	a.lines = append(a.lines, sp)

	if len(sp) > 1 {
		a.space.UpdateHeadWidth(s)
		a.padding.UpdateWidth(sp)
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
	return a.space.Adjust(a.margin.Join(a.padding.Format(sp)))
}

func (a *Aligner) Flush(out io.Writer) error {
	for _, sp := range a.lines {
		_, err := fmt.Fprintln(out, a.format(sp))
		if err != nil {
			return err
		}
	}
	return nil
}
