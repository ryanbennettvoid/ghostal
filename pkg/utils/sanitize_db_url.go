package utils

import "net/url"

func SanitizeDBURL(dbURL string) (string, error) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Redacted(), nil
}
