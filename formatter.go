package suite

import (
	"strings"
)

func title(title string) string {
	border := strings.Repeat("-", len(title))
	return border + "\n" + cyan(title) + "\n" + border + "\n"
}

func focusTitle(title string) string {
	border := strings.Repeat("-", len(title))
	return border + "\n" + orange(title) + "\n" + border + "\n"
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
