package main

import (
	"regexp"
	"strings"
)

var SPACES = regexp.MustCompile(`\s+`)

type Delimiter struct {
	re        *regexp.Regexp
	Count     int
	UseRegexp bool
}

func NewDelimiter(expr string, useRegexp bool, count int) (d *Delimiter, err error) {
	d = &Delimiter{
		Count:     count,
		UseRegexp: useRegexp,
	}
	switch {
	case expr == "":
		d.re = nil
	case d.UseRegexp:
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

func NewDelimiterDefault() *Delimiter {
	d, _ := NewDelimiter("", false, -1)
	return d
}

func (d *Delimiter) Split(s string) []string {
	if d.re == nil {
		return SPACES.Split(s, -1)
	}

	count := d.Count
	if count < 1 {
		count = -1
	}
	useCount := count > 0

	matches := d.re.FindAllStringIndex(s, count)
	if len(matches) == 0 {
		return []string{strings.TrimSpace(s)}
	}

	a := make([]string, 0, len(matches)*2+1)
	beg, end := 0, 0
	for _, match := range matches {
		end = match[0]

		a = append(a, s[beg:end])
		beg, end = match[0], match[1]
		if useCount && len(a) >= count {
			break
		}

		a = append(a, s[beg:end])
		beg = match[1]
		if useCount && len(a) >= count {
			break
		}
	}
	a = append(a, s[beg:])
	for i := 0; i < len(a); i++ {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
