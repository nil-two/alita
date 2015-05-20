package alita

import (
	"bytes"
	"reflect"
	"testing"
)

type AlignSimpleTest struct {
	src []byte
	dst []byte
}

var indexTestsAlignSimple = []AlignSimpleTest{
	{[]byte(`
a
bb
ccc
`[1:]), []byte(`
a
bb
ccc
`[1:])},

	{[]byte(`
1 100
10 10
100 1
`[1:]), []byte(`
1   100
10  10
100 1
`[1:])},

	{[]byte(`
1 10 100
10 100 1
100 1 10
`[1:]), []byte(`
1   10  100
10  100 1
100 1   10
`[1:])},

	{[]byte(`
1 10 100 1000
10 100 1000 1
100 1000 1 10
1000 1 10 100
`[1:]), []byte(`
1    10   100  1000
10   100  1000 1
100  1000 1    10
1000 1    10   100
`[1:])},
}

func TestAlignSimple(t *testing.T) {
	for _, test := range indexTestsAlignSimple {
		w := bytes.NewBuffer(make([]byte, 0))
		a := NewAligner(w)

		r := bytes.NewReader(test.src)
		if err := a.ReadAll(r); err != nil {
			t.Errorf("ReadAll(%q) returns err; want nil", test.src)
		}
		if err := a.Flush(); err != nil {
			t.Errorf("Flush(%q) returns err; want nil", test.src)
		}

		actual := w.Bytes()
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("got %q;want %q", actual, expect)
		}
	}
}
