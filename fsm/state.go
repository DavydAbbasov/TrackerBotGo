package fsm

type UserState struct {
	State    string
	SubState string
	Data     map[string]string
}

