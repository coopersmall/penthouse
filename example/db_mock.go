package example

import (
	. "github.com/coopersmall/penthouse"
)

type DBMock struct {
	Mock
}

func NewDBMock() *DBMock {
	return &DBMock{
		Mock: NewMock(),
	}
}

func (m *DBMock) Get(id int) (string, error) {
	rets := m.CallMethod("Get", id)
	if rets[1] != nil {
		return "", rets[1].(error)
	}

	return rets[0].(string), nil
}

func (m *DBMock) Set(id int, value string) error {
	rets := m.CallMethod("Set", id, value)
	return rets[0].(error)
}
