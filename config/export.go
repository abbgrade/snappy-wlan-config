package config

import (
	"fmt"
	"strings"
)

type Export struct {
	_keyValueFormat string
	_prefix         string
	_suffix         string
	_lines          []string
}

func (export *Export) Append(key, value string, optional bool, defaults ...string) {

	// apply defaults
	if value == "" && len(defaults) > 0 {
		value = defaults[0]
	}

	if value != "" {

		// export key value pair
		export._lines = append(export._lines, fmt.Sprintf(export._keyValueFormat, key, value))

	} else if optional == false {

		// fail on missing non optional value
		Warning.Fatalf("%v is required but not set", key)
	}
}

func (export *Export) Extend(lines ...string) {
	export._lines = append(export._lines, lines...)
}

func (export *Export) Truncate() {
	export._lines = []string{}
}

func (export *Export) Dump() string {

	// wrap with prefix and suffix
	if export._prefix != "" {
		export._lines = append([]string{export._prefix}, export._lines...)
	}

	if export._suffix != "" {
		export._lines = append(export._lines, export._suffix)
	}

	return strings.Join(export._lines, "\n")

}
