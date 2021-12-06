package main

import (
	"io"
	"os"
)

// This struct stores information about options for processing files.  It's primarily
// used as an argument to the processFile function
type Options struct {
	PrintNums bool      // whether to print line numbers in the output
	Language  string    // string indicating which language to use when parsing
	Colorize  bool      // whether to colorize the output
	Out       io.Writer // where to write the results
	Err       io.Writer // where to write errors
}

func (o Options) AutoDetect() bool {
	return o.Language == ""
}

func NewOptions() Options {
	return Options{Out: os.Stdout, Err: os.Stderr}
}
