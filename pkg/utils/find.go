package utils

import "errors"

func Find[T any](list []T, filter func(t T) bool) (T, error) {
	for _, item := range list {
		if filter(item) {
			return item, nil
		}
	}
	var t T
	return t, errors.New("not found")
}
