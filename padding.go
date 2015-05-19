package alita

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type Justify int

const (
	JustLeft Justify = iota
	JustRight
	JustCenter
)

func (j Justify) Just(width int, s string) string {
	w := runewidth.StringWidth(s)
	switch j {
	case JustLeft:
		return s + strings.Repeat(" ", width-w)
	case JustRight:
		return strings.Repeat(" ", width-w) + s
	case JustCenter:
		n := width - w
		l, r := n/2, n/2
		if n%2 != 0 {
			r += 1
		}
		return strings.Repeat(" ", l) + s + strings.Repeat(" ", r)
	}
	return s + strings.Repeat(" ", width-w)
}

type Padding struct {
	justfies []Justify
	width    []int
}

func NewPadding() *Padding {
	return &Padding{}
}

func (p *Padding) SetJustfies(a []Justify) {
	p.justfies = a
}

func (p *Padding) UpdateWidth(a []string) {
	if len(a) == 1 {
		return
	}
	for i, s := range a {
		w := runewidth.StringWidth(s)
		switch {
		case i == len(p.width):
			p.width = append(p.width, w)
		case w > p.width[i]:
			p.width[i] = w
		}
	}
}

func (p *Padding) justKind(i int) Justify {
	if len(p.justfies) == 0 {
		return JustLeft
	}
	j := (i-1)%(len(p.justfies)-1) + 1
	return p.justfies[j]
}

func (p *Padding) Format(a []string) []string {
	if len(a) == 1 {
		return a
	}
	for i := 0; i < len(a); i++ {
		a[i] = p.justKind(i).Just(p.width[i], a[i])
	}
	return a
}
