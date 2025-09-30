package storage

import (
	"sync"
	"time"
)

type Activity struct {
	ID     int64
	UserID int64
	// Name         string
	// Emoji        string
	// IsArchived   bool
	// CreatedAt    time.Time
	NameActivity string
	TimeEntry    []TimeEntry
}

type TimeEntry struct {
	Timestamp time.Time
	Start     time.Time
	End       time.Time
	Duration  time.Duration
}

type ActivityStorage interface {
	Add(userID int64, activity Activity)
	List(userID int64) []Activity
	Reset(userID int64)
	Delete(userID int64, name string)
}

type MemoryActivityStorage struct {
	mu         sync.RWMutex
	activities map[int64][]Activity
}

// Конструктор
func NewMemoryActivityStorage() *MemoryActivityStorage {
	return &MemoryActivityStorage{
		activities: make(map[int64][]Activity),
	}
}

// Добавляет активность
func (m *MemoryActivityStorage) Add(userID int64, activity Activity) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activities[userID] = append(m.activities[userID], activity)
}

// Возвращает все активности пользователя
func (m *MemoryActivityStorage) List(userID int64) []Activity {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activities[userID]
}

// Удаляет все активности пользователя
func (m *MemoryActivityStorage) Reset(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activities, userID)
}

// Удаляет конкретную активность по имени
func (m *MemoryActivityStorage) Delete(userID int64, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var filtered []Activity
	for _, a := range m.activities[userID] {
		if a.NameActivity != name {
			filtered = append(filtered, a)
		}
	}
	m.activities[userID] = filtered
}
