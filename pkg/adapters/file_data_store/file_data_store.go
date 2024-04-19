package file_data_store

import (
	"ghostal/pkg/values"
	"os"
	"path"
)

type FileDataStore struct {
	filepath string
}

func NewFileDataStore(filepath string) *FileDataStore {
	return &FileDataStore{filepath: filepath}
}

func (d *FileDataStore) resolveFilepath() (string, error) {
	filepath := d.filepath
	for depth := 0; depth < values.ConfigScanClimbMaxDepth; depth++ {
		_, err := os.Stat(filepath)
		if err != nil {
			if os.IsNotExist(err) {
				filepath = path.Join("../", filepath)
				continue
			}
			return "", err
		}
		return filepath, nil
	}
	// if not found, use the original filepath
	return d.filepath, nil
}

func (d *FileDataStore) Load() ([]byte, error) {
	filepath, err := d.resolveFilepath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// if not exists, return empty data
			return make([]byte, 0), nil
		}
		// if any other type of error, return error
		return nil, err
	}
	return data, nil
}

func (d *FileDataStore) Save(data []byte) error {
	filepath, err := d.resolveFilepath()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}
