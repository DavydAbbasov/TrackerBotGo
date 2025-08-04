// реализация FSM в памяти (временная, простая)
package fsm

import "sync"

type memoryFSM struct {
	states map[int64]*UserState
	mu     sync.RWMutex
}

var memory = make(map[int64]*UserState)

func NewMemoryFSM() *FSM {
	return &FSM{
		states: make(map[int64]*UserState),
	}
}

func (m *memoryFSM) Get(userID int64) *UserState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.states[userID]
}

func (m *memoryFSM) Set(userID int64, state string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.states[userID] == nil {
		m.states[userID] = &UserState{Data: map[string]string{}}
	}
	m.states[userID].State = state
}

func (m *memoryFSM) SetData(userID int64, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.states[userID] == nil {
		m.states[userID] = &UserState{Data: map[string]string{}}
	}
	m.states[userID].Data[key] = value
}

func (m *memoryFSM) GetData(userID int64, key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.states[userID] == nil {
		return ""
	}
	return m.states[userID].Data[key]
}

func (m *memoryFSM) Reset(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, userID)
}
