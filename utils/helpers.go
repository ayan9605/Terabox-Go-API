package utils

import (
	"fmt"
	"strings"
)

// GetSize converts bytes to human-readable format
func GetSize(sizeBytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case sizeBytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(sizeBytes)/float64(GB))
	case sizeBytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(sizeBytes)/float64(MB))
	case sizeBytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(sizeBytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", sizeBytes)
	}
}

// FindBetween extracts string between two delimiters
func FindBetween(str, start, end string) string {
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)

	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	return str[startIndex : startIndex+endIndex]
}
