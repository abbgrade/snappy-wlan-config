package config

import (
	"fmt"
	"github.com/deckarep/golang-set"
)

type InterfaceExport struct {
	Export
	InterfaceId string
}

func NewInterfaceExport(config *WifiConfig) InterfaceExport {
	export := InterfaceExport{}
	export._keyValueFormat = "\t%v %v"

	export.AddLines(config)

	return export
}

func (export *InterfaceExport) AddLines(config *WifiConfig) {

	export.InterfaceId = config.GetInterfaceId()
	addressType := config.IP.GetAddressType()

	export.Extend(fmt.Sprintf("allow-hotplug %v", export.InterfaceId))
	export.Extend(fmt.Sprintf("iface %v inet %v", export.InterfaceId, addressType))

	switch config.GetConnectionType() {
	case CONNECTION_TYPE_CLIENT:
		wifiConfigPath := GetWifiConfigPath(config.Interface)
		export.Append("wpa-conf", wifiConfigPath, false)
	case CONNECTION_TYPE_ACCESSPOINT:
		wifiConfigPath := GetAccesspointConfigPath(config.Interface)
		export.Append("hostapd", wifiConfigPath, false)
	}

	if config.IP.Address != "" {
		export.Append("address", config.IP.Address, false)
		export.Append("netmask", config.IP.Netmask, false)
		export.Append("network", config.IP.Network, false)
		export.Append("gateway", config.IP.Gateway, false)
	}
}

// Controller extension

func GetNetworkConfigPath(interfaceName string) string {
	return fmt.Sprintf("/etc/network/interfaces.d/%v", interfaceName)
}

func (config *Controller) ExportInterface(interfaceName string, networks []WifiConfig) {
	path := GetNetworkConfigPath(interfaceName)
	export := OpenExportFile(path)
	defer export.Close()

	// export unique networks
	networkExports := mapset.NewSet()
	for _, network := range networks {

		networkExport := NewInterfaceExport(&network)
		networkExports.Add(networkExport.Dump())

	}

	// add a file header
	export.AddHeader("interfaces")

	// add each  export
	for network := range networkExports.Iter() {

		export.Extend(network.(string))
		export.Flush()

	}
}
