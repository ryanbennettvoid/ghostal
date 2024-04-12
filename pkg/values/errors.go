package values

import "errors"

var NoSpecialCharsErr = errors.New("snapshot name can only contain alphanumeric characters or underscores")
var SnapshotNameTakenErr = errors.New("snapshot name already used")
var SnapshotNotExistsErr = errors.New("snapshot does not exist")
