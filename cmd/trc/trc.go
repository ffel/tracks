package main

import (
	"github.com/ffel/pandocfilter"
	"github.com/ffel/tracks"
)

func main() {
	pandocfilter.Run(&tracks.Tracks{})
}
