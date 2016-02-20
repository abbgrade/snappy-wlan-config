package config

import (
	"fmt"
	"strings"
)

type Export struct {
	_keyValueFormat string
	_prefix         string
	_suffix         string
	Lines           []string
}

func (export *Export) Append(key, value string, optional bool, defaults ...string) {

	// apply defaults
	if value == "" && len(defaults) > 0 {
		value = defaults[0]
	}

	if value != "" {

		// export key value pair
		export.Lines = append(export.Lines, fmt.Sprintf(export._keyValueFormat, key, value))

	} else if optional == false {

		// fail on missing non optional value
		Warning.Fatalf("%v is required but not set", key)
	}
}

func (export *Export) Extend(lines ...string) {
	export.Lines = append(export.Lines, lines...)
}

func (export *Export) Dump() string {

	// wrap with prefix and suffix
	if export._prefix != "" {
		export.Lines = append([]string{export._prefix}, export.Lines...)
	}

	if export._suffix != "" {
		export.Lines = append(export.Lines, export._suffix)
	}

	return strings.Join(export.Lines, "\n")

}
