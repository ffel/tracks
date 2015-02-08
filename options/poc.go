package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/rakyll/globalconf"
)

var (
	current = flag.Int("current", -1, "do not set unless you know what you're doing...")
)

func main() {
	// read from / write to ~/.config/appname/config.ini
	conf, err := globalconf.New("appname")
	conf.ParseAll()

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	println(*current)

	if *current < 0 {
		*current = 0
	}

	*current++

	conf.Set("", &flag.Flag{Name: "current", Value: newFlagValue(*current)})
}

// from https://github.com/rakyll/globalconf/blob/master/globalconf_test.go

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
