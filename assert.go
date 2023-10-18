package gotesting

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
)

var assert = newAsserter

type Assert interface {
	Equal(expected, actual any)
}

type asserter struct {
	t *T
}

// newAsserter returns a new asserter.
func newAsserter(t *T) Assert {
	return &asserter{
		t: t,
	}
}

func (a *asserter) Equal(expected, actual any) {
	fn := func(t *testing.T) error {
		var err error
		t.Run(t.Name(), func(t *testing.T) {
			if expected == actual {
				return
			}
			_, file, line, _ := runtime.Caller(1)

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

			err = fmt.Errorf("%s", json)
		})
		return err
	}

	a.t.runs = append(a.t.runs, fn)
}
