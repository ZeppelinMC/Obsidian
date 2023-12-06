package color

import (
	"fmt"
	"strings"
)

type Attribute = string
type Color []Attribute

const (
	Reset         Attribute = "0"
	Bold          Attribute = "1"
	Dim           Attribute = "2"
	Italic        Attribute = "3"
	Underline     Attribute = "4"
	Blinking      Attribute = "5"
	Inverse       Attribute = "6"
	Hidden        Attribute = "7"
	Strikethrough Attribute = "8"

	FgBlack   Attribute = "30"
	FgRed     Attribute = "31"
	FgGreen   Attribute = "32"
	FgYellow  Attribute = "33"
	FgBlue    Attribute = "34"
	FgMagenta Attribute = "35"
	FgCyan    Attribute = "36"
	FgWhite   Attribute = "37"

	BgBlack   Attribute = "40"
	BgRed     Attribute = "41"
	BgGreen   Attribute = "42"
	BgYellow  Attribute = "43"
	BgBlue    Attribute = "44"
	BgMagenta Attribute = "45"
	BgCyan    Attribute = "46"
	BgWhite   Attribute = "47"

	FgBrBlack   Attribute = "90"
	FgBrRed     Attribute = "91"
	FgBrGreen   Attribute = "92"
	FgBrYellow  Attribute = "93"
	FgBrBlue    Attribute = "94"
	FgBrMagenta Attribute = "95"
	FgBrCyan    Attribute = "96"
	FgBrWhite   Attribute = "97"

	BgBrBlack   Attribute = "100"
	BgBrRed     Attribute = "101"
	BgBrGreen   Attribute = "102"
	BgBrYellow  Attribute = "103"
	BgBrBlue    Attribute = "104"
	BgBrMagenta Attribute = "105"
	BgBrCyan    Attribute = "106"
	BgBrWhite   Attribute = "107"
)

const escape = "\x1b"

func (c Color) Colorize(str string) string {
	return fmt.Sprintf("%s[%sm%s%s[0m", escape, strings.Join(c, ";"), str, escape)
}

func (c Color) colorize() string {
	return fmt.Sprintf("%s[%sm", escape, strings.Join(c, ";"))
}

func FromCode(str string) string {
	var s strings.Builder
	var cuco Color
	for i, c := range str {
		if c == '&' && len(str) > i {
			switch {
			case str[i+1] == 'r':
				clear(cuco)
			case isColor(str[i+1]):
				cuco = append(cuco, colors[str[i+1]])
			}
			continue
		}
		if isColor(byte(c)) && i != 0 && str[i-1] == '&' {
			continue
		}
		s.WriteString(cuco.colorize())
		s.WriteRune(c)
	}
	return s.String()
}

func isColor(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'k' && c <= 'o') || c == 'r'
}

var colors = map[byte]Attribute{
	'0': FgBlack,
	'1': FgBlue,
	'2': FgGreen,
	'3': FgCyan,
	'4': FgRed,
	'5': FgMagenta,
	'6': FgYellow,
	'7': FgWhite,
	'8': FgBrBlack,
	'9': FgBlue,
	'a': FgBrGreen,
	'b': FgBrCyan,
	'c': FgBrRed,
	'd': FgBrMagenta,
	'e': FgBrYellow,
	'f': FgBrWhite,

	'k': Hidden,
	'l': Bold,
	'm': Strikethrough,
	'n': Underline,
	'o': Italic,
}
