package store

// Use memory as storage, use with caution in build environment
type MemoryDAO struct {
	store map[string]KeyInfo
	Crypto
}

func NewMemory(crypto Crypto) MemoryDAO {
	if crypto == nil {
		crypto = AES{}
	}
	return MemoryDAO{
		store:  make(map[string]KeyInfo),
		Crypto: crypto,
	}
}
func (m MemoryDAO) Write(name, password string, store KeyInfo) error {
	m.store[name] = store
	return nil
}

func (m MemoryDAO) Read(name, password string) (KeyInfo, error) {
	return m.store[name], nil
}

// ReadMetadata read a key information from the local store
func (m MemoryDAO) ReadMetadata(name string) (store KeyInfo, err error) {
	return m.store[name], nil
}

func (m MemoryDAO) Delete(name, password string) error {
	delete(m.store, name)
	return nil
}

func (m MemoryDAO) Has(name string) bool {
	_, ok := m.store[name]
	return ok
}
