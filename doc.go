package tracks

import "github.com/ffel/pandocfilter"

// docCheck performs the doc check and returns true if
// it could have been performed
func (tr *Tracks) docCheck(value interface{}) bool {
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

			meta["tracks"] = map[string]interface{}{
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

		return true
	}

	return false
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
