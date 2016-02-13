package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type WlanConfig struct {
	SSID string
	PSK  string
}

type ConfigMessage struct {
	Config struct {
		WLAN WlanConfig
	}
}

func (config *WlanConfig) Save(path string) {

	// Dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// Write the file
	ioutil.WriteFile(path, data, 0644)
}

func (config *WlanConfig) Load(path string) {

	// Does the file exist?
	if _, err := os.Stat(path); err == nil {

		// Read the file
		data, err := ioutil.ReadFile(path)
		if err != nil {
			Warning.Fatalf("load: %v", err)
		}

		// Parse the YAML
		if err := yaml.Unmarshal([]byte(data), &config); err != nil {
			Warning.Fatalf("parse: %v", err)
		}
	}

}

func (config *WlanConfig) Merge(branch WlanConfig) {

	if branch.SSID != "" {
		config.SSID = branch.SSID
	}

	if branch.PSK != "" {
		config.PSK = branch.PSK
	}
}

func (config *ConfigMessage) Scan() {

	// Read from StandardIn
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Warning.Fatalf("scan: %v", err)
	}

	// Parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}
}

func (config *ConfigMessage) Print() {

	// Dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// Print the dump
	fmt.Printf("%s", string(data))
}
