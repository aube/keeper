package common

func StringOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
