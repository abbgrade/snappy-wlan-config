package config

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/funkygao/golib/observer"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Controller struct {
	Model             Config
	InterfacesDirPath string `yaml:"-"`
	ConfigPath        string `yaml:"-"`
}

func (config *Controller) Save() {

	// dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// write the file
	if err := ioutil.WriteFile(config.ConfigPath, data, 0644); err != nil {
		Warning.Fatalf("write %v : %v", config.ConfigPath, err)
	}
}

func (config *Controller) Load() {

	// does the file exist?
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		return
	}

	// read the file
	data, err := ioutil.ReadFile(config.ConfigPath)
	if err != nil {
		Warning.Fatalf("load: %v", err)
	}

	Trace.Printf("loaded %v", string(data))

	// parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}

}

func (config *Controller) Merge(request Transaction) {

	// configure observer
	observation := make(chan interface{})
	observationActive := true
	observer.Subscribe(EVENT_DELETE_INTERFACE, observation)

	// event handler
	go func() {
		for observationActive == true {
			interfaceName := <-observation
			if interfaceName == nil {
				continue
			}

			Trace.Printf("deleted interface: %#v\n", interfaceName)
			config.DeleteInterface(interfaceName.(string))
		}
	}()

	// act
	config.Model.Merge(request.Config.Model)

	// clean up observer
	observer.UnSubscribe(EVENT_DELETE_INTERFACE, observation)
	observationActive = false

}

func (config *Controller) GetWifiConfigPath(interfaceName string) string {

	fileName := fmt.Sprintf("interface_%v.conf", interfaceName)
	return path.Join(config.InterfacesDirPath, fileName)

}

func (config *Controller) GetNetworkConfigPath(interfaceName string) string {

	return fmt.Sprintf("/etc/network/interfaces.d/%v", interfaceName)

}

func (config *Controller) Export() {

	// sort networks by the interface
	interfaces := make(map[string][]WifiConfig)

	for _, network := range config.Model.Networks {

		// set default interface
		if network.Interface == "" {
			network.Interface = defaultInterface
		}

		// add network to the networks of the same interface
		interfaces[network.Interface] = append(interfaces[network.Interface], network)
	}

	for interfaceName, networks := range interfaces {

		// create a config file for each wifi interface
		wifiConfigPath := config.GetWifiConfigPath(interfaceName)
		Trace.Printf("export: %v", wifiConfigPath)
		wifiConfigFile, err := os.Create(wifiConfigPath)
		if err != nil {
			Warning.Fatalf("export %v : %v", wifiConfigPath, err)
		}
		defer wifiConfigFile.Close()

		// add a file header
		fmt.Fprint(wifiConfigFile, "# DO NOT CHANGE THIS FILE\n")
		fmt.Fprint(wifiConfigFile, "# This file is generated by snappy-wlan-config,\n")
		fmt.Fprint(wifiConfigFile, "# manual changes may become reversed.\n")
		fmt.Fprint(wifiConfigFile, "\n")

		for _, network := range networks {

			// add each network configuration
			export := NewWifiExport(&network)
			fmt.Fprint(wifiConfigFile, export.Dump())
			fmt.Fprint(wifiConfigFile, "\n")
			wifiConfigFile.Sync()
		}

		// create a config file for each interface
		networkConfigPath := config.GetNetworkConfigPath(interfaceName)
		Trace.Printf("create %v", networkConfigPath)

		inetConfigFile, err := os.Create(networkConfigPath)
		if err != nil {
			Warning.Fatalf("create %v : %v", networkConfigPath, err)
		}
		defer inetConfigFile.Close()

		// export networks
		networkExports := mapset.NewSet()
		for _, network := range networks {

			export := NewNetworkExport(&network)
			networkExports.Add(export.Dump())

		}

		// add a file header
		fmt.Fprint(inetConfigFile, "# DO NOT CHANGE THIS FILE\n")
		fmt.Fprint(inetConfigFile, "# This file is generated by snappy-wlan-config,\n")
		fmt.Fprint(inetConfigFile, "# manual changes may become reversed.\n")
		fmt.Fprint(inetConfigFile, "\n")

		// add each unique export
		for network := range networkExports.Iter() {

			fmt.Fprint(inetConfigFile, network.(string))
			fmt.Fprint(inetConfigFile, "\n")

		}

	}
}

func (config *Controller) DeleteInterface(interfaceName string) {

	interfacePath := config.GetWifiConfigPath(interfaceName)
	if err := os.Remove(interfacePath); err != nil {
		Warning.Fatalf("delete %v : %v", interfacePath, err)
	}
	Info.Printf("deleted: %v\n", interfacePath)

}
