package alita

import (
	"reflect"
	"regexp"
	"testing"
)

type DelimiterSetTest struct {
	useRegexp bool
	expr      string
	re        *regexp.Regexp
}

var indexTestsDelimiterSet = []DelimiterSetTest{
	{false, `=`, regexp.MustCompile(`=`)},
	{false, `=+`, regexp.MustCompile(`=\+`)},
	{false, `-*>`, regexp.MustCompile(`-\*>`)},

	{true, `=`, regexp.MustCompile(`=`)},
	{true, `=+`, regexp.MustCompile(`=+`)},
	{true, `-*>`, regexp.MustCompile(`-*>`)},
}

func TestDelimiterSet(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsDelimiterSet {
		d.UseRegexp = test.useRegexp
		if err := d.Set(test.expr); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.expr, err)
		}
		if !reflect.DeepEqual(d.re, test.re) {
			t.Errorf("got %q, want %q",
				d.re, test.re)
		}
	}
}

type DelimiterSplitDefaultTest struct {
	src string
	dst []string
}

var indexTestsDelimiterSplitDefault = []DelimiterSplitDefaultTest{
	{"a", []string{"a"}},
	{"a b", []string{"a", "b"}},
	{"a b c", []string{"a", "b", "c"}},
	{"a b c d", []string{"a", "b", "c", "d"}},

	{"a  b", []string{"a", "b"}},
	{"a  b c", []string{"a", "b", "c"}},
}

func TestDelimiterDefaultSplit(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsDelimiterSplitDefault {
		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q).Split(%q) = %q; want %q",
				d.re, test.src, actual, expect)
		}
	}
}

type DelimiterSplitTest struct {
	useRegexp bool
	expr      string
	src       string
	dst       []string
}

var indexTestsDelimiterSplit = []DelimiterSplitTest{
	{false, `=`, "n", []string{"n"}},
	{false, `=`, "n=", []string{"n", "=", ""}},
	{false, `=`, "=n", []string{"", "=", "n"}},

	{false, `=`, "n=100", []string{"n", "=", "100"}},
	{false, `=`, "n  =  100", []string{"n", "=", "100"}},
	{false, `=`, "n  = 100=200", []string{"n", "=", "100", "=", "200"}},
	{false, `=`, "n=100 =  200", []string{"n", "=", "100", "=", "200"}},

	{false, `=`, "n==100", []string{"n", "=", "", "=", "100"}},
	{false, `=`, "n===100", []string{"n", "=", "", "=", "", "=", "100"}},

	{false, `=>`, "a=>b=>c", []string{"a", "=>", "b", "=>", "c"}},
	{false, `=>`, "a => b => c", []string{"a", "=>", "b", "=>", "c"}},

	{true, `=+>`, "a => b",
		[]string{"a", "=>", "b"}},
	{true, `=+>`, "a => b ==> c ===> d",
		[]string{"a", "=>", "b", "==>", "c", "===>", "d"}},
}

func TestDelimiterSplit(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsDelimiterSplit {
		d.UseRegexp = test.useRegexp
		if err := d.Set(test.expr); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.expr, err)
		}
		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q).Split(%q) = %q; want %q",
				d.re, test.src, actual, expect)
		}
	}
}
