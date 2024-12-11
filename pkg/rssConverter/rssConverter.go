package rssConverter

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/tunnelchaos/hopger/pkg/config"
	"github.com/tunnelchaos/hopger/pkg/gopherhelpers"
)

type RSSConverter struct{}

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
		text, err := gopherhelpers.ConvertHTMLToText(item.Description)
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
