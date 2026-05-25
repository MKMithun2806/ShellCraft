package main

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
)

func col(text, color string) string {
	return color + text + Reset
}

func errorf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Red)
}

func successf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Green)
}

func infof(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Yellow)
}

func promptf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Blue)
}

func headerf(format string, a ...interface{}) string {
	return Bold + Yellow + fmt.Sprintf(format, a...) + Reset
}

func payloadf(format string, a ...interface{}) string {
	return col(fmt.Sprintf(format, a...), Cyan)
}
