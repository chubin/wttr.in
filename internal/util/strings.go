package util

func InSlice(what string, where []string) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}
	return false
}
