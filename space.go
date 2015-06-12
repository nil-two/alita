package main

import (
	"strings"
	"unicode"
)

type Space struct {
	tabWidth  int
	headWidth int
	headSpace string
}

func NewSpace() *Space {
	return &Space{
		tabWidth:  8,
		headWidth: int(^uint(0) >> 1),
	}
}

func (s *Space) UpdateHeadWidth(t string) {
	if s.headWidth < 1 {
		return
	}

	w, i := 0, 0
	for _, c := range t {
		switch c {
		case ' ':
			w += 1
			i += 1
		case '\t':
			w += s.tabWidth
			i += 1
		default:
			goto END
		}
	}
END:
	if w < s.headWidth {
		s.headWidth = w
		s.headSpace = t[:i]
	}
}

func (s *Space) Trim(t string) string {
	return strings.TrimSpace(t)
}

func (s *Space) Adjust(t string) string {
	return s.headSpace + strings.TrimRightFunc(t, unicode.IsSpace)
}
