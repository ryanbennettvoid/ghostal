package utils

import "errors"

var FindNotFoundError = errors.New("not found")

func Find[T any](list []T, filter func(t T) bool) (T, error) {
	for _, item := range list {
		if filter(item) {
			return item, nil
		}
	}
	var t T
	return t, FindNotFoundError
}
