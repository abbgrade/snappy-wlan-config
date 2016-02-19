package config

import (
	"github.com/deckarep/golang-set"
	"github.com/funkygao/golib/observer"
)

const defaultInterface = "wlan0"
const defaultAddressType = "dhcp"
const EVENT_DELETE_INTERFACE = "delete interface"

type IPConfig struct {
	AddressType string `yaml:"address_type"` // eg. dhcp, static
	Address     string
	Netmask     string
	Network     string
	Gateway     string
}

type WifiConfig struct {
	Interface     string   // default wlan0
	ID            string   // descriptional name
	Protocol      string   // eg. WPA, WPA2, WEP
	SSID          string   // network id
	ScanSSID      string   `yaml:"scan_ssid"` // default 0, hidden network
	PSK           string   // network password
	KeyManagement string   `yaml:"key_management"` // eg. WPA-PSK
	Pairwise      string   // eg. CCMP or TKIP
	Group         string   // eg. TKIP or CCMP
	AuthAlgorithm string   `yaml:"auth_algorithm"` // SHARED for WEP-shared
	Priority      string   // for WEP-shared
	IPConfig      IPConfig `yaml:"ip_config"`
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
}
