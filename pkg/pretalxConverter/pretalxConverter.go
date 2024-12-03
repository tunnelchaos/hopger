package pretalxConverter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tunnelchaos/hopger/pkg/config"
	"github.com/tunnelchaos/hopger/pkg/helpers"
)

func parseDuration(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	// Parse hours and minutes
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours in duration: %s", input)
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes in duration: %s", input)
	}

	// Convert to time.Duration
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute, nil
}

func formatSection(indent int, header string, content string) string {
	content = strings.ReplaceAll(content, "\n", "")
	indentstring := strings.Repeat(" ", indent)
	header = header + strings.Repeat(" ", indent-len(header))
	words := strings.Fields(content)
	section := ""
	currentline := header
	for _, word := range words {
		if len(currentline)+len(word)+1 > 70 {
			section += currentline + "\n"
			currentline = indentstring
		}
		currentline += word + " "
	}
	section += currentline + "\n"
	return section
}

func sortEvents(events []Event) error {

	sort.Slice(events, func(i, j int) bool {
		// Combine date and time for both events
		timeI, errI := time.Parse(time.RFC3339, events[i].Date)
		timeJ, errJ := time.Parse(time.RFC3339, events[j].Date)

		// Handle parsing errors
		if errI != nil || errJ != nil {
			fmt.Println("Error parsing dates:", errI, errJ)
			return false
		}

		// Compare the parsed times
		if timeI.Equal(timeJ) {
			return events[i].saal < events[j].saal
		}
		return timeI.Before(timeJ)
	})

	return nil
}

func eventToGopher(event Event, loc *time.Location, addDate bool, addSaal bool) string {
	//Parse start time and duration to get end time
	starttime, err := time.Parse(time.RFC3339, event.Date)
	if err != nil {
		return "Failed to parse event start time: " + err.Error() + "\n"
	}
	duration, err := parseDuration(event.Duration)
	if err != nil {
		return "Failed to parse event duration: " + err.Error() + "\n"
	}
	endtime := starttime.Add(duration)

	timestring := starttime.In(loc).Format("15:04") + " - " + endtime.In(loc).Format("15:04") + "   "
	indent := len(timestring)
	eventstring := formatSection(indent, timestring, event.Title)
	if addDate {
		d, err := time.Parse(time.RFC3339, event.Date)
		if err != nil {
			return "Failed to parse event date: " + err.Error() + "\n"
		}
		eventstring = eventstring + formatSection(indent, "Date:", d.In(loc).Format("2006-01-02"))
	}
	if addSaal {
		eventstring = eventstring + formatSection(indent, "Room:", event.saal)
	}
	eventstring += formatSection(indent, "Description:", event.Description)
	speakerHeaer := "Speaker"
	if len(event.Persons) > 1 {
		speakerHeaer += "s"
	}
	speakerHeaer += ":"
	speakerstring := " "
	for i, speaker := range event.Persons {
		speakerstring += speaker.PublicName
		if i != len(event.Persons)-1 {
			speakerstring += ", "
		}
	}
	eventstring += formatSection(indent, speakerHeaer, speakerstring)
	eventstring += formatSection(indent, "Language:", event.Language)
	eventstring += formatSection(indent, "Track:", event.Track)
	eventstring += helpers.CreateMaxLine("-") + "\n"
	return eventstring
}

func Convert(eventname string, info config.Info, server config.Server) error {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get(info.URL)
	if err != nil {
		return errors.New("Failed to fetch pretalx schedule from URL: " + info.URL + ":" + err.Error())
	}
	defer response.Body.Close()
	var fahrplan Fahrplan
	err = json.NewDecoder(response.Body).Decode(&fahrplan)
	if err != nil {
		return errors.New("Failed to parse pretalx schedule: " + info.Name + ":" + err.Error())
	}
	loc, err := time.LoadLocation(fahrplan.Schedule.Conference.TimeZoneName)
	if err != nil {
		return errors.New("Failed to load timezone: " + fahrplan.Schedule.Conference.TimeZoneName + ":" + err.Error())
	}
	tracks := make(map[string][]Event)
	for _, track := range fahrplan.Schedule.Conference.Tracks {
		tracks[track.Name] = []Event{}
	}
	basepath := path.Join(server.GopherDir, eventname, info.Name)
	os.RemoveAll(basepath)
	bydatepath := path.Join(basepath, "By Date")
	for i, day := range fahrplan.Schedule.Conference.Days {
		dayname := "Day " + strconv.Itoa(i) + ": " + day.Date
		daypath := path.Join(bydatepath, dayname)
		os.MkdirAll(daypath, 0755)
		for name, events := range day.Rooms {
			roomstring := "Room: " + name + "\n"
			roomstring += helpers.CreateMaxLine("=") + "\n"
			roomstring += "Time          | Event\n"
			roomstring += helpers.CreateMaxLine("-") + "\n"
			for _, event := range events {
				e := event
				e.saal = name
				tracks[event.Track] = append(tracks[event.Track], e)
				roomstring += eventToGopher(event, loc, false, false)
			}
			roompath := path.Join(daypath, name+".txt")
			roomfile, err := os.Create(roompath)
			if err != nil {
				return errors.New("Failed to create room file: " + err.Error())
			}
			defer roomfile.Close()
			_, err = roomfile.WriteString(roomstring)
			if err != nil {
				return errors.New("Failed to write to room file: " + err.Error())
			}
		}
	}
	bytrackpath := path.Join(basepath, "By Track")
	os.MkdirAll(bytrackpath, 0755)
	for k, v := range tracks {
		trackpath := path.Join(bytrackpath, k)
		trackstring := "Track: " + k + "\n"
		trackstring += helpers.CreateMaxLine("=") + "\n"
		trackstring += "Time          | Event\n"
		trackstring += helpers.CreateMaxLine("-") + "\n"
		sortEvents(v)
		for _, event := range v {
			trackstring += eventToGopher(event, loc, true, true)
		}
		trackfile, err := os.Create(trackpath + ".txt")
		if err != nil {
			return errors.New("Failed to create track file: " + err.Error())
		}
		defer trackfile.Close()
		_, err = trackfile.WriteString(trackstring)
		if err != nil {
			return errors.New("Failed to write to track file: " + err.Error())
		}
	}
	return nil
}
