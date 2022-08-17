package helper

import 	"regexp"

type RegexStat struct {
	Before, Match, After string
	Start, Length        int
}

//─────────────┤ MapFromPatSlice ├─────────────

func MapFromPatSlice(pats []string) map[string]*regexp.Regexp {
	var regs = make(map[string]*regexp.Regexp)
	for _, p := range pats {
		regs[p] = regexp.MustCompile(p)
	}

	return regs
}

//─────────────┤ RegexStatFromPat ├─────────────

func RegexStatFromPat(pat, search string, rx ...*regexp.Regexp) RegexStat {
	if len(search) == 0 || len(pat) == 0 {
		return RegexStat{}
	}

	var r *regexp.Regexp
	if len(rx) == 0 {
		r = regexp.MustCompile(pat)
	} else {
		r = rx[0]
	}

	loc := r.FindStringIndex(search)
	if loc == nil {
		return RegexStat{Before: search}
	}

	return RegexStat{
		Before: search[:loc[0]],
		Match:  search[loc[0]:loc[1]],
		After:  search[loc[1]:],
		Start:  loc[0],
		Length: loc[1] - loc[0],
	}
}
