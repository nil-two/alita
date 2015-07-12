package main

import (
	"github.com/jessevdk/go-flags"
)

type Option struct {
	Delimiter string `short:"d" long:"delimiter"`
	UseRegexp bool   `short:"r" long:"regexp"`
	Count     int    `short:"c" long:"count"`
	Margin    string `short:"m" long:"margin"`
	Justify   string `short:"j" long:"justify"`
	IsHelp    bool   `short:"h" long:"help"`
	IsVersion bool   `          long:"version"`
	Files     []string
}

func ParseOption(args []string) (*Option, error) {
	opt := &Option{}
	f := flags.NewParser(opt, flags.PassDoubleDash)

	files, err := f.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	opt.Files = files

	return opt, nil
}
