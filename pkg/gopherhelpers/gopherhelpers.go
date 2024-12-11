package gopherhelpers

import (
	"strings"

	"git.mills.io/prologic/go-gopher"
)

const MaxLine = 70

func FillLineWithChar(line string, length int, char string) string {
	return line + strings.Repeat(char, length-len(line))
}

func CreateMaxLine(line string) string {
	return strings.Repeat(line, MaxLine)
}

func FillLine(line string, length int) string {
	return FillLineWithChar(line, length, " ")
}

func FormatInfo(indent int, header string, content string) string {
	content = strings.ReplaceAll(content, "\n", "")
	indentstring := strings.Repeat(" ", indent)
	header = header + strings.Repeat(" ", indent-len(header))
	words := strings.Fields(content)
	section := ""
	currentline := header
	for _, word := range words {
		if len(currentline)+len(word)+1 > MaxLine {
			section += currentline + "\n"
			currentline = indentstring
		}
		currentline += word + " "
	}
	section += currentline + "\n"
	return section
}

func FormatForGopherMap(indent int, header string, content string) string {
	content = strings.ReplaceAll(content, "\n", "")
	header = header + strings.Repeat(" ", indent-len(header))
	section := header + content
	return section
}

func CreateEventHeader(title string) string {
	result := title + "\n"
	result += CreateMaxLine("=") + "\n"
	result += "Time          | Event\n"
	result += CreateMaxLine("-") + "\n"
	return result
}

func CreateGopherEntry(selectortype gopher.ItemType, Name string, selector string, host string, port int) string {
	item := gopher.Item{
		Type:        selectortype,
		Description: Name,
		Selector:    selector,
		Host:        host,
		Port:        port,
	}
	result, _ := item.MarshalText()
	return string(result)
}

func CreateGopherURL(Name string, URL string, Server string, Port int) string {
	return CreateGopherEntry(gopher.HTML, Name, "URL:"+URL, Server, Port)
}

func CreateGopherInfo(Heading string) string {
	return CreateGopherEntry(gopher.INFO, Heading, "fake", "(NULL)", 0)
}
