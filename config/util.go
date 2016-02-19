package config

func StringCoalesce(values ...string) string {

	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
