package main

import (
	"flag"
	"log"

	"github.com/rakyll/globalconf"
)

var (
	flagName    = flag.String("name", "default name", "Name of the person.")
	flagAddress = flag.String("addr", "default address", "Address of the person.")
)

func main() {
	// read from / write to ~/.config/appname/config.ini
	conf, err := globalconf.New("appname")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	println(*flagName)
	println(*flagAddress)

	// this does write in ~/.config/appname/config.ini
	conf.Set("demo", &flag.Flag{Name: "a", Value: newFlagValue("test")})

}

// from https://github.com/rakyll/globalconf/blob/master/globalconf_test.go

type flagValue struct {
	str string
}

func (f *flagValue) String() string {
	return f.str
}

func (f *flagValue) Set(value string) error {
	f.str = value
	return nil
}

func newFlagValue(val string) *flagValue {
	return &flagValue{str: val}
}
