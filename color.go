package gotesting

import (
	"github.com/fatih/color"
)

func Green(str string) string {
	return color.GreenString(str)
}

func Red(str string) string {
	return color.RedString(str)
}

func Yellow(str string) string {
	return color.YellowString(str)
}
