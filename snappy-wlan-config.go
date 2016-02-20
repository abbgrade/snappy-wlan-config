package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"

	config "github.com/abbgrade/snappy-wlan-config/config"
)

const INPUT_FILE_STDIN = "#stdin"
const INPUT_FILE_DEFAULT = INPUT_FILE_STDIN

func main() {

	dryRun := flag.Bool("d", false, "Dry Run: \tDon't change anything on the system.")
	inputPath := flag.String("i", INPUT_FILE_DEFAULT, "Input File: \tRead Input from file instead of stdin")

	flag.Parse()

	config.InitLogging(ioutil.Discard, os.Stdout, os.Stderr, os.Stderr)

	config.Info.Printf("dryRun = %v", *dryRun)
	config.Info.Printf("inputPath = %v", *inputPath)

	if *dryRun == true {
		config.EnableDryRun()
	}

	// get environment variables
	appDataDirPath := config.StringCoalesce(os.Getenv("SNAP_APP_DATA_PATH"), ".")

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
	inputFile, _ := os.Open(os.DevNull)
	switch {
	case *inputPath == INPUT_FILE_STDIN:
		inputFile = os.Stdin
	case *inputPath != "":

		inputFile_, err := os.Open(*inputPath)
		if err != nil {
			config.Warning.Fatalf("read %v : %v", *inputPath, err)
		} else {
			inputFile = inputFile_
		}
	}

	request.Scan(inputFile)
	config.Trace.Print("scanned: %v", request)

	// upgrade scan
	request.Config.Model.Upgrade()
	config.Trace.Print("upgraded: %v", request)

	// merge load and scan
	controller.Merge(request)
	config.Trace.Print("merged: %v", controller)

	if *dryRun == false {
		// save merge
		controller.Save()
		config.Trace.Print("saved: %v", controller)
	}

	// print save
	response := config.Transaction{}
	response.Config.Model = controller.Model
	response.Print()
	config.Trace.Print("printed: %v", response)

	// export wpa supplicant config
	controller.Export()

}
