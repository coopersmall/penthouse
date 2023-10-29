package suite

import (
	"fmt"
	"testing"
)

type Output interface {
	Log(string, *testing.T)
	Error(string, *testing.T)
	Skip(string, *testing.T)
}

var NewOutput = newOutput

type output struct{}

func newOutput() Output {
	return &output{}
}

func (o *output) Log(title string, t *testing.T) {
	t.Logf(title)
	fmt.Print(title)
}

func (o *output) Error(message string, t *testing.T) {
	t.Errorf(message)
}

func (o *output) Skip(message string, t *testing.T) {
	t.Skip(message)
	fmt.Print(message)
}
