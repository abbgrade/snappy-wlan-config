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

	// handle arguments
	dryRun := flag.Bool("d", false, "Dry Run: \tDon't change anything on the system.")
	inputPath := flag.String("i", INPUT_FILE_DEFAULT, "Input File: \tRead Input from file instead of stdin")
	flag.Parse()

	// setup logging
	config.InitLogging(ioutil.Discard, os.Stdout, os.Stderr, os.Stderr)

	config.Trace.Printf("dryRun = %v", *dryRun)
	config.Trace.Printf("inputPath = %v", *inputPath)

	// get environment variables
	appDataDirPath := config.StringCoalesce(os.Getenv("SNAP_APP_DATA_PATH"), ".")
	configPath := path.Join(appDataDirPath, "config.yaml")

	// init controller
	controller := config.InitController(appDataDirPath, configPath, *dryRun)

	// load
	controller.Load()
	controller.Model.Upgrade()

	// scan request
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
	request.Config.Model.Upgrade()

	// merge load and scan
	controller.Merge(request)

	if *dryRun == false {
		// save merge
		controller.Save()
	}

	// print response
	response := config.Transaction{}
	response.Config.Model = controller.Model
	response.Print()

	// export effected config files
	controller.Export()

}
