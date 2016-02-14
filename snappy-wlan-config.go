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
	config.Trace.Printf("app data path = %v", appDataPath)
	configPath := path.Join(appDataPath, "config.yaml")

	// load
	data := config.Config{}
	data.Load(configPath)
	config.Trace.Printf("loaded: %v from %v", data, configPath)

	// upgrade load
	data.Upgrade()
	config.Trace.Print("upgraded: %v", data)

	// scan
	request := config.Transaction{}
	request.Scan()
	config.Trace.Print("scanned: %v", request)

	// upgrade scan
	request.Config.WLAN.Upgrade()
	config.Trace.Print("upgraded: %v", request)

	// merge load and scan
	data.Merge(request.Config.WLAN)
	config.Trace.Print("merged: %v", data)

	// save merge
	data.Save(configPath)
	config.Trace.Print("saved: %v", data)

	// print save
	response := config.Transaction{}
	response.Config.WLAN = data
	response.Print()
	config.Trace.Print("printed: %v", response)

	data.Export(appDataPath)
}
