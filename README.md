alita
===
Align FILE(s), or standard input.

It's inspired by [h1mesuke/vim-alignta](https://github.com/h1mesuke/vim-alignta).

	$ cat user.conf
	[user]
	name=        Tom
	age  =17
	userid = 10001

	$ cat user.conf | alita -d==
	[user]
	name   = Tom
	age    = 17
	userid = 10001

Usage
------
	$ alita [OPTION]... [FILE]...

	Delimiter control:
	  -r, --regexp             DELIM is a regular expression
	  -d, --delimiter=DELIM    delimit line by DELIM

	Output control:
	  -m, --margin=FORMAT      join cells with FORMAT
	  -j, --justfy=SEQUENCE    justfy cells by SEQUENCE

	Miscellaneous:
	  -h, --help               show this help message
	      --version            print the version

Installation
--------
###compiled binary
See [releases](https://github.com/kusabashira/alita/releases)

###go get
	go get github.com/kusabashira/alita

Command Line Options
------
###-h, --help
Display a help message.

###--version
Display the version of alita.

###-r, --regexp
Enable delimit line with regexp.

###-d, --delimiter=DELIM
Delimit line by DELIM.
Default DELIM is `spaces (/\s+/)`.

	$ cat nums.txt
	1 100 10000
	100 10000 1
	10000 1 100

	$ cat nums.txt | alita
	1     100   10000
	100   10000 1
	10000 1     100

DELIM will interpreted as fixed string.

	$ cat user
	name=Tom
	age=17

	$ cat user | alita -d==
	(delimit line with '=')
	name = Tom
	age  = 17

	$ cat snip.cpp
	cout    <<    "9 * 2 = "<<9 * 2 << endl;
	cout << "9 / 2 = "<<9 / 2 << ".." << 9 % 2<< endl;

	$ alita -d="<<"
	(delimit line with '<<')
	cout << "9 * 2 = " << 9 * 2 << endl;
	cout << "9 / 2 = " << 9 / 2 << ".."  << 9 % 2 << endl;

If you enabled regexp with `-r` or `--regexp`.
DELIM will interpreted as regexp.

	$ cat root
	a=>b ===>  c
	c ==>    d ==>e
	f===> g =>   h

	$ cat root | alita -r -d="=+>"
	(delimit line with /=+>/)
	a =>   b ===> c
	c ==>  d ==>  e
	f ===> g =>   h

	$ cat url
	https://github.com/vim-scripts/Align
	https://github.com/h1mesuke/vim-alignta
	https://github.com/kusabashira/alita

	$ cat url | alita -r -d="[:/]+"
	(delimit line with /[:\/]+/)
	https :// github.com / vim-scripts / Align
	https :// github.com / h1mesuke    / vim-alignta
	https :// github.com / kusabashira / alita

###-m, --margin=FORMAT
Join cells with margin which described in FORMAT.
Default FORMAT is `1:1`.

FORMAT needs to be `{left-margin}:{right-margin}` or `{margin}`.

If FORMAT is colon separated digits.
left side will interpreted as `left-margin`,
right side will interpreted as `right-margin`.

	$ cat user | alita -d== -m=0:1
	(left-margin: 0, right-margin: 1)
	name= Tom
	age = 17

	$ cat user | alita -d== -m=3:2
	(left-margin: 3, right-margin: 2)
	name   =  Tom
	age    =  17

If FORMAT is `digit only`,
digit will interpreted as both `left-margin` and `right-margin`

	$ cat user | alita --margin=2
	(left-margin: 2, right-margin: 2)
	name  =  Tom
	age   =  17

	$ cat user | alita --margin=0
	(left-margin will 0 space)
	name=Tom
	age =17

### -j, --justfy=SEQUENCE
Justfy cells by format which described in SEQUENCE.

SEQUENCE include `l`, `r` or `c`.

| char | justfy        |
|:-----|:--------------|
| l    | left-justfy   |
| r    | right-justfy  |
| c    | center-justfy |

SEQUENCE will interpreted as following format.

(Default SEQUENCE is `l`)

`{L-fld-align} [ {M-fld-align} {R-fld-align} ]...`

You can specify any number of `{M-fld-align}` and `{R-fld-align}`.

If the match there is more than one,
they will continue to be applied to the order.
(next of last `R-fld-align` is first `M-fld-align`)

	$ cat text
	a = bbbbb =  c = ddddd =  e = fffff =  1
	 aaa = bbb = ccc =  ddd = eee = fff = 10
	aaaaa =  b = ccccc = d = eeeee =  f = 100
	
	$ cat text | alita -d==
	(all cells left-justified)
	a     = bbbbb = c     = ddddd = e     = fffff = 1
	aaa   = bbb   = ccc   = ddd   = eee   = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

	$ cat text | alita -d== -j=r
	(all cells right-justified)
	    a = bbbbb =     c = ddddd =     e = fffff =   1
	  aaa =   bbb =   ccc =   ddd =   eee =   fff =  10
	aaaaa =     b = ccccc =     d = eeeee =     f = 100

	$ cat text | alita -d== -j=rl
	(Only the first cell right-justified,
	the rest of the cell left-justified)
	    a = bbbbb = c     = ddddd = e     = fffff = 1
	  aaa = bbb   = ccc   = ddd   = eee   = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

	$ alita --justfy=rllcc
	(cell[0] right-justified.  cell[1] left-justified.
	 cell[2] left-justified.   cell[3] center-justified.
	 cell[4] center-justified. cell[5] left-justified.
	 cell[6] left-justified.   cell[7] center-justified.
	 cell[8] center-justified. ...)
	    a = bbbbb =   c   = ddddd =   e   = fffff =  1
	  aaa = bbb   =  ccc  = ddd   =  eee  = fff   = 10
	aaaaa = b     = ccccc = d     = eeeee = f     = 100

Other specification
------
###head space
If input text has head space.
It will remain shortest head space.

License
--------
MIT License

Author
-------
wara <kusabashira227@gmail.com>
