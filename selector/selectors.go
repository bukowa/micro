package selector

import "strings"

const (
	XRegexpMatch = "string_regexp"

	XStringPrefix    = "string_prefix"
	XStringSuffix    = "string_suffix"
	XStringEqual     = "string_equal"
	XStringEqualFold = "string_equal_fold"
	XStringContains  = "string_contains"
)

var RegexpMatch = XRegexpSelector{
	XName: XRegexpMatch,
	Regxp: nil,
}

var StringPrefix = XStringSelector{
	XName: XStringPrefix,
	ScoreFunc: func(value string, match string) bool {
		return strings.HasPrefix(value, match)
	},
	Match: nil,
}

var StringSuffix = XStringSelector{
	XName: XStringSuffix,
	ScoreFunc: func(value string, match string) bool {
		return strings.HasSuffix(value, match)
	},
	Match: nil,
}

var StringEqual = XStringSelector{
	XName: XStringEqual,
	ScoreFunc: func(value string, match string) bool {
		return value == match
	},
	Match: nil,
}

var StringEqualFold = XStringSelector{
	XName: XStringEqualFold,
	ScoreFunc: func(value string, match string) bool {
		return strings.EqualFold(value, match)
	},
	Match: nil,
}

var StringContains = XStringSelector{
	XName: XStringContains,
	ScoreFunc: func(value string, match string) bool {
		return strings.Contains(value, match)
	},
	Match: nil,
}
