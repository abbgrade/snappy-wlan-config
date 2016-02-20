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
	export._keyValueFormat = "\t%v=\"%v\""

	export.AddLines(config)

	return export
}

func (export *InterfaceExport) AddLines(config *WifiConfig) {

	export.InterfaceId = config.GetInterfaceId()
	addressType := config.IPConfig.GetAddressType()

	export.Extend(fmt.Sprintf("iface %v inet %v", export.InterfaceId, addressType))

	if config.IPConfig.AddressType == "static" {
		export.Append("address", config.IPConfig.Address, false)
		export.Append("netmask", config.IPConfig.Netmask, false)
		export.Append("network", config.IPConfig.Network, false)
		export.Append("gateway", config.IPConfig.Gateway, false)
	}
}

// Controller extension

func (config *Controller) GetNetworkConfigPath(interfaceName string) string {
	return fmt.Sprintf("/etc/network/interfaces.d/%v", interfaceName)
}

func (config *Controller) ExportInterface(interfaceName string, networks []WifiConfig) {
	path := config.GetNetworkConfigPath(interfaceName)
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
