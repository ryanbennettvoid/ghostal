package json_file_config

import (
	"encoding/json"
	"errors"
	"ghostel/pkg/values"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// FileConfig implements IConfig interface using a JSON file
type FileConfig struct {
	mutex sync.Mutex // ensures atomic writes; protects the following fields
	data  map[string]string
}

// CreateJSONFileConfig initializes and returns a new FileConfig instance
func CreateJSONFileConfig() (*FileConfig, error) {
	fc := &FileConfig{
		data: make(map[string]string),
	}

	// Load existing data from the file
	err := fc.load()
	if err != nil {
		return nil, err
	}

	return fc, nil
}

// Get retrieves a value for a given key from the config
func (fc *FileConfig) Get(key string) (string, error) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	value, ok := fc.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

// Set updates or adds a new key-value pair in the config
func (fc *FileConfig) Set(key, value string) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.data[key] = value
	return fc.save()
}

// load reads the configuration from the JSON file
func (fc *FileConfig) load() error {
	file, err := os.Open(values.ConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // It's okay if the file doesn't exist
		}
		return err
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileData, &fc.data)
	if err != nil {
		return err
	}
	return nil
}

// save writes the current in-memory configuration to the JSON file
func (fc *FileConfig) save() error {
	file, err := json.MarshalIndent(fc.data, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(values.ConfigFilename, file, 0644)
}
