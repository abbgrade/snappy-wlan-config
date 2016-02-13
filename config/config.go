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

type Config struct {
	Config struct {
		WLAN WlanConfig
	}
}

func (config *Config) Scan() {

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Warning.Fatalf("read: %v", err)
	}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}
}

func (config *Config) Print() {
	dump, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}
	fmt.Printf("%s", string(dump))
}

func (config *Config) Save(path string) {
	dump, err := yaml.Marshal(&config.Config.WLAN)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}
	ioutil.WriteFile(path, dump, 0644)
}
