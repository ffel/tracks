// Package tracks keeps track of notes
package tracks

import (
	"regexp"
	"strconv"

	"github.com/ffel/pandocfilter"
)

var refPatt *regexp.Regexp

func init() {
	refPatt = regexp.MustCompile(`^([[:alpha:]]+)([[:digit:]]+)$`)
}

type Tracks struct {
	Prefix  string
	Current int

	checkedDoc bool   // meta is read and track attribute is assessed
	trackDoc   bool   // doc has tracks attribute
	node       string // doc node (iff trackDoc)
}

// Value implements pandocfilter Filter interface
func (tr *Tracks) Value(key string, value interface{}) (bool, interface{}) {

	if !tr.checkedDoc {
		if performed := tr.docCheck(value); performed {
			return false, value
		}
	} else if !tr.trackDoc {
		return false, value
	}

	ok, t, c := pandocfilter.IsTypeContents(value)

	if ok && t == "Header" {
		slice, err := pandocfilter.GetSlice(c, "1")

		if err != nil || len(slice) < 1 {
			return true, value
		}

		if tr.exists(slice[0].(string)) {
			return true, value
		}

		slice[0] = tr.nextNode()
		// slice[0] = tr.Prefix + strconv.Itoa(tr.Current)
		// tr.Current++

		return false, value
	}

	return true, value
}

// assign next free node
func (tr *Tracks) nextNode() string {
	defer func() { tr.Current++ }()

	return tr.Prefix + strconv.Itoa(tr.Current)
}

// exists checks whether ref is valid tracks ref
func (tr *Tracks) exists(ref string) bool {
	match := refPatt.FindStringSubmatch(ref)

	if match == nil || len(match) < 3 {
		return false
	}

	if match[1] != tr.Prefix {
		return false
	}

	return true
}
