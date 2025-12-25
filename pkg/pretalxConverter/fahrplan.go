package pretalxConverter

// Root struct for the entire JSON
type Fahrplan struct {
	Schema    string    `json:"$schema"`
	Generator Generator `json:"generator"`
	Schedule  EventData `json:"schedule"`
}

type Generator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type EventData struct {
	URL        string     `json:"url"`
	Version    string     `json:"version"`
	BaseURL    string     `json:"base_url"`
	Conference Conference `json:"conference"`
}

type Conference struct {
	Acronym          string  `json:"acronym"`
	Title            string  `json:"title"`
	Start            string  `json:"start"`
	End              string  `json:"end"`
	DaysCount        int     `json:"daysCount"`
	TimeslotDuration string  `json:"timeslot_duration"`
	TimeZoneName     string  `json:"time_zone_name"`
	Colors           Colors  `json:"colors"`
	Rooms            []Room  `json:"rooms"`
	Tracks           []Track `json:"tracks"`
	Days             []Day   `json:"days"`
}

type Colors struct {
	Primary string `json:"primary"`
}

type Room struct {
	Name        string  `json:"name"`
	GUID        string  `json:"guid"`
	Description *string `json:"description"`
	Capacity    int     `json:"capacity"`
}

type Track struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Day struct {
	Index    int                `json:"index"`
	Date     string             `json:"date"`
	DayStart string             `json:"day_start"`
	DayEnd   string             `json:"day_end"`
	Rooms    map[string][]Event `json:"rooms"` // Dynamic keys for room names
}

type Event struct {
	URL              string        `json:"url"`
	ID               any           `json:"id"`
	GUID             string        `json:"guid"`
	Date             string        `json:"date"`
	Start            string        `json:"start"`
	Logo             *string       `json:"logo"`
	Duration         string        `json:"duration"`
	Room             string        `json:"room"`
	Slug             string        `json:"slug"`
	Title            string        `json:"title"`
	Subtitle         string        `json:"subtitle"`
	Track            string        `json:"track"`
	Type             string        `json:"type"`
	Language         string        `json:"language"`
	Abstract         string        `json:"abstract"`
	Description      string        `json:"description"`
	RecordingLicense *string       `json:"recording_license"`
	DoNotRecord      bool          `json:"do_not_record"`
	Persons          []Person      `json:"persons"`
	Links            []interface{} `json:"links"`
	Attachments      []interface{} `json:"attachments"`
	Answers          []interface{} `json:"answers"`
	saal             string
}

type Person struct {
	GUID       string        `json:"guid"`
	ID         int           `json:"id"`
	Code       string        `json:"code"`
	PublicName string        `json:"public_name"`
	Avatar     *string       `json:"avatar"`
	Biography  *string       `json:"biography"`
	Answers    []interface{} `json:"answers"`
}
