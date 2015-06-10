package alita

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
	// Fixed
	{false, `=`, regexp.MustCompile(`=`)},
	{false, `=+`, regexp.MustCompile(`=\+`)},
	{false, `-*>`, regexp.MustCompile(`-\*>`)},
	{false, `abc`, regexp.MustCompile(`abc`)},
	{false, `\w+:`, regexp.MustCompile(`\\w\+:`)},
	{false, `[:/]+`, regexp.MustCompile(`\[:/\]\+`)},

	// Regexp
	{true, `=`, regexp.MustCompile(`=`)},
	{true, `=+`, regexp.MustCompile(`=+`)},
	{true, `-*>`, regexp.MustCompile(`-*>`)},
	{true, `abc`, regexp.MustCompile(`abc`)},
	{true, `\w+:`, regexp.MustCompile(`\w+:`)},
	{true, `[:/]+`, regexp.MustCompile(`[:/]+`)},
}

func TestDelimiterSet(t *testing.T) {
	d := NewDelimiter()
	for _, test := range delimiterSetTests {
		d.UseRegexp = test.useRegexp
		if err := d.Set(test.expr); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.expr, err)
		}
		if !reflect.DeepEqual(d.re, test.re) {
			t.Errorf("got %q; want %q",
				d.re, test.re)
		}
	}
}

var delimiterSplitDefaultTests = []struct {
	src string
	dst []string
}{
	// Normal
	{"a", []string{"a"}},
	{"a b", []string{"a", "b"}},
	{"a b c", []string{"a", "b", "c"}},
	{"a b c d", []string{"a", "b", "c", "d"}},
	{"ab cd", []string{"ab", "cd"}},
	{"日本 語", []string{"日本", "語"}},

	// Long spaces
	{"a  b", []string{"a", "b"}},
	{"a  b c", []string{"a", "b", "c"}},
	{"a \t b c", []string{"a", "b", "c"}},

	// Head and tail spaces
	{"  a b c", []string{"", "a", "b", "c"}},
	{"a b c ", []string{"a", "b", "c", ""}},
}

func TestDelimiterDefaultSplit(t *testing.T) {
	d := NewDelimiter()
	for _, test := range delimiterSplitDefaultTests {
		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q).Split(%q) = %q; want %q",
				d.re, test.src, actual, expect)
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
	d := NewDelimiter()
	for _, test := range delimiterSplitTests {
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
