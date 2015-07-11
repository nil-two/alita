package main

import (
	"fmt"
	"regexp"
	"strings"
)

var SPACES = regexp.MustCompile(`\s+`)

type Delimiter struct {
	re        *regexp.Regexp
	Count     int
	UseRegexp bool
}

func NewDelimiter(expr string, useRegexp bool, count int) (*Delimiter, error) {
	d := &Delimiter{
		Count:     count,
		UseRegexp: useRegexp,
	}
	if err := d.Set(expr); err != nil {
		return nil, err
	}
	return d, nil
}

func NewDelimiterDefault() *Delimiter {
	d, _ := NewDelimiter("", false, -1)
	return d
}

func (d *Delimiter) String() string {
	return fmt.Sprint(*d)
}

func (d *Delimiter) Set(expr string) error {
	if expr == "" {
		d.re = nil
		return nil
	}
	if !d.UseRegexp {
		expr = regexp.QuoteMeta(expr)
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	d.re = re
	return nil
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
