package config

import (
	"io/ioutil"
	"os"

	"github.com/funkygao/golib/observer"
	"gopkg.in/yaml.v2"
)

var (
	_dryRun            = false
	_interfacesDirPath string
	_configPath        string
)

type Controller struct {
	Model Config
}

func InitController(interfacesDirPath, configPath string, dryRun bool) Controller {
	_interfacesDirPath = interfacesDirPath
	_configPath = configPath
	_dryRun = dryRun

	config := Controller{}
	return config
}

func (config *Controller) Save() {

	// dump the YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		Warning.Fatalf("dump: %v", err)
	}

	// write the file
	if err := ioutil.WriteFile(_configPath, data, 0644); err != nil {
		Warning.Fatalf("write %v : %v", _configPath, err)
	}

	Trace.Print("saved: %v", config)
}

func (config *Controller) Load() {

	// does the file exist?
	if _, err := os.Stat(_configPath); os.IsNotExist(err) {
		return
	}

	// read the file
	data, err := ioutil.ReadFile(_configPath)
	if err != nil {
		Warning.Fatalf("load: %v", err)
	}

	Trace.Printf("loaded %v", string(data))

	// parse the YAML
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		Warning.Fatalf("parse: %v", err)
	}

	Trace.Printf("loaded: %v from %v", config, _configPath)
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

			Trace.Printf("remove interface: %#v\n", interfaceName)
			config.RemoveInterface(interfaceName.(string))
		}
	}()

	// act
	config.Model.Merge(request.Config.Model)

	// clean up observer
	observer.UnSubscribe(EVENT_DELETE_INTERFACE, observation)
	observationActive = false

}

func (config *Controller) Export() {

	// sort networks by the interface
	interfaces := make(map[string][]WifiConfig)

	for _, network := range config.Model.Networks {

		// set default interface
		network.Interface = StringCoalesce(network.Interface, INTERFACE_DEFAULT)

		// add network to the networks of the same interface
		interfaces[network.Interface] = append(interfaces[network.Interface], network)
	}

	for interfaceName, networks := range interfaces {

		// create a config file for each wifi interface
		config.ExportWifiClient(interfaceName, networks)

		// create a config file for accesspoints
		config.ExportWifiAccesspoint(networks)

		// create a config file for each interface
		config.ExportInterface(interfaceName, networks)
	}
}

func (config *Controller) RemoveInterface(interfaceName string) {

	interfacePath := GetWifiConfigPath(interfaceName)
	if _, err := os.Stat(interfacePath); os.IsNotExist(err) {

	} else if err := os.Remove(interfacePath); err != nil {
		Warning.Fatalf("delete %v : %v", interfacePath, err)
	}
	Info.Printf("deleted: %v\n", interfacePath)

}
