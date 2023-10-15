package gotesting

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
)

type Assert interface {
	Equal(expected, actual any)
	NotEqual(expected, actual any) bool
	True(value bool) bool
	False(value bool) bool
}

type asserter struct {
	name string
	t    *testing.T
}

// newAsserter returns a new asserter.
func newAsserter(name string, t *testing.T) Assert {
	return &asserter{
		name: name,
		t:    t,
	}
}

func (a *asserter) Equal(expected, actual any) {
	_, file, line, _ := runtime.Caller(1)
	a.t.Run(a.name, func(t *testing.T) {
		if expected == actual {
			fmt.Print(Green("•"))
			return
		}
		fmt.Print(Red("•"))

		error := struct {
			Expected any    `json:"expected"`
			Actual   any    `json:"actual"`
			File     string `json:"file"`
			Line     int    `json:"line"`
		}{
			Expected: expected,
			Actual:   actual,
			File:     file,
			Line:     line,
		}

		json, _ := json.Marshal(error)

		t.Errorf("\n%s\n", json)
	})
}

func (a *asserter) NotEqual(expected, actual any) bool {
	_, file, line, _ := runtime.Caller(1)
	if expected != actual {
		return true
	}

	a.t.Errorf("%s:%d: expected %v, got %v", file, line, expected, actual)
	return false
}

func (a *asserter) True(value bool) bool {
	_, file, line, _ := runtime.Caller(1)
	if value {
		return true
	}

	a.t.Errorf("%s:%d: expected true, got false", file, line)
	return false
}

func (a *asserter) False(value bool) bool {
	_, file, line, _ := runtime.Caller(1)
	if !value {
		return true
	}

	a.t.Errorf("%s:%d: expected false, got true", file, line)
	return false
}
