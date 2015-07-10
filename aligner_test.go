package main

import (
	"bytes"
	"reflect"
	"testing"
)

func testAlign(t *testing.T, a *Aligner, src, dst []byte) {
	out := bytes.NewBuffer(make([]byte, 0))

	r := bytes.NewReader(src)
	if err := a.ReadAll(r); err != nil {
		t.Errorf("ReadAll(%q) returns err; want nil", src)
	}
	if err := a.Flush(out); err != nil {
		t.Errorf("Flush(%q) returns err; want nil", src)
	}

	actual := out.Bytes()
	expect := dst
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got:\n%swant:\n%s", actual, expect)
	}
}

var alignSimpleTests = []struct {
	src []byte
	dst []byte
}{
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
	for _, test := range alignSimpleTests {
		a := NewAligner()
		testAlign(t, a, test.src, test.dst)
	}
}

var alignFixedTests = []struct {
	delim string
	src   []byte
	dst   []byte
}{
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

	{`\d\+`, []byte(`
a \d\+ 1 aaaaa
bbb \d\+ 10 bbb
ccccc \d\+ 100 c
`[1:]), []byte(`
a     \d\+ 1 aaaaa
bbb   \d\+ 10 bbb
ccccc \d\+ 100 c
`[1:])},

	{`:`, []byte(`
one:two:three
four:five:six
seven:eight:nine
`[1:]), []byte(`
one   : two   : three
four  : five  : six
seven : eight : nine
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

	{`=`, []byte(`
a =  1
 bbb = 10
ccccc = 100
a =  2
 bbb = 20
ccccc = 200
a =  3
 bbb = 30
ccccc = 300
a =  4
 bbb = 40
ccccc = 400
`[1:]), []byte(`
a     = 1
bbb   = 10
ccccc = 100
a     = 2
bbb   = 20
ccccc = 200
a     = 3
bbb   = 30
ccccc = 300
a     = 4
bbb   = 40
ccccc = 400
`[1:])},

	{`=`, []byte(`
a = 1

bbb = 10

ccccc = 100
`[1:]), []byte(`
a     = 1

bbb   = 10

ccccc = 100
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
}

func TestAlignFixed(t *testing.T) {
	for _, test := range alignFixedTests {
		d := NewDelimiter()
		if err := d.Set(test.delim); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		a := NewAlignerWithModules(d, nil, nil, nil)

		testAlign(t, a, test.src, test.dst)
	}
}

var alignRegexpTests = []struct {
	delim string
	src   []byte
	dst   []byte
}{
	{`=+>`, []byte(`
a=>b ===>  c
c ==>    d ==>e
f===> g =>   h
`[1:]), []byte(`
a =>   b ===> c
c ==>  d ==>  e
f ===> g =>   h
`[1:])},

	{`[:/]+`, []byte(`
https://github.com/vim-scripts/Align
https://github.com/h1mesuke/vim-alignta
https://github.com/kusabashira/alita
`[1:]), []byte(`
https :// github.com / vim-scripts / Align
https :// github.com / h1mesuke    / vim-alignta
https :// github.com / kusabashira / alita
`[1:])},

	{`\\$`, []byte(`
one \
two three \
four five six \
seven \\ \
eight \nine \
ten \
`[1:]), []byte(`
one           \
two three     \
four five six \
seven \\      \
eight \nine   \
ten           \
`[1:])},

	{`(&|\\\\)`, []byte(`
one&two&three\\ \hline
four&five&six\\
seven&eight&nine\\
`[1:]), []byte(`
one   & two   & three \\ \hline
four  & five  & six   \\
seven & eight & nine  \\
`[1:])},

	{`\d+`, []byte(`
a \d\+ 1 aaaaa
bbb \d\+ 10 bbb
ccccc \d\+ 100 c
`[1:]), []byte(`
a \d\+     1   aaaaa
bbb \d\+   10  bbb
ccccc \d\+ 100 c
`[1:])},

	{`\d+`, []byte(`
a \d\+ 1 \u\+ AAAAA a
bbb \d\+ 10 \u\+ BBB b
ccccc \d\+ 100 \u\+ C c
`[1:]), []byte(`
a \d\+     1   \u\+ AAAAA a
bbb \d\+   10  \u\+ BBB b
ccccc \d\+ 100 \u\+ C c
`[1:])},

	{`(=|/\*|\*/)`, []byte(`
a =  1 /* AAAAA */
 bbb = 10 /*  BBB */
ccccc = 100 /* C */
`[1:]), []byte(`
a     = 1   /* AAAAA */
bbb   = 10  /* BBB   */
ccccc = 100 /* C     */
`[1:])},
}

func TestAlignRegexp(t *testing.T) {
	for _, test := range alignRegexpTests {
		d := NewDelimiter()
		d.UseRegexp = true
		if err := d.Set(test.delim); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		a := NewAlignerWithModules(d, nil, nil, nil)

		testAlign(t, a, test.src, test.dst)
	}
}

var alignMarginTests = []struct {
	margin string
	delim  string
	src    []byte
	dst    []byte
}{
	{`1:0`, `=`, []byte(`
`[1:]), []byte(`
`[1:])},

	{`0:1`, `=`, []byte(`
name=Tom
age=17
`[1:]), []byte(`
name= Tom
age = 17
`[1:])},

	{`3:2`, `=`, []byte(`
name=Tom
age=17
`[1:]), []byte(`
name   =  Tom
age    =  17
`[1:])},

	{`2`, `=`, []byte(`
name=Tom
age=17
`[1:]), []byte(`
name  =  Tom
age   =  17
`[1:])},

	{`0`, `=`, []byte(`
name=Tom
age=17
`[1:]), []byte(`
name=Tom
age =17
`[1:])},

	{`1:3`, `=`, []byte(`
a=bbb=ccccc
aaa=b=ccc
`[1:]), []byte(`
a   =   bbb =   ccccc
aaa =   b   =   ccc
`[1:])},

	{`0:2`, `=`, []byte(`
a=bbb=ccccc
aaa=b=ccc
`[1:]), []byte(`
a  =  bbb=  ccccc
aaa=  b  =  ccc
`[1:])},

	{`0`, `|`, []byte(`
|one|two|three|
|four|five|six|
|seven|eight|nine|
`[1:]), []byte(`
|one  |two  |three|
|four |five |six  |
|seven|eight|nine |
`[1:])},

	{`2`, `|`, []byte(`
|one|two|three|
|four|five|six|
|seven|eight|nine|
`[1:]), []byte(`
  |  one    |  two    |  three  |
  |  four   |  five   |  six    |
  |  seven  |  eight  |  nine   |
`[1:])},
}

func TestAlignMargin(t *testing.T) {
	for _, test := range alignMarginTests {
		d := NewDelimiter()
		if err := d.Set(test.delim); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		m := NewMargin()
		if err := m.Set(test.margin); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		a := NewAlignerWithModules(d, nil, m, nil)

		testAlign(t, a, test.src, test.dst)
	}
}

var alignJustifyTests = []struct {
	justfy string
	delim  string
	src    []byte
	dst    []byte
}{
	{`l`, `=`, []byte(`
a = 1
bbb = 10
ccccc = 100
`[1:]), []byte(`
a     = 1
bbb   = 10
ccccc = 100
`[1:])},

	{`r`, `=`, []byte(`
a = 1
bbb = 10
ccccc = 100
`[1:]), []byte(`
    a =   1
  bbb =  10
ccccc = 100
`[1:])},

	{`l`, `=`, []byte(`
a = bbbbb =  c = ddddd =  e = fffff =  1
 aaa = bbb = ccc =  ddd = eee = fff = 10
aaaaa =  b = ccccc = d = eeeee =  f = 100
`[1:]), []byte(`
a     = bbbbb = c     = ddddd = e     = fffff = 1
aaa   = bbb   = ccc   = ddd   = eee   = fff   = 10
aaaaa = b     = ccccc = d     = eeeee = f     = 100
`[1:])},

	{`r`, `=`, []byte(`
a = bbbbb =  c = ddddd =  e = fffff =  1
 aaa = bbb = ccc =  ddd = eee = fff = 10
aaaaa =  b = ccccc = d = eeeee =  f = 100
`[1:]), []byte(`
    a = bbbbb =     c = ddddd =     e = fffff =   1
  aaa =   bbb =   ccc =   ddd =   eee =   fff =  10
aaaaa =     b = ccccc =     d = eeeee =     f = 100
`[1:])},

	{`rl`, `=`, []byte(`
a = bbbbb =  c = ddddd =  e = fffff =  1
 aaa = bbb = ccc =  ddd = eee = fff = 10
aaaaa =  b = ccccc = d = eeeee =  f = 100
`[1:]), []byte(`
    a = bbbbb = c     = ddddd = e     = fffff = 1
  aaa = bbb   = ccc   = ddd   = eee   = fff   = 10
aaaaa = b     = ccccc = d     = eeeee = f     = 100
`[1:])},

	{`rllcc`, `=`, []byte(`
a = bbbbb =  c = ddddd =  e = fffff =  1
 aaa = bbb = ccc =  ddd = eee = fff = 10
aaaaa =  b = ccccc = d = eeeee =  f = 100
`[1:]), []byte(`
    a = bbbbb =   c   = ddddd =   e   = fffff =  1
  aaa = bbb   =  ccc  = ddd   =  eee  = fff   = 10
aaaaa = b     = ccccc = d     = eeeee = f     = 100
`[1:])},

	{`r`, `＝`, []byte(`
あ ＝  壱
 あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
        あ ＝     壱
    あいう ＝   壱十
あいうえお ＝ 壱十百
`[1:])},

	{`c`, `＝`, []byte(`
あ ＝  壱
 あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
    あ     ＝   壱
  あいう   ＝  壱十
あいうえお ＝ 壱十百
`[1:])},

	{`lcr`, `＝`, []byte(`
あ ＝  壱
 あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
あ         ＝     壱
あいう     ＝   壱十
あいうえお ＝ 壱十百
`[1:])},

	{`rcl`, `＝`, []byte(`
あ ＝  壱
 あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:]), []byte(`
        あ ＝ 壱
    あいう ＝ 壱十
あいうえお ＝ 壱十百
`[1:])},
}

func TestAlignJustify(t *testing.T) {
	for _, test := range alignJustifyTests {
		d := NewDelimiter()
		if err := d.Set(test.delim); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		p := NewPadding()
		if err := p.Set(test.justfy); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.delim, err)
		}
		a := NewAlignerWithModules(d, p, nil, nil)

		testAlign(t, a, test.src, test.dst)
	}
}
