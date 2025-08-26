package interfaces

import "github.com/DavydAbbasov/trecker_bot/fsm"

type FSMManager interface {
	Get(userID int64) *fsm.UserState
	Set(userID int64, state string)
	SetData(userID int64, key, value string)
	GetData(userID int64, key string) string
	Reset(userID int64)
}
