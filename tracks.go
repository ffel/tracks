// Package tracks keeps track of notes
package tracks

// implement pandoc Filter interface

import "github.com/ffel/pandocfilter"

type Tracks struct {
	Provider

	checkedDoc bool    // meta is read and track attribute is assessed
	trackDoc   bool    // doc has tracks attribute
	node       TrackId // doc node (iff trackDoc)
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
			return false, value
		}

		if tr.exists(slice[0].(string)) {
			return false, value
		}

		slice[0] = string(tr.Provide())

		return false, value
	}

	return true, value
}
