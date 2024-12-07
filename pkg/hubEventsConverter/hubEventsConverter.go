package hubEventsConverter

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/tunnelchaos/hopger/pkg/config"
	"github.com/tunnelchaos/hopger/pkg/helpers"
)

type EventKind string

const (
	KindOffical              EventKind = "official"
	KindSelfOrganizedSession EventKind = "sos"
	KindAssembly             EventKind = "assembly"
)

type HubEvent struct {
	ID               string    `json:"id"`
	Kind             EventKind `json:"kind"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	URL              string    `json:"url"`
	Track            string    `json:"track"`
	Assembly         string    `json:"assembly"`
	Room             string    `json:"room"`
	Language         string    `json:"language"`
	Description      string    `json:"description"`
	ScheduleStart    time.Time `json:"schedule_start"`
	ScheduleDuration string    `json:"schedule_duration"`
	ScheduleEnd      time.Time `json:"schedule_end"`
}

type HubEventsConverter struct {
}

func sortEvents(events []HubEvent) {

	sort.Slice(events, func(i, j int) bool {
		if events[i].ScheduleStart.Equal(events[j].ScheduleStart) {
			return events[i].Assembly < events[j].Assembly
		}
		return events[i].ScheduleStart.Before(events[j].ScheduleStart)
	})

}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (e *HubEvent) eventToGopher(loc *time.Location, addDate bool) string {

	timestring := e.ScheduleStart.In(loc).Format("15:04") + " - " + e.ScheduleEnd.In(loc).Format("15:04") + "   "
	indent := len(timestring)
	eventstring := helpers.FormatInfo(indent, timestring, e.Name)
	if addDate {
		eventstring = eventstring + helpers.FormatInfo(indent, "Date:", e.ScheduleStart.In(loc).Format("2006-01-02"))
	}
	eventstring += helpers.FormatInfo(indent, "Description:", e.Description)
	eventstring += helpers.FormatInfo(indent, "Language:", e.Language)
	eventstring += helpers.FormatInfo(indent, "Track:", e.Track)
	eventstring += helpers.CreateMaxLine("-") + "\n"
	return eventstring
}

func writeEventFiles(path, content string) error {
	file, err := os.Create(path + ".txt")
	if err != nil {
		return errors.New("Failed to create track file: " + err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return errors.New("Failed to write to track file: " + err.Error())
	}
	return nil
}

func (c *HubEventsConverter) Convert(eventname string, info config.Info, server config.Server) error {
	httpClient := helpers.CreateHttpClient()
	response, err := httpClient.Get(info.URL)
	if err != nil {
		return errors.New("Failed to fetch hub schedule from URL: " + info.URL + ":" + err.Error())
	}
	defer response.Body.Close()
	var hubEvents []HubEvent
	err = json.NewDecoder(response.Body).Decode(&hubEvents)
	if err != nil {
		return errors.New("Failed to parse hub schedule: " + info.Name + ":" + err.Error())
	}
	//Sort events by start time
	sortEvents(hubEvents)
	// Convert events to gopher
	basepath := path.Join(server.GopherDir, eventname, info.Name)
	os.RemoveAll(basepath)
	bydatepath := path.Join(basepath, "By Date")
	var oldDay time.Time
	daycount := -1
	daypath := ""
	daystring := ""
	//tracks := make(map[string]string)
	assemblies := make(map[string]string)
	os.MkdirAll(bydatepath, 0755)
	for _, event := range hubEvents {
		if event.Kind == KindOffical {
			continue
		}
		if !DateEqual(event.ScheduleStart, oldDay) {
			err := writeEventFiles(daypath, daystring)
			if err != nil {
				return errors.New("Failed to write event day file: " + err.Error())
			}
			daycount++
			dayname := "Day " + strconv.Itoa(daycount) + ": " + event.ScheduleStart.Format("2006-01-02")
			daypath = path.Join(bydatepath, dayname)
			oldDay = event.ScheduleStart
			daystring = helpers.CreateEventHeader(dayname)
		}
		daystring += event.eventToGopher(time.Local, false)
		/*if event.Track != "" {
			fmt.Println("Event Track: ", event.Track)
			if tracks[event.Track] == "" {
				tracks[event.Track] = helpers.CreateEventHeader(event.Track)
			}
			tracks[event.Track] += event.eventToGopher(time.Local, true)
		}*/
		if assemblies[event.Assembly] == "" {
			assemblies[event.Assembly] = helpers.CreateEventHeader(event.Assembly)
		}
		assemblies[event.Assembly] += event.eventToGopher(time.Local, true)
	}
	err = writeEventFiles(daypath, daystring)
	if err != nil {
		return errors.New("Failed to write event day file: " + err.Error())
	}
	/*bytrackpath := path.Join(basepath, "By Track")
	os.MkdirAll(bytrackpath, 0755)
	for k, v := range tracks {
		err := writeEventFiles(path.Join(bytrackpath, k), v)
		if err != nil {
			return errors.New("Failed to write track file: " + err.Error())
		}
	}*/
	byassemblypath := path.Join(basepath, "By Assembly")
	os.MkdirAll(byassemblypath, 0755)
	for k, v := range assemblies {
		err := writeEventFiles(path.Join(byassemblypath, k), v)
		if err != nil {
			return errors.New("Failed to write assembly file: " + err.Error())
		}
	}

	return nil
}
