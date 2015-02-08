package tracks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

Header 1 {#n987}
========

Header 2
========
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

func Example_Tracks() {
	filter := &Tracks{Prefix: "n", Current: 100}

	json := pandoc2json(newtracks)

	newjson := pandocfilter.Walk(filter, "", json)

	markdown := json2pandoc(newjson)

	fmt.Println(markdown)

	// output:
	// ---
	// title: Titel
	// tracks: n100
	// ...
	//
	// Tracks {#n101}
	// ======
}

func Example_Tracks2() {
	filter := &Tracks{Prefix: "n", Current: 100}

	json := pandoc2json(existingtracks)

	newjson := pandocfilter.Walk(filter, "", json)

	markdown := json2pandoc(newjson)

	fmt.Println(markdown)

	// output:
	// ---
	// title: Titel
	// tracks: n12345
	// ...
	//
	// Header 1 {#n987}
	// ========
	//
	// Header 2 {#n100}
	// ========
}

func pandoc2json(doc string) interface{} {
	// for some reason
	// cmd := exec.Command("pandoc", "-t json")
	// lets pandoc produce the error "pandoc: Unknown writer:  plain"
	// this is a pity, for tests will not work on windows
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

	cmd.Wait()

	return pandoc
}

func json2pandoc(data interface{}) string {

	// we kunnen de huidige applicate testen, niet als de standaard
	// --filter, maar op de omslachtige manier: vanuit json
	// weer naar markdown

	cmd := exec.Command("bash", "-c", "pandoc -s -f json -t markdown")

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

	// from runner - we need pandoc to feed with json strings

	encoder := json.NewEncoder(stdin)

	if err := encoder.Encode(&data); err != nil {
		log.Println(err)
		return ""
	}

	stdin.Close()

	out, err := ioutil.ReadAll(stdout)

	if err != nil {
		return "no output"
	}

	return string(out)
}
