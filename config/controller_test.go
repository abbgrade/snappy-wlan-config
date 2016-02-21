package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	//"regexp"
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

	wifiConfigDirPath := tempPath
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

	assert.Equal(t, len(files), 1, "wrong number of generated interface files")
	assert.Equal(t, files[0].Name(), interfaceName, "wrong interface file generated")

	interfaceFilePath := path.Join(interfacesDirPath, interfaceName)
	{
		data, _ := ioutil.ReadFile(interfaceFilePath)

		expectedLine := fmt.Sprintf("iface %v inet dhcp\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), string(data), expectedLine)

		expectedLine = fmt.Sprintf("allow-hotplug %v\n", interfaceName)
		assert.True(t, strings.Contains(string(data), expectedLine), string(data), expectedLine)

		expectedLine = "wpa-conf "
		assert.True(t, strings.Contains(string(data), expectedLine), string(data), expectedLine)

		/*
			re := regexp.MustCompile("wpa-conf (.+)?")
			wpaFilePath := re.FindString(string(data))
			assert.True(t, strings.Contains(string(data), expectedLine), wpaFilePath, expectedLine)

			wpaFilePath = strings.Replace(wpaFilePath, expectedLine, "", 1)

			fmt.Printf("%v\n", wpaFilePath)

			if _, err := ioutil.ReadFile(wpaFilePath); err != nil {
				t.Fatal("missing referenced wpa file")
			}
		*/

	}

	Info.Println(interfacesDirPath)
	Info.Println(files)
}
