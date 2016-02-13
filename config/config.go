package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type NetworkConfig struct {
	Interface     string
	ID            string
	Protocol      string
	SSID          string
	ScanSSID      string
	PSK           string
	KeyManagement string
	Pairwise      string
	Group         string
	AuthAlgorithm string
	Priority      string
}

type Config struct {
	Networks []NetworkConfig
}

type ConfigMessage struct {
	Config struct {
		WLAN Config
	}
}

func (config *NetworkConfig) Export(file *os.File) {

	if config.Interface == "" {
		Info.Fatalln("Using defualt interface wlan0")
		config.Interface = "wlan0"
	}

	fmt.Fprintf(file, "network={\n")
	if config.ID != "" {
		fmt.Fprintf(file, "\tid_str=%v\n", config.ID)
	}

	switch {
	case config.Protocol == "WPA2" || config.Protocol == "RSN":
		fmt.Fprintf(file, "\tproto=%v\n", "RSN")
	case config.Protocol == "WPA":
		fmt.Fprintf(file, "\tproto=%v\n", "WPA")
	case config.Protocol == "WEP":
	case config.Protocol == "":
	default:
		Warning.Fatalln("Protocol must be in WPA2,RSN,WPA,WEP")
	}

	if config.SSID == "" {
		Warning.Fatalln("SSID is required")
	}
	fmt.Fprintf(file, "\tssid=%v\n", config.SSID)

	if config.ScanSSID != "" {
		fmt.Fprintf(file, "\tscan_ssid=%v\n", config.ScanSSID)
	}

	if config.PSK == "" {
		Warning.Fatalln("PSK is required")
	}

	if config.Protocol == "WEP" {
		fmt.Fprintf(file, "\twep_tx_keyidx=%v\n", 0)
		fmt.Fprintf(file, "\twep_key0=%v\n", config.PSK)
	} else {
		fmt.Fprintf(file, "\tpsk=%v\n", config.PSK)
	}

	if config.KeyManagement != "" {
		fmt.Fprintf(file, "\tkey_mgmt=%v\n", config.KeyManagement)
	}

	if config.Pairwise != "" {
		fmt.Fprintf(file, "\tpairwise=%v\n", config.Pairwise)
	}

	if config.Group != "" {
		fmt.Fprintf(file, "\tgroup=%v\n", config.Group)
	}

	if config.AuthAlgorithm != "" {
		fmt.Fprintf(file, "\tauth_alg=%v\n", config.AuthAlgorithm)
	}

	if config.Priority != "" {
		fmt.Fprintf(file, "\tpriority=%v\n", config.Priority)
	}

	fmt.Fprintf(file, "}\n")
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

func (config *Config) Export(interfacesDirPath string) {

	interfaces := make(map[string][]NetworkConfig)

	for _, network := range config.Networks {
		interfaces[network.Interface] = append(interfaces[network.Interface], network)
	}

	for iface, networks := range interfaces {

		interfacePath := path.Join(interfacesDirPath, "interface_"+iface+".conf")

		f, err := os.Create(interfacePath)
		if err != nil {
			Warning.Fatalf("export: %v", err)
		}

		//It’s idiomatic to defer a Close immediately after opening a file.

		defer f.Close()

		for _, network := range networks {
			network.Export(f)
			f.Sync()
		}

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
