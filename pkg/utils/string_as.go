package utils

import (
	"fmt"
)

func StringAsBool(str string) (bool, error) {
	switch str {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("failed to parse \"%s\" as bool", str)
	}
}
