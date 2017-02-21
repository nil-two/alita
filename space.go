package main

import (
	"strings"
	"unicode"
)

var IntMax = int(^uint(0) >> 1)

type Space struct {
	tabWidth     int
	leadingWidth int
	leadingSpace string
}

func NewSpace() *Space {
	return &Space{
		tabWidth:     8,
		leadingWidth: IntMax,
	}
}

func (s *Space) UpdateLeadingWidth(t string) {
	if s.leadingWidth < 1 {
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
	if w < s.leadingWidth {
		s.leadingWidth = w
		s.leadingSpace = t[:i]
	}
}

func (s *Space) Trim(t string) string {
	return strings.TrimSpace(t)
}

func (s *Space) Adjust(t string) string {
	return s.leadingSpace + strings.TrimRightFunc(t, unicode.IsSpace)
}
