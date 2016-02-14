package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type NetworkConfig struct {
	Interface     string // default wlan0
	ID            string // descriptional name
	Protocol      string // eg. WPA, WPA2, WEP
	SSID          string // network id
	ScanSSID      string // default 0, hidden network
	PSK           string // network password
	KeyManagement string // eg. WPA-PSK
	Pairwise      string // eg. CCMP or TKIP
	Group         string // eg. TKIP or CCMP
	AuthAlgorithm string // SHARED for WEP-shared
	Priority      string // for WEP-shared
}

type Config struct {
	Networks []NetworkConfig
}

type ConfigMessage struct {
	Config struct {
		WLAN Config
	}
}

type WPASupplicantExport struct {
	Lines []string
}

func (export *WPASupplicantExport) Append(key, value string, optional bool, defaults ...string) {

	if value == "" && len(defaults) > 0 {
		value = defaults[0]
	}

	if value != "" {
		export.Lines = append(export.Lines, fmt.Sprintf("\t%v=\"%v\"", key, value))
	} else if optional == false {
		Warning.Fatalf("%v is required but not set", key)
	}
}

func (export *WPASupplicantExport) Save(file *os.File) {
	fmt.Fprintf(file, "network={\n")

	for _, line := range export.Lines {
		fmt.Fprintf(file, "%v\n", line)
	}

	fmt.Fprintf(file, "}\n")
}

func (config *NetworkConfig) Export(file *os.File) {

	export := WPASupplicantExport{}

	export.Append("id_str", config.ID, true)
	export.Append("ssid", config.SSID, false)
	export.Append("scan_ssid", config.ScanSSID, true)

	switch config.Protocol {
	case "":
		fallthrough
	case "WPA":
		fallthrough
	case "WPA2":
		fallthrough
	case "RSN":

		export.Append("proto", config.Protocol, true)
		export.Append("psk", config.PSK, false)

		export.Append("key_mgmt", config.KeyManagement, true)
		export.Append("pairwise", config.Pairwise, true)
		export.Append("group", config.Group, true)

	case "WEP":

		export.Append("wep_tx_keyidx", "0", false)
		export.Append("wep_key0", config.PSK, false)

		export.Append("key_mgmt", config.KeyManagement, true, "NONE")

		export.Append("auth_alg", config.AuthAlgorithm, true)
		export.Append("priority", config.Priority, true)

	default:

		Warning.Fatalln("Protocol must be in WPA2,RSN,WPA,WEP")

	}

	export.Save(file)

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
		if network.Interface == "" {
			network.Interface = "wlan0"
		}
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
