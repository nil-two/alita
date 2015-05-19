package alita

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type Padding struct {
	width []int
}

func NewPadding() *Padding {
	return &Padding{}
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

func (p *Padding) Format(a []string) []string {
	if len(a) == 1 {
		return a
	}
	for i := 0; i < len(a); i++ {
		w := p.width[i] - runewidth.StringWidth(a[i])
		a[i] = a[i] + strings.Repeat(" ", w)
	}
	return a
}
