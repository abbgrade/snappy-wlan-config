package config

import (
	"fmt"
	"os"
	"path"
)

type ExportFile struct {
	Export
	_file *os.File
}

func OpenExportFile(filePath string) ExportFile {

	export := ExportFile{}

	Trace.Printf("export: %v", filePath)
	if _dryRunPath != "" {
		filePath = path.Join(_dryRunPath, filePath)
		os.MkdirAll(path.Dir(filePath), 0777)
	}

	if file, err := os.Create(filePath); err != nil {
		Warning.Fatalf("export %v : %v", filePath, err)
	} else {
		export._file = file
	}

	// check write access
	if _, err := fmt.Fprint(export._file, ""); err != nil {
		Warning.Fatalf("%v : %v", filePath, err)
	}

	return export
}

func (export *ExportFile) Flush() {
	export.Extend("", "")
	export._file.WriteString(export.Dump())
	export._file.Sync()
	export.Truncate()
}

func (export *ExportFile) Close() {
	export._file.Close()
}

func (export *ExportFile) AddHeader(toolName string) {

	export.Extend(""+
		"##", ""+
		"# DO NOT CHANGE THIS FILE", ""+
		fmt.Sprintf("## %v config", toolName), ""+
		"# This config file is generated by snappy-wlan-config,", ""+
		"# manual changes may become reversed.", ""+
		"##", "")

}