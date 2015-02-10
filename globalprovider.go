package tracks

// the provider provides unique node id's for the real world application

import (
	"flag"
	"log"
	"strconv"
	"sync"

	"github.com/rakyll/globalconf"
)

// this implements a basic singleton like approach as to get
// the global provider
func GetProvider() Provider {
	return provider
}

// no init but injectable
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

	id := "n" + strconv.Itoa(c)

	return TrackId(id)
}

func (p *prvdr) Prefix() string {
	return "n"
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
