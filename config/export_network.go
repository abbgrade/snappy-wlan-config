package config

import (
	"fmt"
	"strings"
)

type NetworkExport struct {
	Lines       []string
	InterfaceId string
}

func NewNetworkExport(config *WifiConfig) NetworkExport {
	export := NetworkExport{}

	export.AddLines(config)

	return export
}

func (export *NetworkExport) Append(key, value string) {

	// export key value pair
	export.Lines = append(export.Lines, fmt.Sprintf("\t%v=\"%v\"", key, value))

}

func (export *NetworkExport) AddLines(config *WifiConfig) {

	export.InterfaceId = config.GetInterfaceId()
	addressType := config.IPConfig.GetAddressType()

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
