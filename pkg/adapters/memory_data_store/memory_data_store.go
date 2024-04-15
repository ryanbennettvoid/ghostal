package memory_data_store

type MemoryDataStore struct {
	Data []byte
}

func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		Data: make([]byte, 0),
	}
}

func (m *MemoryDataStore) Load() ([]byte, error) {
	return m.Data, nil
}

func (m *MemoryDataStore) Save(bytes []byte) error {
	m.Data = bytes
	return nil
}
