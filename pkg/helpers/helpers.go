package helpers

import (
	"net/http"
	"strings"
	"time"
)

const MaxLine = 70

func FillLineWithChar(line string, length int, char string) string {
	return strings.Repeat(char, length-len(line))
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

func CreateGopherEntry(selectortype string, Name string, selector string, host string, port string) string {
	return selectortype + Name + "\t" + selector + "\t" + host + "\t" + port + "\n"
}

func CreateGopherURL(Name string, URL string, Server string, Port string) string {
	return CreateGopherEntry("h", Name, "URL:"+URL, Server, Port)
}

func CreateGopherInfo(Heading string) string {
	return CreateGopherEntry("i", Heading, "fake", "(NULL)", "0")
}

func CreateHttpClient() *http.Client {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return httpClient
}
