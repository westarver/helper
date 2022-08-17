package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	apath "github.com/rhysd/abspath"
)

//─────────────┤ GetLogFile ├─────────────

func GetLogFile(file string) *os.File {
	if file == "std" {
		return os.Stderr
	}
	f, _ := OpenFile(file)
	if f == nil {
		fmt.Println("Failure to open Log file " + file + ". Will default to Stderr.")
		f = os.Stderr
	}
	return f
}

//─────────────┤ OpenFile ├─────────────

func OpenFile(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	}
	return nil, err
}

//─────────────┤ OpenAppend ├─────────────

func OpenAppend(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	}
	return nil, err
}

//─────────────┤ OpenTrunc ├─────────────

func OpenTrunc(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	}
	return nil, err
}

//─────────────┤ CopyFileStr ├─────────────

func CopyFileStr(dst, src string) error {
	if DirExists(dst) {
		_, base := filepath.Split(src)
		dst = filepath.Join(dst, base)
	}
	d, err := OpenTrunc(dst)
	if err != nil {
		return err
	}
	s, err := OpenFile(src)
	if err != nil {
		return err
	}
	_, _ = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}

//─────────────┤ GetPath  ├─────────────

func GetPath(name string) (string, error) {
	a, err := apath.ExpandFrom(name)
	return a.String(), err
}

//─────────────┤ DirExists ├─────────────

func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

//─────────────┤ ValidatePath ├─────────────

func ValidatePath(path string) (string, error) {
	a, err := apath.ExpandFrom(path)
	if err != nil {
		return "", err
	}
	path = a.String()
	d, f := filepath.Split(path)
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return path, err
	}
	return d + f, nil
}

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

type RegexStat struct {
	Before, Match, After string
	Start, Length        int
}

//─────────────┤ mapFromPatSlice ├─────────────

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
