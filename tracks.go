// Package tracks keeps track of notes
package tracks

import (
	"fmt"
	"os"
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
		// if the header contains no ref, pandoc adds the default
		// so, there is always one
		headerId, err := pandocfilter.GetString(c, "1", "0")

		if err == nil {
			fmt.Fprintf(os.Stderr, "%v (%T)\n", headerId, headerId)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}

	if ok && t == "Header" {
		// //obj, _ := pandocfilter.GetObject(value, "c")

		// slice := c.([]interface{})

		// slice[0] = tr.Prefix + strconv.Itoa(tr.Current)

		// tr.Current++

		// return false, value

		obj, _ := pandocfilter.GetObject(c, "1")

		slice := obj.([]interface{})

		slice[0] = tr.Prefix + strconv.Itoa(tr.Current)

		tr.Current++

		return false, value

	}

	// kan ik schrijven wanneer ik 'm als object heb

	return true, value
}
