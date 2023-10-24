package gotesting

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

var assert = newAsserter

type Assert interface {
	// Equal asserts that the actual value is equal to the expected value.
	// If the values are not equal, it panics.
	Equal(actual, expected any)
}

type asserter struct {
	t *t
}

// newAsserter returns a new asserter.
func newAsserter(t *t) Assert {
	return &asserter{
		t: t,
	}
}

func (a *asserter) Equal(actual, expected any) {
	_, file, line, _ := runtime.Caller(1)
	fn := func() error {
		if reflect.DeepEqual(actual, expected) {
			return nil
		}

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

		json, _ := json.MarshalIndent(error, "", "  ")

		return fmt.Errorf("%s", json)
	}

	a.t.runs = append(a.t.runs, fn)
}
