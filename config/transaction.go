package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Transaction struct {
	Config struct {
		WLAN Config
	}
}

func (config *Transaction) Scan() {

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

func (config *Transaction) Print() {

	// dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// print the dump
	fmt.Printf("%s", string(data))
}
