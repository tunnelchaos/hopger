package assemblyConverter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/tunnelchaos/go-packages/config"
	"github.com/tunnelchaos/go-packages/gopherhelpers"
	"github.com/tunnelchaos/go-packages/helpers"
)

type AssemblyConverter struct{}

type assembly struct {
	Slug             string `json:"slug"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	Parent           any    `json:"parent"`
	AssemblyLocation string `json:"assembly_location"`
	AssemblyLink     string `json:"assembly_link"`
	IsOfficial       bool   `json:"is_official"`
	EventsURL        string `json:"events_url"`
	RoomsURL         string `json:"rooms_url"`
	BadgesURL        string `json:"badges_url"`
}

func (a *assembly) toGopher(Server string, Port int) string {
	maxline := len("Parent: ")
	parent := ""
	result := gopherhelpers.CreateGopherInfo(gopherhelpers.FormatForGopherMap(maxline, "Name:", a.Name))
	if a.Parent != nil {
		//Interpret parent as string
		parent = a.Parent.(string)
	}
	result += gopherhelpers.CreateGopherInfo(gopherhelpers.FormatForGopherMap(maxline, "Parent:", parent))
	result += gopherhelpers.CreateGopherURL(gopherhelpers.FormatForGopherMap(maxline, "Link:", ""), a.AssemblyLink, Server, Port)
	result += gopherhelpers.CreateGopherInfo(gopherhelpers.FillLineWithChar("", gopherhelpers.MaxLine-1, "-"))
	return result
}

func (ac *AssemblyConverter) Convert(eventname string, info config.Info, server config.Server) error {
	httpClient := helpers.CreateHttpClient()
	// Fetch assemblies
	response, err := httpClient.Get(info.URL)
	if err != nil {
		return errors.New("Failed to fetch assembly data from URL: " + info.URL + ":" + err.Error())
	}
	defer response.Body.Close()
	var assemblies []assembly
	err = json.NewDecoder(response.Body).Decode(&assemblies)
	if err != nil {
		return errors.New("Failed to parse assembly data: " + info.Name + ":" + err.Error())
	}
	// Convert assemblies to gopher
	result := gopherhelpers.CreateGopherInfo(fmt.Sprintf("All Assemblies from %s", eventname))
	result += gopherhelpers.CreateGopherInfo(gopherhelpers.FillLineWithChar("", gopherhelpers.MaxLine-1, "="))
	for _, assembly := range assemblies {
		result += assembly.toGopher(server.Hostname, server.GopherPort)
	}
	// Save gopher data
	basepath := path.Join(server.GopherDir, eventname, info.Name)
	os.RemoveAll(basepath)
	os.MkdirAll(basepath, 0755)
	gophermapPath := path.Join(basepath, "gophermap")
	f, err := os.Create(gophermapPath)
	if err != nil {
		return errors.New("Failed to create gophermap file: " + err.Error())
	}
	defer f.Close()
	_, err = f.WriteString(result)
	if err != nil {
		return errors.New("Failed to write to gophermap file: " + err.Error())
	}

	return nil
}
