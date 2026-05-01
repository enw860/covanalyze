package enricher

func typeSwitch(i interface{}) string {
	switch v := i.(type) {
	case int:
		return "int"
	case string:
		return "string"
	default:
		_ = v
		return "unknown"
	}
}
