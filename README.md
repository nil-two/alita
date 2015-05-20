alita
===
Align FILE(s), or standard input.

It's inspired by [h1mesuke/vim-alignta](https://github.com/h1mesuke/vim-alignta).

	$ cat user.conf
	[user]
	name=Tom
	age=17
	userid=10010

	$ cat user.conf | alita -d==
	[user]
	name   = Tom
	age    = 17
	userid = 10010


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
`alita` can be easily installed as an executable.
Download the latest
[compiled binaries](https://github.com/kusabashira/alita/releases)
and put it anywhere in your executable path.

Or, if you've done Go development before
and your $GOPATH/bin directory is already in your PATH:

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

DELIM will interpreted as fixed string.

	$ alita --delimiter==
	delimit line with '='

	$ alita --delimiter="->"
	delimit line with '->'


If you enabled regexp with `-r` or `--regexp`.
DELIM will interpreted as regexp.

	$ alita --regexp --delimiter="=+>"
	delimit line with /=+>/

	$ alita --regexp --delimiter="[:/]+"
	delimit line with /[:\/]+/


###-m, --margin=FORMAT
Join cells with margin which described in FORMAT.
Default FORMAT is `1:1`.

FORMAT is `colon separated digits` or `digit only`.

If FORMAT is `colon separated digits`,
left side will interpreted as `left-margin`
right side will interpreted as `right-margin`

	$ alita --margin=0:1
	left-margin will 0 space
	left-margin will 1 space

	$ alita --margin=3:2
	left-margin will 3 space
	left-margin will 2 space

If FORMAT is `digit only`,
digit will interpreted as both `left-margin` and `right-margin`

	$ alita --margin=2
	left-margin will 2 space
	left-margin will 2 space

	$ alita --margin=0
	left-margin will 0 space
	left-margin will 0 space

### -j, --justfy=SEQUENCE
Justfy cells by format which described in SEQUENCE.
SEQUENCE include `l`, `r` or `c`.

| char | justfy         |
|:-----|:---------------|
| l    | left-justfy    |
| r    | right-justfy   |
| c    | center-justfy  |

SEQUENCE will interpreted as following format.
(Default SEQUENCE is `l`)

`{L-fld-align} [ {M-fld-align} {R-fld-align} ]...`

You can specify any number of `{M-fld-align}` and `{R-fld-align}`.
If the match there is more than one,
they will continue to be applied to the order .
(next of last `R-fld-align` is first `M-fld-align`)

	$ alita --justfy=r
	all cells right-justified

	$ alita --justfy=lc
	Only the first cell right-justified,
	the rest of the cell center-justified.

	$ alita --justfy=rlc
	cell[0] right-justified.  cell[1] left-justified.
	cell[2] center-justified. cell[3] left-justified.
	cell[4] center-justified. ...

License
--------
MIT License

Author
-------
wara <kusabashira227@gmail.com>
