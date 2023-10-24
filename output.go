package gotesting

import (
	"fmt"
	"testing"
)

type Output interface {
	Log(string)
	Error(string)
	Skip(string)
}

var NewOutput = newOutput

type output struct {
	t *testing.T
}

func newOutput(t *testing.T) Output {
	return &output{
		t: t,
	}
}

func (o *output) Log(title string) {
	o.t.Logf(title)
	fmt.Print(title)
}

func (o *output) Error(message string) {
	o.t.Errorf(message)
	fmt.Print(message)
}

func (o *output) Skip(message string) {
	o.t.Skip(message)
	fmt.Print(message)
}
