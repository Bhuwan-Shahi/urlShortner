package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP ensures the URL starts with http:// or https://
func EnforceHTTP(url string) string {
	// Check if URL starts with either http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// RemoveDomainError checks if the URL matches the domain
// Returns false if the URL matches our own domain to prevent self-referential shortcuts
func RemoveDomainError(url string) bool {
	// First check the raw URL
	if url == os.Getenv("DOMAIN") {
		return false
	}

	// Clean the URL by removing protocols and www
	cleanURL := url
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "www.")

	// Get the domain part (everything before the first slash)
	domain := strings.Split(cleanURL, "/")[0]

	// Compare with our domain
	if domain == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
