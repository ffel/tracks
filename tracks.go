// Package tracks keeps track of notes
package tracks

import (
	"strconv"

	"github.com/ffel/pandocfilter"
)

type Tracks struct {
	Prefix  string
	Current int
}

// Value implements pandocfilter Filter interface
func (tr *Tracks) Value(key string, value interface{}) (bool, interface{}) {

	ok, t, c := pandocfilter.IsTypeContents(value)

	if ok && t == "Header" {
		// we need access to the collection, then we can set a new value
		slice, err := pandocfilter.GetSlice(c, "1")

		if err != nil || len(slice) < 1 {
			return true, value
		}

		slice[0] = tr.Prefix + strconv.Itoa(tr.Current)

		tr.Current++

		return false, value
	}

	// value appears to be significant here ...
	return true, value
}
