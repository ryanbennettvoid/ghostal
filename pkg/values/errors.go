package values

import "errors"

var NoProgramArgsProvidedError = errors.New("no program arguments provided")
var NoSpecialCharsErr = errors.New("snapshot name can only contain alphanumeric characters and underscores")
var SnapshotNameTakenErr = errors.New("snapshot name already used")
var SnapshotNotExistsErr = errors.New("snapshot does not exist")
var UnsupportedURLSchemeError = errors.New("url scheme is unsupported")
