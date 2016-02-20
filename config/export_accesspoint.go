package config

import (
	"fmt"
	"strings"
)

type AccesspointExport struct {
	Lines []string
}

func NewAccesspointExport(config *WifiConfig) AccesspointExport {

	export := AccesspointExport{}

	export.AddLines(config)

	return export
}

func (export *AccesspointExport) Append(key, value string, optional bool, defaults ...string) {

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

func (export *AccesspointExport) AddLines(config *WifiConfig) {
	export.Append("interface", config.Interface, false, INTERFACE_DEFAULT)
	export.Append("driver", "", false)
	export.Append("ssid", config.SSID, false)
	export.Append("channel", "1", true)
	export.Append("ignore_broadcast_ssid", config.ScanSSID, true)
	export.Append("country_code", "", false)
	export.Append("ieee80211d", "", false)
	export.Append("hw_mode", "", false)
	export.Append("ieee80211n", "", true)

	switch config.Protocol {
	case "WPA2":
		fallthrough
	case "RSN":
		fallthrough
	case "":
		export.Append("wpa", "2", false)
		export.Append("rsn_preauth", "1", false)
		export.Append("rsn_preauth_interfaces", config.Interface, false)
		export.Append("wpa_key_mhmt", "WPA-PSK", false)
		export.Append("rsn_pairwise", config.Pairwise, false)
		export.Append("wpa_group_rekey", "600", true)
		export.Append("wpa_ptk_rekey", "600", true)
		export.Append("wpa_gmk_rekey", "86400", true)
		export.Append("wpa_passphrase", config.PSK, false)
	default:
		Warning.Fatalln("Protocol must be in WPA2,RSN,WPA,WEP")

	}
}

func (export *AccesspointExport) Dump() string {

	return strings.Join(export.Lines, "\n")

}
