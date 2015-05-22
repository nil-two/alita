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
`[1:]), []byte(`
`[1:])},

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

	{[]byte(`
(x == ABC) ? "abc" :
(x == DEFG) ? "defg" :
(x == HIJKL) ? "hijkl" : "???");
`[1:]), []byte(`
(x == ABC)   ? "abc"   :
(x == DEFG)  ? "defg"  :
(x == HIJKL) ? "hijkl" : "???");
`[1:])},

	{[]byte(`
one two three four five
six seven eight nine ten
eleven twelve thirteen fourteen fifteen
`[1:]), []byte(`
one    two    three    four     five
six    seven  eight    nine     ten
eleven twelve thirteen fourteen fifteen
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
			t.Errorf("got %q; want %q", actual, expect)
		}
	}
}

type AlignWithFixedTest struct {
	delim string
	src   []byte
	dst   []byte
}

var indexTestsAlignFixedTest = []AlignWithFixedTest{
	{`=`, []byte(`
a =  1
 bbb = 10
ccccc = 100
`[1:]), []byte(`
a     = 1
bbb   = 10
ccccc = 100
`[1:])},

	{`=`, []byte(`
[user]
name=        Tom
age  =17
userid = 10001
`[1:]), []byte(`
[user]
name   = Tom
age    = 17
userid = 10001
`[1:])},

	{`=`, []byte(`
name=Tom
age=17
`[1:]), []byte(`
name = Tom
age  = 17
`[1:])},

	{`<<`, []byte(`
cout    <<    "9 * 2 = "<<9 * 2 << endl;
cout << "9 / 2 = "<<9 / 2 << ".." << 9 % 2<< endl;
`[1:]), []byte(`
cout << "9 * 2 = " << 9 * 2 << endl;
cout << "9 / 2 = " << 9 / 2 << ".."  << 9 % 2 << endl;
`[1:])},

	{`<<`, []byte(`
cin << x;
cin      << y;
cout << "this is x=" << x;
cout << "but y=" << y << "is not";
`[1:]), []byte(`
cin  << x;
cin  << y;
cout << "this is x=" << x;
cout << "but y="     << y  << "is not";
`[1:])},

	{`:=`, []byte(`
aa:=bb:=cc:=1;
a:=b:=c:=1;
aaa:=bbb:=ccc:=1;
`[1:]), []byte(`
aa  := bb  := cc  := 1;
a   := b   := c   := 1;
aaa := bbb := ccc := 1;
`[1:])},

	{`=`, []byte(`
aa=bb=cc=1;
a=b=c=1;
aaa=bbb=ccc=1;
`[1:]), []byte(`
aa  = bb  = cc  = 1;
a   = b   = c   = 1;
aaa = bbb = ccc = 1;
`[1:])},
	{`＝`, []byte(`
あ ＝  壱
あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
あ         ＝ 壱
あいう     ＝ 壱十
あいうえお ＝ 壱十百
`[1:])},

	{`＝`, []byte(`
あ ＝ 壱
ΑΒΓ ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
あ         ＝ 壱
ΑΒΓ     ＝ 壱十
あいうえお ＝ 壱十百
`[1:])},
}

func TestAlignFixed(t *testing.T) {
	for _, test := range indexTestsAlignFixedTest {
		w := bytes.NewBuffer(make([]byte, 0))
		a := NewAligner(w)
		if err := a.Delimiter.Set(test.delim); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}

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
			t.Errorf("got %q; want %q", actual, expect)
		}
	}
}
