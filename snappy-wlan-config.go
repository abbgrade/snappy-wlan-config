package main

import (
	config "github.com/abbgrade/snappy-wlan-config/config"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	config.InitLogging(ioutil.Discard, os.Stdout, os.Stderr, os.Stderr)

	// get environment variables
	appDataDirPath := os.Getenv("SNAP_APP_DATA_PATH")
	if appDataDirPath == "" {
		appDataDirPath = "."
	}

	// init controller
	controller := config.Controller{}

	controller.ConfigPath = path.Join(appDataDirPath, "config.yaml")
	controller.InterfacesDirPath = appDataDirPath

	config.Trace.Printf("parameters = %v", controller)

	// load
	controller.Load()
	config.Trace.Printf("loaded: %v from %v", controller, controller.ConfigPath)

	// upgrade load
	controller.Model.Upgrade()
	config.Trace.Print("upgraded: %v", controller)

	// scan
	request := config.Transaction{}
	request.Scan()
	config.Trace.Print("scanned: %v", request)

	// upgrade scan
	request.Config.Model.Upgrade()
	config.Trace.Print("upgraded: %v", request)

	// merge load and scan
	controller.Merge(request)
	config.Trace.Print("merged: %v", controller)

	// save merge
	controller.Save()
	config.Trace.Print("saved: %v", controller)

	// print save
	response := config.Transaction{}
	response.Config.Model = controller.Model
	response.Print()
	config.Trace.Print("printed: %v", response)

	// export wpa supplicant config
	controller.Export()
}
