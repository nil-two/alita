package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Justify int

const (
	JustLeft Justify = iota
	JustRight
	JustCenter
)

var JUSTFIES_SEQUENCE = regexp.MustCompile("^[lrc]+$")

func ParseJustifies(seq string) ([]Justify, error) {
	switch {
	case seq == "":
		return []Justify{JustLeft}, nil
	case JUSTFIES_SEQUENCE.MatchString(seq):
		js := make([]Justify, 0, len(seq))
		for _, ch := range seq {
			switch ch {
			case 'l':
				js = append(js, JustLeft)
			case 'r':
				js = append(js, JustRight)
			case 'c':
				js = append(js, JustCenter)
			}
		}
		return js, nil
	default:
		return nil, fmt.Errorf("padding: invalid format: %s", seq)
	}
}

func (j Justify) Just(width int, s string) string {
	w := runewidth.StringWidth(s)
	if width <= w {
		return s
	}
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

func NewPadding(seq string) (p *Padding, err error) {
	p = &Padding{}
	p.justfies, err = ParseJustifies(seq)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func NewPaddingDefault() *Padding {
	p, _ := NewPadding("l")
	return p
}

func (p *Padding) SetJustfies(a []Justify) {
	p.justfies = a
}

func (p *Padding) UpdateWidth(a []string) {
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
	switch len(p.justfies) {
	case 0:
		return JustLeft
	case 1:
		return p.justfies[0]
	}
	if i <= 0 {
		return p.justfies[0]
	}
	j := (i-1)%(len(p.justfies)-1) + 1
	return p.justfies[j]
}

func (p *Padding) Format(a []string) []string {
	for i := 0; i < len(a); i++ {
		j := p.justKind(i)
		w := 0
		if i < len(p.width) {
			w = p.width[i]
		}
		a[i] = j.Just(w, a[i])
	}
	return a
}
