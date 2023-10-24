package gotesting

import (
	"encoding/json"
	"strings"
)

type Formatter interface {
	Title(string) string
	Focus(string) string
	Success() string
	Failure(...error) string
	Skip() string
}

var NewFormatter = newFormatter

type formatter struct{}

func newFormatter() Formatter {
	return &formatter{}
}

func (f *formatter) Title(title string) string {
	border := strings.Repeat("-", len(title))
	return border + "\n" + Cyan(title) + "\n" + border + "\n"
}

func (f *formatter) Focus(title string) string {
	border := strings.Repeat("-", len(title))
	return border + "\n" + Orange(title) + "\n" + border + "\n"
}

func (f *formatter) Success() string {
	return (Green("•"))
}

func (f *formatter) Failure(errs ...error) string {
	var s strings.Builder

	for i, err := range errs {
		if i == 0 {
			s.WriteString("\n")
		}

		json, _ := json.MarshalIndent(err, "", "  ")
		s.WriteString(string(json))
		if i != len(errs)-1 {
			s.WriteString(",")
		}

		s.WriteString("\n")
	}

	s.WriteString(Red("•"))

	return s.String()
}

func (f *formatter) Skip() string {
	return (Yellow("•"))
}
