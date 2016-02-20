package config

import (
	"fmt"
	"path"
)

type WifiClientExport struct {
	Export
}

func NewWifiClientExport(config *WifiConfig) WifiClientExport {

	export := WifiClientExport{}
	export._keyValueFormat = "\t%v=\"%v\""
	export._prefix = "network={"
	export._suffix = "}"

	export.AddLines(config)

	return export
}

func (export *WifiClientExport) AddLines(config *WifiConfig) {
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

// Controller extension

func (config *Controller) GetWifiConfigPath(interfaceName string) string {
	fileName := fmt.Sprintf("interface_%v.conf", interfaceName)
	return path.Join(config.InterfacesDirPath, fileName)
}

func (config *Controller) ExportWifiClient(interfaceName string, networks []WifiConfig) {
	path := config.GetWifiConfigPath(interfaceName)
	export := OpenExportFile(path)
	defer export.Close()

	// add a file header
	export.AddHeader("wpasupplicant")

	for _, network := range networks {

		if network.GetConnectionType() != CONNECTION_TYPE_CLIENT {
			Trace.Printf("skip %v", network.GetConnectionType())
			continue
		}

		// add each network configuration
		networkExport := NewWifiClientExport(&network)
		export.Extend(networkExport.Dump())
		export.Flush()
	}
}
