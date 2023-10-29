package example

import (
    "github.com/coopersmall/penthouse/mock"
)

type DBMock struct {
    Mock mock.Mock
}

func NewDBMock() *DBMock {
    return &DBMock{
        Mock: mock.NewMock(),
    }
}


func (m *DBMock) Get(id int) (string, error) {
    args := m.Mock.CallMethod("Get", id)
    return args[0].(string), mock.Error(args[1])
}

func (m *DBMock) Set(id int, value string) (error) {
    args := m.Mock.CallMethod("Set", id, value)
    return mock.Error(args[0])
}
