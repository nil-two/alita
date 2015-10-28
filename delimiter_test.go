package main

import (
	"reflect"
	"regexp"
	"testing"
)

var delimiterSetTests = []struct {
	useRegexp bool
	expr      string
	re        *regexp.Regexp
}{
	// fixed
	{false, `=`, regexp.MustCompile(`=`)},
	{false, `=+`, regexp.MustCompile(`=\+`)},
	{false, `-*>`, regexp.MustCompile(`-\*>`)},
	{false, `abc`, regexp.MustCompile(`abc`)},
	{false, `\w+:`, regexp.MustCompile(`\\w\+:`)},
	{false, `[:/]+`, regexp.MustCompile(`\[:/\]\+`)},

	// regexp
	{true, `=`, regexp.MustCompile(`=`)},
	{true, `=+`, regexp.MustCompile(`=+`)},
	{true, `-*>`, regexp.MustCompile(`-*>`)},
	{true, `abc`, regexp.MustCompile(`abc`)},
	{true, `\w+:`, regexp.MustCompile(`\w+:`)},
	{true, `[:/]+`, regexp.MustCompile(`[:/]+`)},
}

func TestDelimiterSet(t *testing.T) {
	for _, test := range delimiterSetTests {
		d, err := NewDelimiter(test.expr, test.useRegexp, -1)
		if err != nil {
			t.Errorf("NewDelimiter(%q, %v, %v) returns %q; want nil",
				test.expr, test.useRegexp, -1, err)
			continue
		}
		if !reflect.DeepEqual(d.re, test.re) {
			t.Errorf("NewDelimiter(%q, %v, %v).re got %q; want %q",
				test.expr, test.useRegexp, -1, d.re, test.re)
		}
	}
}

var delimiterSplitDefaultTests = []struct {
	src string
	dst []string
}{
	// normal
	{"a", []string{"a"}},
	{"a b", []string{"a", "b"}},
	{"a b c", []string{"a", "b", "c"}},
	{"a b c d", []string{"a", "b", "c", "d"}},
	{"ab cd", []string{"ab", "cd"}},
	{"日本 語", []string{"日本", "語"}},

	// long spaces
	{"a  b", []string{"a", "b"}},
	{"a  b c", []string{"a", "b", "c"}},
	{"a \t b c", []string{"a", "b", "c"}},

	// head and tail spaces
	{"  a b c", []string{"", "a", "b", "c"}},
	{"a b c ", []string{"a", "b", "c", ""}},
}

func TestDelimiterDefaultSplit(t *testing.T) {
	d, _ := NewDelimiter("", false, 0)
	for _, test := range delimiterSplitDefaultTests {
		expect := test.dst
		actual := d.Split(test.src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("NewDelimiter(%q, %v, %v).Split(%q) = %q; want %q",
				"", false, 0, test.src, actual, expect)
		}
	}
}

var delimiterSplitTests = []struct {
	useRegexp bool
	expr      string
	src       string
	dst       []string
}{
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
	{false, `=>`, "a==>=b==>=c", []string{"a=", "=>", "=b=", "=>", "=c"}},

	{true, `=+>`, "a => b",
		[]string{"a", "=>", "b"}},
	{true, `=+>`, "a => b ==> c ===> d",
		[]string{"a", "=>", "b", "==>", "c", "===>", "d"}},
	{true, `=+>`, "a=>b==>=c",
		[]string{"a", "=>", "b", "==>", "=c"}},
}

func TestDelimiterSplit(t *testing.T) {
	for _, test := range delimiterSplitTests {
		d, err := NewDelimiter(test.expr, test.useRegexp, -1)
		if err != nil {
			t.Errorf("NewDelimiter(%q, %v, %v) returns %q; want nil",
				test.expr, test.useRegexp, -1, err)
			continue
		}

		expect := test.dst
		actual := d.Split(test.src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("NewDelimiter(%q, %v, %v).Split(%q) = %q; want %q",
				test.expr, test.useRegexp, -1, test.src, actual, expect)
		}
	}
}

var delimiterSplitWithCountTests = []struct {
	count int
	src   string
	dst   []string
}{
	// less than 1
	{-2, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
	{-1, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
	{0, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},

	// greater than 0
	{1, "n =  m   =    100", []string{"n", "=  m   =    100"}},
	{2, "n =  m   =    100", []string{"n", "=", "m   =    100"}},
	{3, "n =  m   =    100", []string{"n", "=", "m", "=    100"}},
	{4, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
	{5, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
}

func TestSplitWithCount(t *testing.T) {
	delim := "="
	for _, test := range delimiterSplitWithCountTests {
		d, err := NewDelimiter(delim, false, test.count)
		if err != nil {
			t.Errorf("NewDelimiter(%q, %v, %v) returns %q; want nil",
				delim, false, test.count, err)
			continue
		}

		expect := test.dst
		actual := d.Split(test.src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("NewDelimiter(%q, %v, %v).Split(%q) = %q; want %q",
				delim, false, test.count, test.src, actual, expect)
		}
	}
}
