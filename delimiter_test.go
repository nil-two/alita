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

var delimiterSplitWithCountTests = []struct {
	count int
	src   string
	dst   []string
}{
	// less than 0
	{-2, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
	{-1, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},

	// equal 0
	{0, "n =  m   =    100", []string{"n =  m   =    100"}},

	// greater than 0
	{1, "n =  m   =    100", []string{"n", "=  m   =    100"}},
	{2, "n =  m   =    100", []string{"n", "=", "m   =    100"}},
	{3, "n =  m   =    100", []string{"n", "=", "m", "=    100"}},
	{4, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
	{5, "n =  m   =    100", []string{"n", "=", "m", "=", "100"}},
}

func TestSplitWithCount(t *testing.T) {
	d := NewDelimiter()
	if err := d.Set(`=`); err != nil {
		t.Errorf("Delimiter(\"=\").Split(%q) = %q; want %q", err)
	}
	for _, test := range delimiterSplitWithCountTests {
		d.Count = test.count

		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q, %d).Split(%q) = %q; want %q",
				d.re, d.Count, test.src, actual, expect)
		}
	}
}
