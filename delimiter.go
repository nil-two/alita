package alita

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

func NewDelimiter() *Delimiter {
	return &Delimiter{
		Count:     -1,
		UseRegexp: false,
	}
}

func (d *Delimiter) String() string {
	return fmt.Sprint(*d)
}

func (d *Delimiter) Set(expr string) error {
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

	matches := d.re.FindAllStringIndex(s, -1)
	if len(matches) == 0 {
		return []string{strings.TrimSpace(s)}
	}

	a := make([]string, 0, len(matches)*2+1)
	beg, end := 0, 0
	for _, match := range matches {
		end = match[0]
		a = append(a, s[beg:end])
		beg, end = match[0], match[1]
		a = append(a, s[beg:end])
		beg = match[1]
	}
	a = append(a, s[beg:])
	for i := 0; i < len(a); i++ {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
