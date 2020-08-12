package color

import "runtime"

var Reset = "\033[0m"
var RedColor = "\033[31m"
var GreenColor = "\033[32m"
var YellowColor = "\033[33m"
var BlueColor = "\033[34m"
var PurpleColor = "\033[35m"
var CyanColor = "\033[36m"
var GrayColor = "\033[37m"
var WhiteColor = "\033[97m"

func Red(s string) string {
	return Colored(RedColor, s)
}

func Green(s string) string {
	return Colored(GreenColor, s)
}

func Yellow(s string) string {
	return Colored(YellowColor, s)
}

func Blue(s string) string {
	return Colored(BlueColor, s)
}

func Purple(s string) string {
	return Colored(PurpleColor, s)
}

func Cyan(s string) string {
	return Colored(CyanColor, s)
}

func Gray(s string) string {
	return Colored(GrayColor, s)
}

func White(s string) string {
	return Colored(WhiteColor, s)
}

func Colored(color, s string) string {
	return color + s + Reset
}

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		RedColor = ""
		GreenColor = ""
		YellowColor = ""
		BlueColor = ""
		PurpleColor = ""
		CyanColor = ""
		GrayColor = ""
		WhiteColor = ""
	}
}
