package helpers

import "strings"

const maxLine = 70

func CreateMaxLine(line string) string {
	return strings.Repeat(line, maxLine)
}
