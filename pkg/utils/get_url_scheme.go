package utils

import "net/url"

func GetURLScheme(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return u.Scheme, nil
}
