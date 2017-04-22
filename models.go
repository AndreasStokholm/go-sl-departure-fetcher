package main

// DepartureSite - Config structure for departure sites
type DepartureSite struct {
	Directions    []int
	SiteID        int    `mapstructure:"site_id"`
	TransportMode string `mapstructure:"transport_mode"`
}

// Departure - Structure of departures from the SL API
type Departure struct {
	GroupOfLine          string
	DisplayTime          string
	TransportMode        string
	LineNumber           string
	Destination          string
	JourneyDirection     int
	StopAreaName         string
	StopAreaNumber       int
	StopPointNumber      int
	StopPointDesignation string
	TimeTabledDateTime   string
	ExpectedDateTime     string
	JourneyNumber        int
	Deviations           []string
}

// MQTTDeparture - The structure of the object to be published to the MQTT broker
type MQTTDeparture struct {
	TransportMode    string `json:"transport_mode"`
	TimeOfDeparture  string `json:"time_of_departure"`
	JourneyDirection int    `json:"journey_direction"`
}
