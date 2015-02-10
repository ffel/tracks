package tracks

// abstraction for getting next available node id

import "regexp"

// prepare for clean

type TrackId string

type Provider interface {
	Provide() TrackId
	Prefix() string
}

// implementation

var refPatt *regexp.Regexp

func init() {
	refPatt = regexp.MustCompile(`^([[:alpha:]]+)([[:digit:]]+)$`)
}

// exists checks whether ref is valid tracks ref
func (tr *Tracks) exists(ref string) bool {
	match := refPatt.FindStringSubmatch(ref)

	if match == nil || len(match) < 3 {
		return false
	}

	if match[1] != tr.Provider.Prefix() {
		return false
	}

	return true
}
