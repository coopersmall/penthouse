package gotesting

const (
	escape = "\x1b"
	red    = "31"
	green  = "32"
	yellow = "33"
	orange = "38;5;208"
	cyan   = "36"
)

func Green(str string) string {
	return addColor(str, green)
}

func Red(str string) string {
	return addColor(str, red)
}

func Yellow(str string) string {
	return addColor(str, yellow)
}

func Orange(str string) string {
	return addColor(str, orange)
}

func Cyan(str string) string {
	return addColor(str, cyan)
}

func addColor(str string, color string) string {
	return escape + "[" + color + "m" + str + escape + "[0m"
}
