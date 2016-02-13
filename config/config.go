package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type WlanConfig struct {
	Interface string
	SSID      string
	PSK       string
}

type Config struct {
	Networks  []WlanConfig
	Interface string // legacy
	SSID      string // legacy
	PSK       string // legacy
}

type ConfigMessage struct {
	Config struct {
		WLAN Config
	}
}

func (config *Config) Save(path string) {

	// Dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// Write the file
	ioutil.WriteFile(path, data, 0644)
}

func (config *Config) Load(path string) {

	// Does the file exist?
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	// Read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		Warning.Fatalf("load: %v", err)
	}

	Trace.Printf("loaded %v", string(data))

	// Parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}

}

func (config *Config) Extend(element WlanConfig) {
	n := len(config.Networks)
	if n == cap(config.Networks) {
		newSlice := make([]WlanConfig, len(config.Networks), 2*len(config.Networks)+1)
		copy(newSlice, config.Networks)
		config.Networks = newSlice
	}
	config.Networks = config.Networks[0 : n+1]
	config.Networks[n] = element

}

func (config *Config) Upgrade() {

	// Move network outside of the array
	if config.Interface == "" &&
		config.SSID == "" &&
		config.PSK == "" {

		return
	}

	// Move network into the array
	legacyNetwork := WlanConfig{}

	legacyNetwork.Interface = config.Interface
	legacyNetwork.SSID = config.SSID
	legacyNetwork.PSK = config.PSK

	config.Interface = ""
	config.SSID = ""
	config.PSK = ""

	config.Extend(legacyNetwork)
}

func (config *Config) Merge(branch Config) {

	// Merge networks
	if len(branch.Networks) > 0 {
		config.Networks = branch.Networks
	}
}

func (config *ConfigMessage) Scan() {

	// Read from standard-in
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
