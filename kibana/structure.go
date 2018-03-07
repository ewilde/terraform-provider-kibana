package kibana

func boolOrDefault(value interface{}, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}

	return value.(bool)
}

func stringOrDefault(value interface{}, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	return value.(string)
}

type stringApply func(string)

func stringApplyIfExists(value interface{}, apply stringApply) {
	if value == nil {
		return
	}

	apply(value.(string))
}
