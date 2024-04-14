package file_data_store

import (
	"io/ioutil"
	"os"
)

type FileDataStore struct {
	filepath string
}

func NewFileDataStore(filepath string) *FileDataStore {
	return &FileDataStore{filepath: filepath}
}

func (d *FileDataStore) Load() ([]byte, error) {
	data, err := ioutil.ReadFile(d.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]byte, 0), nil
		}
		return nil, err
	}
	return data, nil
}

func (d *FileDataStore) Save(data []byte) error {
	return ioutil.WriteFile(d.filepath, data, 0644)
}
