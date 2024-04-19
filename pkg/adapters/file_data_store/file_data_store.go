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

func (d *FileDataStore) loadFilepath(fp string) ([]byte, error) {
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *FileDataStore) Load() ([]byte, error) {
	filepath := d.filepath
	for depth := 0; depth < values.ConfigScanClimbMaxDepth; depth++ {
		data, err := d.loadFilepath(filepath)
		if err != nil {
			if os.IsNotExist(err) {
				filepath = path.Join("../", filepath)
				continue
			}
			return nil, err
		}
		return data, nil
	}
	// file not found after climbing parent directories
	return make([]byte, 0), nil
}

func (d *FileDataStore) Save(data []byte) error {
	return os.WriteFile(d.filepath, data, 0644)
}
