package utils

import (
	"strings"
)

func ResolveCollection(path string) string {
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if segment != "" {
			return segment
		}
	}

	return ""
}