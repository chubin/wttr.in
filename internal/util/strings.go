package util

import "strings"

func InSlice(what string, where []string) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}
	return false
}

func HasPrefixInSlice(what string, where []string, prefix string) bool {
	for _, item := range where {
		item = item + prefix
		if strings.HasPrefix(what, item) {
			return true
		}
	}
	return false
}
