package helpers

func StringOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
