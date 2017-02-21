package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ogier/pflag"
)

var (
	cmdName    = "alita"
	cmdVersion = "0.7.1"
)

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	delimiter string
	useRegexp bool
	count     int
	margin    string
	justify   string
	isHelp    bool
	isVersion bool
}

func NewCLI(stdin io.Reader, stdout io.Writer, stderr io.Writer) *CLI {
	return &CLI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *CLI) printUsage() {
	fmt.Fprintf(c.stderr, `
Usage: %s [OPTION]... [FILE]...
Align FILE(s), or standard input.

Delimiter control:
  -d, --delimiter=DELIM      separate lines by DELIM
  -r, --regexp               DELIM is a regular expression
  -c, --count=COUNT          separate lines only COUNT times

Output control:
  -m, --margin=N[:M]         put N or N and M spaces at both ends of DELIM
  -j, --justify=[l|c|r]...   justify cells to the left, center, or right

Miscellaneous:
  -h, --help                 display this help and exit
      --version              display version information and exit
`[1:], cmdName)
}

func (c *CLI) printVersion() {
	fmt.Fprintf(c.stderr, "%s\n", cmdVersion)
}

func (c *CLI) printErr(err interface{}) {
	fmt.Fprintf(c.stderr, "%s: %s\n", cmdName, err)
}

func (c *CLI) guideToHelp() {
	fmt.Fprintf(c.stderr, "Try '%s --help' for more information.\n", cmdName)
}

func (c *CLI) parseOption(args []string) (argFiles []string, err error) {
	f := pflag.NewFlagSet(cmdName, pflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	f.StringVarP(&c.delimiter, "delimiter", "d", "", "")
	f.BoolVarP(&c.useRegexp, "regexp", "r", false, "")
	f.IntVarP(&c.count, "count", "c", -1, "")
	f.StringVarP(&c.margin, "margin", "m", "", "")
	f.StringVarP(&c.justify, "justify", "j", "", "")
	f.BoolVarP(&c.isHelp, "help", "h", false, "")
	f.BoolVarP(&c.isVersion, "version", "", false, "")

	if err := f.Parse(args); err != nil {
		return nil, err
	}
	return f.Args(), nil
}

func (c *CLI) newAligner() (a *Aligner, err error) {
	return NewAligner(&Option{
		Delimiter: c.delimiter,
		UseRegexp: c.useRegexp,
		Count:     c.count,
		Margin:    c.margin,
		Justify:   c.justify,
	})
}

func (c *CLI) newArgf(argFiles []string) (r io.Reader, err error) {
	switch len(argFiles) {
	case 0:
		return c.stdin, nil
	default:
		rs := make([]io.Reader, len(argFiles))
		for i, argFile := range argFiles {
			f, err := os.Open(argFile)
			if err != nil {
				return nil, err
			}
			rs[i] = f
		}
		return io.MultiReader(rs...), nil
	}
}

func (c *CLI) do(a *Aligner, r io.Reader) error {
	if err := a.ReadAll(r); err != nil {
		return err
	}
	if err := a.Flush(c.stdout); err != nil {
		return err
	}
	return nil
}

func (c *CLI) Run(args []string) int {
	f, err := c.parseOption(args)
	if err != nil {
		c.printErr(err)
		c.guideToHelp()
		return 2
	}
	if c.isHelp {
		c.printUsage()
		return 0
	}
	if c.isVersion {
		c.printVersion()
		return 0
	}

	a, err := c.newAligner()
	if err != nil {
		c.printErr(err)
		c.guideToHelp()
		return 2
	}
	r, err := c.newArgf(f)
	if err != nil {
		c.printErr(err)
		c.guideToHelp()
		return 2
	}

	if err = c.do(a, r); err != nil {
		c.printErr(err)
		return 1
	}
	return 0
}

func main() {
	c := NewCLI(os.Stdin, os.Stdout, os.Stderr)
	e := c.Run(os.Args[1:])
	os.Exit(e)
}
