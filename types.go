package main

import "time"

// MLHEvent represents an event in the Major League Hacking API
type MLHEvent struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		CreatedAt        time.Time `json:"createdAt"`
		UpdatedAt        time.Time `json:"updatedAt"`
		Status           string    `json:"status"`
		Digital          bool      `json:"digital"`
		RegistrationOpen bool      `json:"registrationOpen"`
		Title            string    `json:"title"`
		Images           struct {
			Icon   string `json:"icon"`
			Banner string `json:"banner"`
		} `json:"images"`
		StartDate string `json:"startDate"`
		StartTime string `json:"startTime"`
		EndDate   string `json:"endDate"`
		EndTime   string `json:"endTime"`
		Location  struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
		Host struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Website     string `json:"website"`
		} `json:"host"`
		Schedule        string    `json:"schedule"`
		ParticipantType []string  `json:"participantType"`
		TimeZone        string    `json:"timeZone"`
		StartDatetime   time.Time `json:"startDatetime"`
		EndDatetime     time.Time `json:"endDatetime"`
	} `json:"attributes"`
	Links struct {
		Register string `json:"register"`
		View     string `json:"view"`
	} `json:"links"`
}

// MLHGetEventsAPIResponse represents a response from the Major League Hacking API for getting events
type MLHGetEventsAPIResponse struct {
	Data []MLHEvent `json:"data"`
}
