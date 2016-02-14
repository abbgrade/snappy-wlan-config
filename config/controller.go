package config

import (
	"fmt"
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
	ioutil.WriteFile(config.ConfigPath, data, 0644)
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
	observation := make(chan interface{})
	observationActive := true
	observer.Subscribe(EVENT_DELETE_INTERFACE, observation)
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

	config.Model.Merge(request.Config.Model)

	observer.UnSubscribe(EVENT_DELETE_INTERFACE, observation)
	observationActive = false

}

func (config *Controller) GetInterfacePath(interfaceName string) string {
	fileName := fmt.Sprintf("interface_%v.conf", interfaceName)
	return path.Join(config.InterfacesDirPath, fileName)
}

func (config *Controller) Export() {

	// sort networks by the interface
	interfaces := make(map[string][]NetworkConfig)

	for _, network := range config.Model.Networks {

		// set default interface
		if network.Interface == "" {
			network.Interface = defaultInterface
		}

		// add network to the networks of the same interface
		interfaces[network.Interface] = append(interfaces[network.Interface], network)
	}

	for interfaceName, networks := range interfaces {

		// create a config file for each interface
		file, err := os.Create(config.GetInterfacePath(interfaceName))
		if err != nil {
			Warning.Fatalf("export: %v", err)
		}
		defer file.Close()

		// add a file header
		fmt.Fprintf(file, "# DO NOT CHANGE THIS FILE\n")
		fmt.Fprintf(file, "# This file is generated by snappy-wlan-config,\n")
		fmt.Fprintf(file, "# manual changes may become reversed.\n")
		fmt.Fprintf(file, "\n")

		for _, network := range networks {

			// add each network configuration
			export := NewNetworkExport(&network)
			export.Save(file)
			fmt.Fprintf(file, "\n")
			file.Sync()
		}

	}
}

func (config *Controller) DeleteInterface(interfaceName string) {
	interfacePath := config.GetInterfacePath(interfaceName)
	os.Remove(interfacePath)
	Info.Printf("deleted: %v\n", interfacePath)
}
