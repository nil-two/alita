package main

import (
	"regexp"
	"strings"
)

var SPACES = regexp.MustCompile(`\s+`)

type Delimiter struct {
	re    *regexp.Regexp
	count int
}

func NewDelimiter(expr string, useRegexp bool, count int) (d *Delimiter, err error) {
	d = &Delimiter{}
	d.count = count
	if d.count < 1 {
		d.count = -1
	}
	switch {
	case expr == "":
		d.re = nil
		if d.count != -1 {
			d.count += 1
		}
	case useRegexp:
		d.re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	default:
		expr = regexp.QuoteMeta(expr)
		d.re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (d *Delimiter) Split(s string) []string {
	if d.re == nil {
		return SPACES.Split(s, d.count)
	}

	matches := d.re.FindAllStringIndex(s, d.count)
	if len(matches) == 0 {
		return []string{strings.TrimSpace(s)}
	}

	a := make([]string, 0, len(matches)*2+1)
	beg, end := 0, 0
	useCount := d.count > 0
	for _, match := range matches {
		end = match[0]

		a = append(a, s[beg:end])
		beg, end = match[0], match[1]
		if useCount && len(a) >= d.count {
			break
		}

		a = append(a, s[beg:end])
		beg = match[1]
		if useCount && len(a) >= d.count {
			break
		}
	}
	a = append(a, s[beg:])
	for i := 0; i < len(a); i++ {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
