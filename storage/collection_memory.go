package storage

import "sync"

type Collection struct {
	NameCollection string
	Collection     []WordPair
}
type WordPair struct {
	TextInput1 string
	TextInput2 string
}
type LearningStorage interface {
	AddCollection(userID int64, collection Collection)
	ListCollections(userID int64) []Collection
	DeleteCollection(userID int64, name string)
	ResetCollections(userID int64)
}
type MemoryLearningStorage struct {
	mu          sync.RWMutex
	collections map[int64][]Collection
}

func NewMemoryLearningStorage() *MemoryLearningStorage {
	return &MemoryLearningStorage{
		collections: make(map[int64][]Collection),
	}
}
func (m *MemoryLearningStorage) AddCollection(userID int64, collection Collection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.collections[userID] = append(m.collections[userID], collection)
}
func (m *MemoryLearningStorage) ListCollections(userID int64) []Collection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.collections[userID]
}
func (m *MemoryLearningStorage) DeleteCollection(userID int64, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var filtered []Collection
	for _, c := range m.collections[userID] {
		if c.NameCollection != name {
			filtered = append(filtered, c)
		}
	}
	m.collections[userID] = filtered
}
func (m *MemoryLearningStorage) ResetCollections(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.collections, userID)
}
