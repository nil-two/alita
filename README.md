alita
=====

[![Build Status](https://travis-ci.org/kusabashira/alita.svg?branch=master)](https://travis-ci.org/kusabashira/alita)

Align FILE(s), or standard input.

It's inspired by [h1mesuke/vim-alignta](https://github.com/h1mesuke/vim-alignta).

	$ cat user.conf
	[user]
	name=        Tom
	age  =17
	userid = 10001

	$ cat user.conf | alita -d=
	[user]
	name   = Tom
	age    = 17
	userid = 10001

Usage
-----

	$ alita [OPTION]... [FILE]...

	Delimiter control:
	  -d, --delimiter=DELIM    separate lines by DELIM
	  -r, --regexp             DELIM is a regular expression
	  -c, --count=COUNT        separate lines only COUNT times

	Output control:
	  -m, --margin=FORMAT      join cells by FORMAT
	  -j, --justify=SEQUENCE   justify cells by SEQUENCE

	Miscellaneous:
	  -h, --help               display this help and exit
	      --version            display version information and exit

Installation
------------

### compiled binary

See [releases](https://github.com/kusabashira/alita/releases)

### go get

	go get github.com/kusabashira/alita

Options
-------

### -h, --help

Display a help message.

### --version

Output the version of alita.

### -d, --delimiter=DELIM

Separate lines by DELIM.
Default DELIM is `spaces (/\s+/)`.

	$ cat nums.txt
	1 100 10000
	100 10000 1
	10000 1 100

	$ cat nums.txt | alita
	1     100   10000
	100   10000 1
	10000 1     100

DELIM will interpreted as a fixed string.

	$ cat user
	name=Tom
	age=17

	$ cat user | alita -d=
	(delimit line by '=')
	name = Tom
	age  = 17

	$ cat snip.cpp
	cout    <<    "9 * 2 = "<<9 * 2 << endl;
	cout << "9 / 2 = "<<9 / 2 << ".." << 9 % 2<< endl;

	$ cat snip.cpp | alita -d'<<'
	(delimit line by '<<')
	cout << "9 * 2 = " << 9 * 2 << endl;
	cout << "9 / 2 = " << 9 / 2 << ".."  << 9 % 2 << endl;

### -r, --regexp

Separate lines by a regular expression.

	$ cat root
	a=>b ===>  c
	c ==>    d ==>e
	f===> g =>   h

	$ cat root | alita -rd'=+>'
	(delimit line by /=+>/)
	a =>   b ===> c
	c ==>  d ==>  e
	f ===> g =>   h

	$ cat url
	https://github.com/vim-scripts/Align
	https://github.com/h1mesuke/vim-alignta
	https://github.com/kusabashira/alita

	$ cat url | alita -rd[:/]+
	(delimit line by /[:\/]+/)
	https :// github.com / vim-scripts / Align
	https :// github.com / h1mesuke    / vim-alignta
	https :// github.com / kusabashira / alita

### -c, --count=COUNT

Separate lines only COUNT times.
Default COUNT is `-1`.

If COUNT is smaller than 0, delimiter separates as much as possible.

	$ cat graph
	1]
	10]]]]]]]]]]
	3]]]
	7]]]]]]]

	$ cat graph | alita -d]
	1  ]
	10 ]  ]  ]  ]  ]  ]  ]  ]  ]  ]
	3  ]  ]  ]
	7  ]  ]  ]  ]  ]  ]  ]

	$ cat graph | alita -d] -c1
	1  ]
	10 ]]]]]]]]]]
	3  ]]]
	7  ]]]]]]]

### -m, --margin=FORMAT

Join cells with a margin which described in FORMAT.
Default FORMAT is `1:1`.

FORMAT needs to be `{left-margin}:{right-margin}` or `{margin}`.

If FORMAT is a `colon separated digits`.
The left side will interpreted as `left-margin`,
and the right side will interpreted as `right-margin`.

	$ cat user | alita -d= -m0:1
	(left-margin: 0, right-margin: 1)
	name= Tom
	age = 17

	$ cat user | alita -d= -m3:2
	(left-margin: 3, right-margin: 2)
	name   =  Tom
	age    =  17

If FORMAT is a `digit only`,
The digit will interpreted as both `left-margin` and `right-margin`.

	$ cat user | alita -m2
	(left-margin: 2, right-margin: 2)
	name  =  Tom
	age   =  17

	$ cat user | alita -m0
	(left-margin: 0, right-margin: 0)
	name=Tom
	age =17

### -j, --justify=SEQUENCE

Justify cells by a format which described in SEQUENCE.

SEQUENCE includes only `l`, `r` and `c`.

| char | justify        |
|:-----|:---------------|
| l    | left-justify   |
| r    | right-justify  |
| c    | center-justify |

SEQUENCE will interpreted as the following format.

(Default SEQUENCE is `l`)

`{L-fld-align} [ {M-fld-align} {R-fld-align} ]...`

You can specify any number of `{M-fld-align}` and `{R-fld-align}`.

Justifies are applied to cells in order from left.

(The next of the last `R-fld-align` is the first `M-fld-align`)

	$ cat text
	a = bbbbb =  c = ddddd =  e = fffff =  1
	 aaa = bbb = ccc =  ddd = eee = fff = 10
	aaaaa =  b = ccccc = d = eeeee =  f = 100
	
	$ cat text | alita -d=
	(all cells left-justified)
	a     = bbbbb = c     = ddddd = e     = fffff = 1
	aaa   = bbb   = ccc   = ddd   = eee   = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

	$ cat text | alita -d= -jr
	(all cells right-justified)
	    a = bbbbb =     c = ddddd =     e = fffff =   1
	  aaa =   bbb =   ccc =   ddd =   eee =   fff =  10
	aaaaa =     b = ccccc =     d = eeeee =     f = 100

	$ cat text | alita -d= -jrl
	(Only the first cell right-justified,
	the rest of the cell left-justified)
	    a = bbbbb = c     = ddddd = e     = fffff = 1
	  aaa = bbb   = ccc   = ddd   = eee   = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

	$ cat text | alita -jrllcc
	(cell[0] right-justified.  cell[1] left-justified.
	 cell[2] left-justified.   cell[3] center-justified.
	 cell[4] center-justified. cell[5] left-justified.
	 cell[6] left-justified.   cell[7] center-justified.
	 cell[8] center-justified. ...)
	    a = bbbbb =   c   = ddddd =   e   = fffff =  1
	  aaa = bbb   =  ccc  = ddd   =  eee  = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

Other Specification
-------------------

#### head spaces

If the input text includes head spaces.
alita leaves the shortest head spaces.

#### trailing spaces

alita removes all trailing spaces.

License
-------

MIT License

Author
-------

kusabashira <kusabashira227@gmail.com>
