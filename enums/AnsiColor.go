package enums

type AnsiColor string

// ANSI Color Codes
const (
	Reset  AnsiColor = "\033[0m"
	Blue   AnsiColor = "\033[1;34m"
	Green  AnsiColor = "\033[1;32m"
	Yellow AnsiColor = "\033[1;33m"
	Red    AnsiColor = "\033[1;31m"
)

func (c AnsiColor) Value() string {
	return string(c)
}
