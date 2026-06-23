package oid

import (
	"strings"
)

/*
condense any consecutive combinations of SPACE, HORIZONTAL
TAB or NEW LINE with a single SPACE rune.
*/
func condenseWHSP(b string) string {
	b = strings.TrimSpace(b)
	a := strings.Builder{}

	var last bool
	for i := 0; i < len(b); i++ {
		c := rune(b[i])
		switch c {
		case ' ', '\n', '\t':
			if !last {
				last = true
				a.WriteRune(' ')
			}
		default:
			if last {
				last = false
			}
			a.WriteRune(c)
		}
	}

	return a.String()
}
