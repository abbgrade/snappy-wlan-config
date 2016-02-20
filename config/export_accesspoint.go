package config

import ()

type AccesspointExport struct {
	Export
}

func NewAccesspointExport(config *WifiConfig) AccesspointExport {

	export := AccesspointExport{}
	export._keyValueFormat = "\t%v=\"%v\""

	export.AddLines(config)

	return export
}

func (export *AccesspointExport) AddLines(config *WifiConfig) {
	export.Append("interface", config.Interface, false, INTERFACE_DEFAULT)
	export.Append("driver", "", true)
	export.Append("ssid", config.SSID, false)
	export.Append("channel", "1", false)
	export.Append("ignore_broadcast_ssid", config.ScanSSID, true)
	export.Append("country_code", "", true)
	export.Append("ieee80211d", "", true)
	export.Append("hw_mode", "", true)
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
		export.Append("rsn_pairwise", config.Pairwise, true)
		export.Append("wpa_group_rekey", "600", true)
		export.Append("wpa_ptk_rekey", "600", true)
		export.Append("wpa_gmk_rekey", "86400", true)
		export.Append("wpa_passphrase", config.PSK, false)
	default:
		Warning.Fatalln("Protocol must be in WPA2,RSN,WPA,WEP")

	}
}

// Conroller extension

func (config *Controller) GetAccesspointConfigPath() string {
	return "/etc/hostapd/hostapd.conf"
}

func (config *Controller) ExportWifiAccesspoint(networks []WifiConfig) {

	for _, network := range networks {

		if network.GetConnectionType() == CONNECTION_TYPE_CLIENT {
			Trace.Printf("skip %v", network.GetConnectionType())
			continue
		}

		path := config.GetAccesspointConfigPath()
		export := OpenExportFile(path)
		defer export.Close()

		// add a file header
		export.AddHeader("hostapd")

		// add each network configuration
		exportAccesspoint := NewAccesspointExport(&network)
		export.Extend(exportAccesspoint.Dump())
		export.Flush()
	}
}
