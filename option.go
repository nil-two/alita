package main

import (
	"flag"
)

type Option struct {
	Delimiter string
	UseRegexp bool
	Count     int
	Margin    string
	Justify   string
	IsHelp    bool
	IsVersion bool
	Files     []string
}

func ParseOption(args []string) (*Option, error) {
	opt := &Option{}

	f := flag.NewFlagSet("alita", flag.ContinueOnError)
	f.StringVar(&opt.Delimiter, "d", "", "")
	f.StringVar(&opt.Delimiter, "delimiter", "", "")
	f.BoolVar(&opt.UseRegexp, "r", false, "")
	f.BoolVar(&opt.UseRegexp, "regexp", false, "")
	f.IntVar(&opt.Count, "c", 0, "")
	f.IntVar(&opt.Count, "count", 0, "")
	f.StringVar(&opt.Margin, "m", "", "")
	f.StringVar(&opt.Margin, "margin", "", "")
	f.StringVar(&opt.Justify, "j", "", "")
	f.StringVar(&opt.Justify, "justify", "", "")
	f.BoolVar(&opt.IsHelp, "h", false, "")
	f.BoolVar(&opt.IsHelp, "help", false, "")
	f.BoolVar(&opt.IsVersion, "version", false, "")

	if err := f.Parse(args); err != nil {
		return nil, err
	}
	opt.Files = f.Args()
	return opt, nil
}
