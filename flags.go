package main

import "flag"

type pprofFlags struct {
	args []string
	*flag.FlagSet
}

func (pprofFlags) ExtraUsage() string { return "" }

func (pprofFlags) AddExtraUsage(eu string) {}

func (f pprofFlags) StringList(o, d, c string) *[]*string {
	return &[]*string{f.String(o, d, c)}
}

func (f pprofFlags) Parse(usage func()) []string {
	f.FlagSet.Usage = usage
	if err := f.FlagSet.Parse(f.args); err != nil {
		panic(err)
	}
	return f.FlagSet.Args()
}
