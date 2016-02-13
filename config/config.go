package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type NetworkConfig struct {
	Interface string
	SSID      string
	PSK       string
}

type Config struct {
	Networks []NetworkConfig
}

type ConfigMessage struct {
	Config struct {
		WLAN Config
	}
}

func (config *Config) Save(path string) {

	// dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// write the file
	ioutil.WriteFile(path, data, 0644)
}

func (config *Config) Load(path string) {

	// does the file exist?
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	// read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		Warning.Fatalf("load: %v", err)
	}

	Trace.Printf("loaded %v", string(data))

	// parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}

}

func (config *Config) Extend(element NetworkConfig) {
	n := len(config.Networks)
	if n == cap(config.Networks) {
		newSlice := make([]NetworkConfig, len(config.Networks), 2*len(config.Networks)+1)
		copy(newSlice, config.Networks)
		config.Networks = newSlice
	}
	config.Networks = config.Networks[0 : n+1]
	config.Networks[n] = element

}

func (config *Config) Upgrade() {
	// hook for future version changes
}

func (config *Config) Merge(branch Config) {

	// merge networks
	if len(branch.Networks) > 0 {
		config.Networks = branch.Networks
	}
}

func (config *ConfigMessage) Scan() {

	// read from standard-in
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Warning.Fatalf("scan: %v", err)
	}

	// parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}
}

func (config *ConfigMessage) Print() {

	// dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// print the dump
	fmt.Printf("%s", string(data))
}
