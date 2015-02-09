package tracks

import (
	"flag"
	"log"
	"regexp"
	"strconv"
	"sync"

	"github.com/rakyll/globalconf"
)

// move nextNode and exists to here ...
// and it should use the provider

// je moet een provider in Tracks stoppen, zodat je fatsoenlijk kan
// unit testen met een mock van de echte provider

// prepare for clean

type TrackId string

type Provider interface {
	Provide() TrackId
}

// implementation

var refPatt *regexp.Regexp

func init() {
	refPatt = regexp.MustCompile(`^([[:alpha:]]+)([[:digit:]]+)$`)
}

var provider Provider

func init() {
	current := flag.Int("current", -1, "do not set unless you know what you're doing...")

	conf, err := globalconf.New("tracks")

	if err != nil {
		log.Fatalln(err)
	}

	conf.ParseAll()

	if *current < 0 {
		*current = 0
	}

	provider = &prvdr{conf: conf, current: *current}
}

type prvdr struct {
	conf    *globalconf.GlobalConf
	current int
}

func (p *prvdr) Provide() TrackId {
	mutex := &sync.Mutex{}

	mutex.Lock()
	defer func() { mutex.Unlock() }()

	c := p.current
	p.current++
	p.conf.Set("", &flag.Flag{Name: "current", Value: newFlagValue(p.current)})

	id := "m" + strconv.Itoa(c)

	return TrackId(id)
}

type flagValue struct {
	val int
}

func (f *flagValue) String() string {
	return strconv.Itoa(f.val)
}

func (f *flagValue) Set(value string) error {
	val, err := strconv.Atoi(value)

	if err == nil {
		f.val = val
	}

	return err
}

func newFlagValue(val int) *flagValue {
	return &flagValue{val: val}
}
