package config

import (
	"github.com/deckarep/golang-set"
	"github.com/funkygao/golib/observer"
)

const defaultInterface = "wlan0"
const EVENT_DELETE_INTERFACE = "delete interface"

type NetworkConfig struct {
	Interface     string // default wlan0
	ID            string // descriptional name
	Protocol      string // eg. WPA, WPA2, WEP
	SSID          string // network id
	ScanSSID      string // default 0, hidden network
	PSK           string // network password
	KeyManagement string // eg. WPA-PSK
	Pairwise      string // eg. CCMP or TKIP
	Group         string // eg. TKIP or CCMP
	AuthAlgorithm string // SHARED for WEP-shared
	Priority      string // for WEP-shared
}

type Config struct {
	Networks []NetworkConfig
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
