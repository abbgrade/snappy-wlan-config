package config

import (
	"fmt"
	"os"
)

type NetworkExport struct {
	Lines []string
}

func (export *NetworkExport) Append(key, value string, optional bool, defaults ...string) {

	// apply defaults
	if value == "" && len(defaults) > 0 {
		value = defaults[0]
	}

	if value != "" {

		// export key value pair
		export.Lines = append(export.Lines, fmt.Sprintf("\t%v=\"%v\"", key, value))

	} else if optional == false {

		// fail on missing non optional value
		Warning.Fatalf("%v is required but not set", key)
	}
}

func (export *NetworkExport) Save(file *os.File) {

	fmt.Fprintf(file, "network={\n")

	for _, line := range export.Lines {
		fmt.Fprintf(file, "%v\n", line)
	}

	fmt.Fprintf(file, "}\n")
}

func (export *NetworkExport) Add(config *NetworkConfig) {
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

}
