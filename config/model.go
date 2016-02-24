package config

import (
	"github.com/deckarep/golang-set"
	"github.com/funkygao/golib/observer"
)

const INTERFACE_DEFAULT = "wlan0"

const ADDRESS_TYPE_DYNAMIC = "dhcp"
const ADDRESS_TYPE_STATIC = "static"

const CONNECTION_TYPE_CLIENT = "client"
const CONNECTION_TYPE_ACCESSPOINT = "accesspoint"
const CONNECTION_TYPE_DEFAULT = CONNECTION_TYPE_CLIENT

var CONNECTION_TYPE_OPTIONS = mapset.NewSetFromSlice([]interface{}{CONNECTION_TYPE_CLIENT, CONNECTION_TYPE_ACCESSPOINT})

const HARDWARE_MODE_DEFAULT = "g"

var HARDWARE_MODE_OPTIONS = mapset.NewSetFromSlice([]interface{}{"a", "b", "g", "n"})

const EVENT_DELETE_INTERFACE = "delete interface"

type IPConfig struct {
	Address string
	Netmask string
	Network string
	Gateway string
}

type WPAConfig struct {
	Protocol      string // eg. WPA, WPA2, WEP
	KeyManagement string `yaml:"key_management"` // eg. WPA-PSK
	Pairwise      string // eg. CCMP or TKIP
	Group         string // eg. TKIP or CCMP

}

type WEPConfig struct {
	AuthAlgorithm string `yaml:"auth_algorithm"` // SHARED for WEP-shared
	Priority      string // for WEP-shared
}

type WifiConfig struct {
	Interface      string // default wlan0
	ID             string // descriptional name
	ConnectionType string `yaml:"connection_type"` // eg client, accesspoint
	SSID           string // network id
	ScanSSID       string `yaml:"scan_ssid"` // default 0, hidden network
	HardwareMode   string // eg. a, b, g, n
	PSK            string // network password
	WPA            WPAConfig
	WEP            WEPConfig
	IP             IPConfig
}

func (config *WifiConfig) GetConnectionType() string {
	return StringCoalesce(config.ConnectionType, CONNECTION_TYPE_DEFAULT)
}

func (config *WifiConfig) GetInterfaceId() string {
	return StringCoalesce(config.ID, config.Interface, INTERFACE_DEFAULT)
}

func (config *IPConfig) GetAddressType() string {
	if config.Address == "" {
		return ADDRESS_TYPE_DYNAMIC
	} else {
		return ADDRESS_TYPE_STATIC
	}
}

type Config struct {
	Networks []WifiConfig
}

func (config *Config) Merge(branch Config) {

	// merge networks
	if len(branch.Networks) > 0 {
		oldConfig := *config
		config.Networks = branch.Networks

		// detect changes
		oldInterfaces := oldConfig.GetInterfaces()
		newInterfaces := config.GetInterfaces()

		deletedInterfaces := oldInterfaces.Difference(newInterfaces)

		for interfaceName := range deletedInterfaces.Iter() {

			observer.Publish(EVENT_DELETE_INTERFACE, interfaceName)

		}
	}

	Trace.Print("merged: %v", config)
}

func (config *Config) GetInterfaces() mapset.Set {
	interfaces := mapset.NewSet()
	for _, network := range config.Networks {
		interfaces.Add(network.Interface)
	}
	return interfaces
}

func (config *Config) Upgrade() {
	// hook for future version changes

	Trace.Print("upgraded: %v", config)
}
