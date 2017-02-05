package main

import (
	"bufio"
	"io"
)

type Option struct {
	Delimiter string
	UseRegexp bool
	Count     int
	Margin    string
	Justify   string
}

type Aligner struct {
	delimiter *Delimiter
	margin    *Margin
	padding   *Padding
	space     *Space
	cells     [][]string
}

func NewAligner(opt *Option) (a *Aligner, err error) {
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
		margin:    m,
		padding:   p,
		space:     s,
	}, nil
}

func (a *Aligner) AddRow(s string) {
	row := a.delimiter.Split(a.space.Trim(s))
	a.cells = append(a.cells, row)

	if len(row) > 1 {
		a.space.UpdateHeadWidth(s)
		a.padding.UpdateWidth(row)
	}
}

func (a *Aligner) ReadAll(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.AddRow(s.Text())
	}
	return s.Err()
}

func (a *Aligner) format(cells []string) string {
	if len(cells) == 1 {
		return cells[0]
	}
	return a.space.Adjust(a.margin.Join(a.padding.Format(cells)))
}

func (a *Aligner) Flush(w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, row := range a.cells {
		if _, err := bw.WriteString(a.format(row) + "\n"); err != nil {
			return err
		}
	}
	return bw.Flush()
}
