package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	controller Controller
	tempPath   string
)

func setupTest(test_name string) {
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	tempPath, _ = ioutil.TempDir("", "")

	Info.Println("create", tempPath)

	wifiConfigDirPath := "/"
	configPath := path.Join(tempPath, "config.yaml")
	dryRunPath := tempPath

	controller = InitController(wifiConfigDirPath, configPath, dryRunPath)

	Info.Println("Test Controller:", test_name)

}

func tearDownTest() {
	Info.Println("remove", tempPath)
	os.RemoveAll(tempPath)
}

func TestExportInterface(t *testing.T) {
	setupTest("export interface")
	defer tearDownTest()

	interfaceName := "wlan42"

	wifi := WifiConfig{}
	wifi.Interface = interfaceName
	networks := []WifiConfig{wifi}

	controller.ExportInterface(interfaceName, networks)

	interfacesDirPath := path.Join(tempPath, "/etc/network/interfaces.d/")
	files, _ := ioutil.ReadDir(interfacesDirPath)

	assert.Equal(t, 1, len(files), "wrong number of generated interface files")
	assert.Equal(t, interfaceName, files[0].Name(), "wrong interface file generated")

	interfaceFilePath := path.Join(interfacesDirPath, interfaceName)
	{
		data, err := ioutil.ReadFile(interfaceFilePath)
		assert.Nil(t, err)
		//Info.Println("interface config\n", string(data))

		expectedLine := fmt.Sprintf("iface %v inet dhcp\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		expectedLine = fmt.Sprintf("allow-hotplug %v\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		expectedLine = "wpa-conf "
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		re := regexp.MustCompile("wpa-conf (.+)?")
		wpaFilePath := re.FindString(string(data))
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, wpaFilePath)

		wpaFilePath = strings.Replace(wpaFilePath, expectedLine, "", 1)
		wpaFilePath = strings.Replace(wpaFilePath, tempPath, "", 1)

		assert.Equal(t, fmt.Sprintf("/interface_%v.conf", interfaceName), wpaFilePath, "wrong wpa-conf path")
	}

}

func TestExportInterface_static(t *testing.T) {
	setupTest("export interface static")
	defer tearDownTest()

	interfaceName := "wlan42"

	wifi := WifiConfig{}
	wifi.Interface = interfaceName
	wifi.IPConfig.AddressType = "static"
	wifi.IPConfig.Address = "192.168.42.2"
	wifi.IPConfig.Netmask = "255.255.255.0"
	wifi.IPConfig.Network = "192.168.42.0"
	wifi.IPConfig.Gateway = "192.168.42.1"
	networks := []WifiConfig{wifi}

	controller.ExportInterface(interfaceName, networks)

	interfacesDirPath := path.Join(tempPath, "/etc/network/interfaces.d/")
	files, _ := ioutil.ReadDir(interfacesDirPath)

	assert.Equal(t, 1, len(files), "wrong number of generated interface files")
	assert.Equal(t, interfaceName, files[0].Name(), "wrong interface file generated")

	interfaceFilePath := path.Join(interfacesDirPath, interfaceName)
	{
		data, err := ioutil.ReadFile(interfaceFilePath)
		assert.Nil(t, err)
		//Info.Println("interface config\n", string(data))

		expectedLine := fmt.Sprintf("iface %v inet static\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		expectedLine = fmt.Sprintf("allow-hotplug %v\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		expectedLine = fmt.Sprintf("address %v\n", wifi.IPConfig.Address)
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		expectedLine = "wpa-conf "
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, string(data))

		re := regexp.MustCompile("wpa-conf (.+)?")
		wpaFilePath := re.FindString(string(data))
		assert.True(t, strings.Contains(string(data), expectedLine), expectedLine, wpaFilePath)

		wpaFilePath = strings.Replace(wpaFilePath, expectedLine, "", 1)
		wpaFilePath = strings.Replace(wpaFilePath, tempPath, "", 1)

		assert.Equal(t, fmt.Sprintf("/interface_%v.conf", interfaceName), wpaFilePath, "wrong wpa-conf path")
	}

}

func TestExportWifiClient(t *testing.T) {
	setupTest("export wifi client")
	defer tearDownTest()

	interfaceName := "wlan42"
	ssid := "FancyWIFI"
	psk := "secret"

	wifi := WifiConfig{}
	wifi.Interface = interfaceName
	wifi.SSID = ssid
	wifi.PSK = psk
	networks := []WifiConfig{wifi}

	controller.ExportWifiClient(interfaceName, networks)

	wifiDirPath := tempPath

	Info.Println(wifiDirPath)

	files, _ := ioutil.ReadDir(wifiDirPath)

	for _, file := range files {
		Info.Println(file.Name())
	}

	assert.Equal(t, 1, len(files), "wrong number of generated interface files")
	wifiConfigName := fmt.Sprintf("interface_%v.conf", interfaceName)
	assert.Equal(t, wifiConfigName, files[0].Name(), "wrong wifi file generated")
	wifiConfigPath := path.Join(wifiDirPath, wifiConfigName)
	{
		data, err := ioutil.ReadFile(wifiConfigPath)
		assert.Nil(t, err)
		//Info.Println("wifi config\n", string(data))

		re := regexp.MustCompile("(?s)network={(.+)")
		content := re.FindString(string(data))

		assert.NotEqual(t, "", content, "must contain \"network={.*}\"", string(data))

		expectedLine := fmt.Sprintf("ssid=\"%v\"\n", ssid)
		assert.True(t, strings.Contains(content, expectedLine), expectedLine, content)

		expectedLine = fmt.Sprintf("psk=\"%v\"\n", psk)
		assert.True(t, strings.Contains(content, expectedLine), expectedLine, content)

	}

}
