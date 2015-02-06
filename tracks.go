// Package tracks keeps track of notes
package tracks

import (
	"fmt"
	"os"

	"github.com/ffel/pandocfilter"
)

type Tracks struct{}

// Value implements pandocfilter Filter interface
func (tr *Tracks) Value(key string, value interface{}) (bool, interface{}) {

	ok, t, c := pandocfilter.IsTypeContents(value)

	if ok && t == "Header" {
		// if the header contains no ref, pandoc adds the default
		// so, there is always one
		obj, err := pandocfilter.GetString(c, "1", "0")

		if err == nil {
			fmt.Fprintf(os.Stderr, "%v (%T)\n", obj, obj)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}

	return true, value
}
