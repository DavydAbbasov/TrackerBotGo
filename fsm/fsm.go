package fsm

import (
	"sync"
)

type FSM struct {
	states map[int64]*UserState
	mu     sync.RWMutex
}

func NewFSM() *FSM {
	return &FSM{
		states: make(map[int64]*UserState),
	}
}

func (f *FSM) Get(userID int64) *UserState {
	f.mu.RLock()
	defer f.mu.RUnlock()

	state, ok := f.states[userID]
	if !ok {
		state = &UserState{
			Data: make(map[string]string),
		}
		f.states[userID] = state
	}
	return state
}

// Устанавливает (или создаёт, если нет) состояние для конкретного пользователя.
func (f *FSM) Set(userID int64, state string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	s, ok := f.states[userID]
	if !ok {
		s = &UserState{
			Data: make(map[string]string),
		}
		f.states[userID] = s
	}
	s.State = state
}

// Возвращает текущую структуру состояния (*UserState) по userID.
func (f *FSM) SetData(userID int64, key, value string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	s := f.Get(userID)
	s.Data[key] = value
}

func (f *FSM) GetData(userID int64, key string) string {
	s := f.Get(userID)
	return s.Data[key]
}

func (f *FSM) Reset(userID int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.states, userID)
}
