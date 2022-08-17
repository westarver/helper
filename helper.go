package helper

import (
	"strings"
	"unicode"
)

//─────────────┤ GetLinesFromString  ├─────────────

func GetLinesFromString(s string) []string {
	return strings.Split(s, "\n")
}

//─────────────┤ CountLinesInString  ├─────────────

func CountLinesInString(s string) int {
	ss := strings.Split(s, "\n")
	return len(ss)
}

//─────────────┤ RemoveDupChar ├─────────────

func RemoveDupChar(target string, char rune, length int) string {
	lng := 0
	var retstr string
	for _, r := range target {
		if r == char {
			if lng < length {
				retstr += string(r)
				lng++
			}
		} else {
			retstr += string(r)
			lng = 0
		}
	}
	return retstr
}

//─────────────┤ LeadingWs ├─────────────

func LeadingWs(s string) string { //returns leading ws
	var ws []rune
	for _, r := range s {
		if !unicode.IsSpace(r) {
			break
		}
		ws = append(ws, r)
	}
	return string(ws)
}

//─────────────┤ StartsWith ├─────────────

func StartsWith(line, search string) bool {
	ws := LeadingWs(line)
	return strings.HasPrefix(line[len(ws):], search)
}

//─────────────┤ StripComment ├─────────────

func StripComment(line, comment string) string {
	for {
		line = strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(line, comment) {
			break
		}
		line = line[len(comment):]
	}
	return line
}

//─────────────┤ SplitOnDashDash ├─────────────

func SplitOnDashDash(inputs, outputs []string) ([]string, []string) {
	for i, s := range outputs {
		if s == "--" {
			inputs = append(inputs, outputs[i+1:]...)
			return inputs, outputs[:i]
		}
	}
	return inputs, outputs
}

type SubStringStat struct {
	Index  int
	Length int
	Before string
	Match  string
	After  string
}

//─────────────┤ FindStringinString ├─────────────

func StringInStringStat(s, sub string) SubStringStat {
	n := strings.Index(s, sub)
	if n == -1 {
		return SubStringStat{Before: s}
	}
	return SubStringStat{Index: n,
		Length: len(sub),
		Before: s[:n],
		Match:  s[n : n+len(sub)],
		After:  s[n+len(sub):],
	}
}

//─────────────┤ getMultilineQuotedStr ├─────────────

func MultilineQuotedStr(lines []string) string {
	mlstr := strings.Join(lines, "\n")
	m := strings.Index(mlstr, "`")
	if m != -1 {
		mlstr = mlstr[m+1:]
		n := strings.LastIndex(mlstr, "`")
		if n != -1 {
			mlstr = mlstr[:n]
		}
	}

	return mlstr
}
