package helpers

func IsModelsContain(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}

	return false
}
