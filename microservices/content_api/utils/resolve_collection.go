package utils

import "strings"

func ResolveCollection(path string) string {
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			return segment
		}
	}

	return ""
}