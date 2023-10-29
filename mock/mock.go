package mock

import (
	"fmt"
)

// Mock is a mock object for testing.
// It is used to mock functions.
type Mock interface {
	// Call calls the function with the given name and arguments.
	// It returns the values that were set with SetReturn.
	// If no values were set, it panics.
	CallMethod(name string, args ...any) []any
	// ClearReturns clears the return values for the given function.
	// If no values were set, it does not perform an operation.
	ClearReturns(name string)
	// SetReturns sets the return values for the given function.
	// If values were already set, they are overwritten.
	SetReturns(name string, values ...any)
	// GetCallParams returns the arguments that were passed to the function.
	// If no values were set, it panics.
	CallParam(name string, call, position int) any
	// GetCallCount returns the number of times the function was called.
	// If no values were set, it panics.
	CallCount(name string) int
}

var NewMock = newMock

type method struct {
	args [][]any
	rets []any
}

type mock struct {
	methods map[string]method
}

func newMock() Mock {
	return &mock{
		methods: make(map[string]method),
	}
}

func (m *mock) CallMethod(name string, args ...any) []any {
	method, ok := m.methods[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	method.args = append(method.args, args)
	m.methods[name] = method
	return method.rets
}

func (m *mock) ClearReturns(name string) {
	mo, ok := m.methods[name]
	if !ok {
		return
	}
	mo.rets = make([]any, 0)
	m.methods[name] = mo
}

func (m *mock) SetReturns(name string, values ...any) {
	mo, ok := m.methods[name]
	if !ok {
		mo = method{
			args: make([][]any, 0),
			rets: make([]any, 0),
		}
	}
	mo.rets = values
	m.methods[name] = mo
}

func (m *mock) CallParam(name string, call, position int) any {
	ret, ok := m.methods[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	return ret.args[call][position]
}

func (m *mock) CallCount(name string) int {
	ret, ok := m.methods[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	return len(ret.args)
}

func Error(a any) error {
	if a == nil {
		return nil
	}

	return a.(error)
}
