package event

import "time"

type UserCreated struct {
	Name    string
	Payload interface{}
}

func NewUserCreated() *UserCreated {
	return &UserCreated{
		Name: "UserCreated",
	}
}

func (e *UserCreated) GetName() string {
	return e.Name
}

func (e *UserCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *UserCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *UserCreated) GetDateTime() time.Time {
	return time.Now()
}