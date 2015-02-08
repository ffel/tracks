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
		checked, hastrack, val := doctrack(value)

		if checked {
			tr.checkedDoc = true
			tr.trackDoc = hastrack

			if tr.trackDoc && !tr.exists(val) {
				tr.node = tr.nextNode()

				ok2, meta := pandocfilter.IsMeta(value)

				if !ok2 {
					panic("ok was already established")
				}

				meta["track"] = map[string]interface{}{
					"t": "MetaInlines",
					"c": []interface{}{
						map[string]interface{}{
							"t": "Str",
							"c": tr.node,
						},
					},
				}
			} else {
				tr.node = val
			}

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

// doctrack checks val for track attribute
//
// meta is true in case having track attribute or not could be established
// track is true if such an attribute could be found
// node contains track string value (if any, otherwise "")
func doctrack(val interface{}) (ismeta, track bool, node string) {
	ok, meta := pandocfilter.IsMeta(val)

	if !ok {
		return false, false, ""
	}

	tr, hasTrack := meta["tracks"]

	if !hasTrack {
		return true, false, ""
	}

	istc, t, c := pandocfilter.IsTypeContents(tr)

	if !istc {
		return true, false, ""
	}

	if t != "MetaInlines" {
		return true, true, ""
	}

	content, err := pandocfilter.GetString(c, "0", "c")

	if err != nil {

		println("error", err.Error())

		return true, true, ""
	}

	return true, true, content
}
