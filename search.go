package main

const MATCH_TYPE_REGEX = "regex"
const MATCH_TYPE_LINE = "line"

// This struct holds information about what to look for in the source code, and type
// of match it is
type Search struct {
	MatchType string
	// If a regex/string match, this is the value to search for.  If not a regex/string
	// match, then this value is empty and meaningless:
	Val string
	// If a line match, these are the line numbers to search for.  If not a line match,
	// then this value is empty and meaningless:
	ValInts []int
}

func (s Search) IsRegexMatch() bool {
	return s.MatchType == MATCH_TYPE_REGEX
}

func (s *Search) SetRegexMatch() {
	(*s).MatchType = MATCH_TYPE_REGEX
}

func NewSearch() Search {
	return Search{MatchType: MATCH_TYPE_LINE}
}
