package main

import (
	"github.com/jessevdk/go-flags"
)

type Option struct {
	Delimiter string `short:"d" long:"delimiter" default:""`
	UseRegexp bool   `short:"r" long:"regexp"    default:"false"`
	Count     int    `short:"c" long:"count"     default:"-1"`
	Margin    string `short:"m" long:"margin"    default:"1:1"`
	Justify   string `short:"j" long:"justify"   default:"l"`
	IsHelp    bool   `short:"h" long:"help"      default:"false"`
	IsVersion bool   `          long:"version"   default:"false"`
	Files     []string
}

func ParseOption(args []string) (*Option, error) {
	opt := &Option{}
	files, err := flags.ParseArgs(opt, args)
	if err != nil {
		return nil, err
	}
	opt.Files = files
	return opt, nil
}
