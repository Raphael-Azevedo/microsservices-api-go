package event

import "time"

type FindUserByEmail struct {
	Name    string
	Payload interface{}
}

func NewFindUserByEmail() *FindUserByEmail {
	return &FindUserByEmail{
		Name: "FindUserByEmail",
	}
}

func (e *FindUserByEmail) GetName() string {
	return e.Name
}

func (e *FindUserByEmail) GetPayload() interface{} {
	return e.Payload
}

func (e *FindUserByEmail) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *FindUserByEmail) GetDateTime() time.Time {
	return time.Now()
}