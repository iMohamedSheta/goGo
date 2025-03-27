package enums

type Color string

// ANSI Color Codes
const (
	Reset  Color = "\033[0m"
	Blue   Color = "\033[1;34m"
	Green  Color = "\033[1;32m"
	Yellow Color = "\033[1;33m"
	Red    Color = "\033[1;31m"
)

func (c Color) Value() string {
	return string(c)
}
