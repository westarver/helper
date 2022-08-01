package helper

import (
	"path/filepath"
	"strings"
)

type IoPair struct {
	in string
	out string
}

const defExt = ".match"
//─────────────────────┤ Matchio ├─────────────────────
//Matchio will match up input files with output files accounting
//for mismatched lengths of slices. Extra outs will be ignored,
//and by default extra ins will match to the last out,
//appending output from multiple templates into a single
//out file. If no output file names are given, matchio will
//simply replace the extension of each template name to a
//default extension to generate the output file names.
//This behaviour can be changed by offering a 'wildcard'
//as an out file name. An out file given as "/s suffix"
//will concatenate "suffix" to in file to create the out
//file name ('source.tpl' becomes 'source.tplsuffix').
//Using "/S suffix" will append the suffix to the basename,
//the extension remains ("source.tpl >> sourcesuffix.tpl").
//A wildcard of "/p prefix" will prepend the given text
//('source.tpl >> 'prefixsource.tpl).  Using "/e ext"
//will replace the extension of the template with the given
//extension ('source.tpl >> source.ext'). To replace the
//extension with a blank use a / as the arg for /e.
//("/e /" results in 'src.ex.tpl >> src.ex')
//These can be combined as in "/d/p/s/S/e dir pre suf ext".
//As expected "/p/e pre ext" will add the prefix and
//replace the extension. "/S/e suffix ext" will produce
//'sourcesuffix.ext'.
//One more wildcard: "/d directory" will prepend the given dir
//to any output file name however derived. The dir will be
//created in the cwd if it does not exist.
//The above examples showed the directives separated with a /
//for visual reasons only. They can be but it's not required.
//"/d/p/e d p e" is equivalent to "/dpe d p e".  The directives
//will work in any order, but the args must be in the same order
//as the directives they refer to.
func Matchio(ins []string, outs []string) []IoPair {
	var matched []IoPair

	for _, f := range ins {
		matched = append(matched, IoPair{f, "x"}) //out is dummy string for now
	}

	ilen := len(ins)
	olen := len(outs)

	if olen == 0 { // easy one first
		for i := 0; i < len(matched); i++ {
			matched[i].out = matched[i].in + defExt
		}
		return matched
	}

	mismatch := ilen - olen
	if mismatch < 0 {
		mismatch = 0 //discard extra outs
	}
	if mismatch == 0 { // all ins have a matching out. match them up.  later we check for wildcards.
		for i := 0; i < len(matched); i++ {
			matched[i].out = outs[i]
		}
	} else { // here we know there is at least one out but fewer outs than ins
		for i := 0; i < olen; i++ { //copy as many outs as were given to the matched slice
			matched[i].out = outs[i]
		}
		for i := olen; i < ilen; i++ { //fill in remainder of matched slice with last out
			matched[i].out = outs[olen-1]
		}
	}

	// now we check for and parse wildcards
	for i, m := range matched {
		if strings.HasPrefix(m.out, "/") {
			// maybe not a file name but a slash directive,  but
			// we do need to check for a rooted absolute path
			// if it is we assume that was intended
			// use 2 slashes to denote that intent
			// ex. //home/me/my/stuff/template.tpl
			if len(m.out) > 1 && m.out[1] == '/' {
				matched[i].out = m.out[1:]
				continue
			}

			ss := strings.Split(m.out, " ")
			//ss holds slash directives in ss[0] and the individual tokens in the remaining indices
			tmp := m.in // file name is going to be the same as input with possible alterations

			// get rid of slashes
			var direct string
			for _, c := range ss[0] {
				if c == '/' {
					continue
				}
				direct += string(c)
			}
			// get slice of args
			var args []string
			args = append(args, ss[1:]...)
			// if there are fewer args tha directives the
			// following loop would panic due to index out of range
			if len(args) < len(direct) {
				matched[i].out = m.in + defExt
				continue
			}
			//start a new loop
			for argnum, op := range direct {
				switch op {
				case 'd':
					d := args[argnum]
					if !strings.HasSuffix(d, "/") {
						d += "/"
					}
					if strings.HasPrefix(tmp, "../") {
						tmp = tmp[2:]
					}
					if strings.HasPrefix(tmp, "./") {
						tmp = tmp[1:]
					}
					tmp = d + tmp
				case 'p':
					pre := args[argnum]
					p, _ := filepath.Split(tmp)
					b := filepath.Base(tmp)
					tmp = p + pre + b
				case 'S':
					s := args[argnum]
					p, _ := filepath.Split(tmp)
					e := filepath.Ext(tmp)
					b := strings.TrimSuffix(filepath.Base(tmp), e)
					tmp = p + b + s + e
				case 'n':
					n := args[argnum]
					p, _ := filepath.Split(tmp)
					e := filepath.Ext(tmp)
					tmp = p + n + e
				case 'e':
					e := args[argnum]
					if !strings.HasPrefix(e, ".") {
						e = "." + e
					}
					// user used "/" as ext meaning remove the ext
					if e == "./" {
						e = ""
					}
					p, _ := filepath.Split(tmp)
					x := filepath.Ext(tmp)
					b := strings.TrimSuffix(filepath.Base(tmp), x)
					tmp = p + b + e
				case 's':
					tmp += args[argnum]
				default: // unrecognized char
					continue
				}
			}
			matched[i].out = tmp
			if len(m.out) == 0 {
				matched[i].out = m.in + defExt
			}
		}
	}
	return matched
} // matchio

