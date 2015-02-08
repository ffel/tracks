package tracks

import (
	"encoding/json"
	"log"
	"os/exec"
	"testing"

	"github.com/ffel/pandocfilter"
)

var notracks string = `---
title: Titel
...

Tracks
======
`

var newtracks string = `---
title: Titel
tracks: 1
...

Tracks
======
`

var newtracks_bool string = `---
title: Titel
tracks: true
...

Tracks
======
`

var newtracks_string string = `---
title: Titel
tracks: this is no valid node
...

Tracks
======
`

var existingtracks string = `---
title: Titel
tracks: n12345
...

Tracks
======
`

func TestDoctrack(t *testing.T) {
	tests := []struct {
		in     string
		tracks bool
		key    string
	}{
		{notracks, false, ""},
		{newtracks, true, "n100"},
		{newtracks_bool, true, "n100"},
		{newtracks_string, true, "n100"},
		{existingtracks, true, "n12345"},
	}

	for i, tst := range tests {
		filter := &Tracks{Prefix: "n", Current: 100}
		_ = pandocfilter.Walk(filter, "", pandoc2json(tst.in))

		if expTracks := tst.tracks; filter.trackDoc != expTracks {
			t.Errorf("track doc %v - expected %v; got %v\n", i, expTracks, filter.trackDoc)
		}

		if expKey := tst.key; filter.node != expKey {
			t.Errorf("node %v - expected %v; got %v\n", i, expKey, filter.node)
		}
	}
}

func pandoc2json(doc string) interface{} {
	// for some reason
	// cmd := exec.Command("pandoc", "-t json")
	// lets pandoc produce the error "pandoc: Unknown writer:  plain"

	cmd := exec.Command("bash", "-c", "pandoc -t json")

	stdin, err := cmd.StdinPipe()

	if err != nil {
		return "no output 1"
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return "no output 2"
	}

	if err := cmd.Start(); err != nil {
		return "no output 3"
	}

	stdin.Write([]byte(doc))
	stdin.Close()

	decoder := json.NewDecoder(stdout)

	var pandoc interface{}

	if err := decoder.Decode(&pandoc); err != nil {
		log.Println(err)
		return "no output 4"
	}

	// out, err := ioutil.ReadAll(stdout)

	// if err != nil {
	// 	return "no output 5"
	// }

	cmd.Wait()

	return pandoc
}
