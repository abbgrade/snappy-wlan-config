package main

import (
	config "github.com/abbgrade/snappy-wlan-config/config"
	"io/ioutil"
	"os"
)

func main() {
	config.InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	data := config.Config{}
	data.Init()
	data.Dump()
}
