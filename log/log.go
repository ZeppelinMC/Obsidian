package log

import (
	"fmt"
	"obsidian/log/color"
)

var blue = color.Color{color.FgBlue, color.Bold}
var red = color.Color{color.FgRed, color.Bold}

func Info(a ...any) {
	str := fmt.Sprint(a...)

	fmt.Printf("%s: %s\n", blue.Colorize("INFO"), str)
}

func Infof(format string, a ...any) {
	str := fmt.Sprintf(format, a...)

	fmt.Printf("%s: %s\n", blue.Colorize("INFO"), str)
}

func Error(a ...any) {
	str := fmt.Sprint(a...)

	fmt.Printf("%s: %s\n", red.Colorize("ERROR"), str)
}

func Errorf(format string, a ...any) {
	str := fmt.Sprintf(format, a...)

	fmt.Printf("%s: %s\n", red.Colorize("ERROR"), str)
}
