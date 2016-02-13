package main

import (
	config "github.com/abbgrade/snappy-wlan-config/config"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	config.InitLogging(ioutil.Discard, os.Stdout, os.Stderr, os.Stderr)

	// get config path
	appDataPath := os.Getenv("SNAP_APP_DATA_PATH")
	if appDataPath == "" {
		appDataPath = "."
	}
	configPath := path.Join(appDataPath, "config.yaml")

	// load
	data := config.Config{}
	data.Load(configPath)
	config.Trace.Printf("loaded: %v from %v", data, configPath)

	// upgrade load
	data.Upgrade()
	config.Trace.Print("upgraded: %v", data)

	// scan
	messageIn := config.ConfigMessage{}
	messageIn.Scan()
	config.Trace.Print("scanned: %v", messageIn)

	// upgrade scan
	messageIn.Config.WLAN.Upgrade()
	config.Trace.Print("upgraded: %v", messageIn)

	// merge load and scan
	data.Merge(messageIn.Config.WLAN)
	config.Trace.Print("merged: %v", data)

	// save merge
	data.Save(configPath)
	config.Trace.Print("saved: %v", data)

	// print save
	messageOut := config.ConfigMessage{}
	messageOut.Config.WLAN = data
	messageOut.Print()
	config.Trace.Print("printed: %v", messageOut)
}
