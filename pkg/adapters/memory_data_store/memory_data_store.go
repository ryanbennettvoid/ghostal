package memory_data_store

type MemoryDataStore struct {
	data []byte
}

func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		data: make([]byte, 0),
	}
}

func (m *MemoryDataStore) Load() ([]byte, error) {
	return m.data, nil
}

func (m *MemoryDataStore) Save(bytes []byte) error {
	m.data = bytes
	return nil
}
