package lib

// Contains returns true if s is in strings array
func Contains(strings []string, s string) bool {
	for _, line := range strings {
		if line == s {
			return true
		}
	}
	return false
}
