package utils

func ToPointer[T any](thing T) *T {
	return &thing
}
