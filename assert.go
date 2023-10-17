package gotesting

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Assert interface {
	Equal(expected, actual any)
	NotEqual(expected, actual any)
	True(value bool)
	False(value bool)
}

type asserter struct {
	name string
	t    *T
}

// newAsserter returns a new asserter.
func newAsserter(name string, t *T) Assert {
	return &asserter{
		name: name,
		t:    t,
	}
}

func (a *asserter) Equal(expected, actual any) {
	fn := func() error {
		if expected == actual {
			return nil
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

		return fmt.Errorf("%s", json)
	}

	a.t.runs = append(a.t.runs, fn)
}

func (a *asserter) NotEqual(expected, actual any) {
	fn := func() error {
		if expected != actual {
			return nil
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

		return fmt.Errorf("%s", json)
	}

	a.t.runs = append(a.t.runs, fn)
}

func (a *asserter) True(value bool) {
	fn := func() error {
		if value {
			return nil
		}
		_, file, line, _ := runtime.Caller(1)

		error := struct {
			Value bool   `json:"value"`
			File  string `json:"file"`
			Line  int    `json:"line"`
		}{
			Value: value,
			File:  file,
			Line:  line,
		}

		json, _ := json.Marshal(error)

		return fmt.Errorf("%s", json)
	}

	a.t.runs = append(a.t.runs, fn)
}

func (a *asserter) False(value bool) {
	fn := func() error {
		if !value {
			return nil
		}
		_, file, line, _ := runtime.Caller(1)

		error := struct {
			Value bool   `json:"value"`
			File  string `json:"file"`
			Line  int    `json:"line"`
		}{
			Value: value,
			File:  file,
			Line:  line,
		}

		json, _ := json.Marshal(error)

		return fmt.Errorf("%s", json)
	}

	a.t.runs = append(a.t.runs, fn)
}
