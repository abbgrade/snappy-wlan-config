package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	SSID string
	PSK  string
}

func (config *Config) Init() {

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Warning.Fatalf("read: %v", err)
	}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}
}

func (config *Config) Dump() {
	dump, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}
	fmt.Printf("%s", string(dump))
}
