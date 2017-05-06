package utils

// InArray returns true if the string 's' is found in the array 'arr', otherwise false
func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
