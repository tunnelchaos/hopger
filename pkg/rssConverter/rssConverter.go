package rssConverter

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/tunnelchaos/hopger/pkg/config"
	"golang.org/x/net/html"
)

type RSSConverter struct{}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		// Get the text from text nodes
		return n.Data
	}
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		// Skip <script> and <style> content
		return ""
	}

	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(extractText(c))
	}

	// Add spacing for certain elements to preserve readability
	if n.Type == html.ElementNode {
		switch n.Data {
		case "p", "br":
			buf.WriteString("\n")
		case "h1", "h2", "h3", "h4", "h5", "h6":
			headerText := strings.TrimSpace(buf.String())
			buf.Reset()
			buf.WriteString(headerText)
			buf.WriteString("\n" + strings.Repeat("=", len(headerText)) + "\n")
		case "li":
			buf.WriteString("- ")
		}
	}

	return buf.String()
}

// ConvertHTMLToText converts HTML content to plain text
func convertHTMLToText(htmlContent string) (string, error) {
	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// Extract text
	text := extractText(doc)

	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n\n\n", "\n\n") // Collapse excessive newlines
	return text, nil
}

func (r *RSSConverter) Convert(eventname string, info config.Info, server config.Server) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(info.URL)
	if err != nil {
		return errors.New("Failed to parse RSS feed: " + err.Error())
	}
	if len(feed.Items) == 0 {
		return nil
	}
	feedPath := path.Join(server.GopherDir, eventname, info.Name)
	os.RemoveAll(feedPath)
	os.MkdirAll(feedPath, 0755)
	gophermap := "Here is the RSS feed for " + feed.Title + "\n\n"
	for _, item := range feed.Items {
		date := item.PublishedParsed.Format("2006-01-02")
		itemPath := path.Join(feedPath, date+".txt")
		//Delete File if it exists
		os.Remove(itemPath)
		f, err := os.Create(itemPath)
		if err != nil {
			return errors.New("Failed to create file: " + err.Error())
		}
		defer f.Close()
		title := item.Title + "\n"
		title += strings.Repeat("=", len(item.Title)) + "\n\n"
		text, err := convertHTMLToText(item.Description)
		if err != nil {
			return errors.New("Failed to convert HTML to text: " + err.Error())
		}
		_, err = f.WriteString(title + text)
		if err != nil {
			return errors.New("Failed to write to file: " + err.Error())
		}
		gophermap += "0" + item.Title + "\t" + date + ".txt\n"
	}
	gophermapPath := path.Join(feedPath, "gophermap")
	f, err := os.Create(gophermapPath)
	if err != nil {
		return errors.New("Failed to create gophermap file: " + err.Error())
	}
	defer f.Close()
	_, err = f.WriteString(gophermap)
	if err != nil {
		return errors.New("Failed to write to gophermap file: " + err.Error())
	}
	return nil
}
