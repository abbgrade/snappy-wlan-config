package config

import (
	"fmt"
	"strings"
)

type NetworkExport struct {
	Lines       []string
	InterfaceId string
}

func NewNetworkExport(config *NetworkConfig) NetworkExport {
	export := NetworkExport{}

	export.AddLines(config)

	return export
}

func (export *NetworkExport) Append(key, value string) {

	// export key value pair
	export.Lines = append(export.Lines, fmt.Sprintf("\t%v=\"%v\"", key, value))

}

func (export *NetworkExport) AddLines(config *NetworkConfig) {

	export.InterfaceId = StringCoalesce(config.ID, config.Interface, defaultInterface)
	addressType := StringCoalesce(config.IPConfig.AddressType, defaultAddressType)

	export.Lines = append(export.Lines, fmt.Sprintf("iface %v inet %v", export.InterfaceId, addressType))

	if config.IPConfig.AddressType == "static" {
		export.Append("address", config.IPConfig.Address)
		export.Append("netmask", config.IPConfig.Netmask)
		export.Append("network", config.IPConfig.Network)
		export.Append("gateway", config.IPConfig.Gateway)
	}
}

func (export *NetworkExport) Dump() string {

	return strings.Join(export.Lines, "\n")

}

type WifiExport struct {
	Lines []string
}

func NewWifiExport(config *NetworkConfig) WifiExport {

	export := WifiExport{}

	export.AddLines(config)

	return export
}

func (export *WifiExport) Append(key, value string, optional bool, defaults ...string) {

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

func (export *WifiExport) AddLines(config *NetworkConfig) {
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

func (export *WifiExport) Dump() string {

	// wrap content
	export.Lines = append([]string{"network={"}, export.Lines...)
	export.Lines = append(export.Lines, "}")

	return strings.Join(export.Lines, "\n")

}
