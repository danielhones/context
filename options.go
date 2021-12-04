package main

type Language string

const (
	PYTHON Language = "py"
	RUBY   Language = "rb"
)

// This struct stores information about options for processing files.  It's primarily
// used as an argument to the processFile function
type Options struct {
	PrintNums bool // whether to print line numbers in the output
	Language  string
}

func (o Options) AutoDetect() bool {
	return o.Language == ""
}
