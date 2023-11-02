package suite

import (
	"fmt"
	"strings"
)

func title(suite *suite) string {
	plural := "s"
	if len(suite.tests) == 1 {
		plural = ""
	}

	title := fmt.Sprintf("%s: %d test%s", cyan(suite.name), len(suite.tests), plural)

	if suite.focused {
		title = fmt.Sprintf("%s: %d test%s (focused)", yellow(suite.name), len(suite.tests), plural)
	}

	border := strings.Repeat("-", len(title))
	return border + "\n" + title + "\n" + border + "\n"
}

func success() string {
	return (green("•"))
}

func failure() string {
	return (red("•"))
}

func skip() string {
	return (yellow("•"))
}

const (
	escape       = "\x1b"
	redSymbol    = "31"
	greenSymbol  = "32"
	yellowSymbol = "33"
	orangeSymbol = "38;5;208"
	cyanSymbol   = "36"
)

func green(str string) string {
	return addColor(str, greenSymbol)
}

func red(str string) string {
	return addColor(str, redSymbol)
}

func yellow(str string) string {
	return addColor(str, yellowSymbol)
}

func orange(str string) string {
	return addColor(str, orangeSymbol)
}

func cyan(str string) string {
	return addColor(str, cyanSymbol)
}

func addColor(str string, color string) string {
	return escape + "[" + color + "m" + str + escape + "[0m"
}
