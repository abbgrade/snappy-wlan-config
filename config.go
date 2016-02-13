package main

import (
	//"bufio"
	config "github.com/abbgrade/snappy-wlan-config/config"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	config.InitLogging(ioutil.Discard, os.Stdout, os.Stderr, os.Stderr)

	appDataPath := os.Getenv("SNAP_APP_DATA_PATH")
	if appDataPath == "" {
		appDataPath = "."
	}
	configPath := path.Join(appDataPath, "config.yaml")

	data := config.WlanConfig{}
	data.Load(configPath)

	message := config.ConfigMessage{}
	message.Scan()

	data.Merge(message.Config.WLAN)
	data.Save(configPath)

	message.Config.WLAN = data
	message.Print()
}
