package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	pageSize          = 200
	channelNamePrefix = "mlh-"
)

func main() {
	mattermostURL := os.Getenv("MM_URL")
	if mattermostURL == "" {
		fmt.Println("env var MM_URL must be set")
		return
	}
	teamID := os.Getenv("MM_TEAM_ID")
	if teamID == "" {
		fmt.Println("env var MM_TEAM_ID must be set")
		return
	}
	apiToken := os.Getenv("MM_API_TOKEN")
	if apiToken == "" {
		fmt.Println("env var MM_API_TOKEN must be set")
		return
	}
	mlhEventsURL := os.Getenv("MLH_EVENTS_URL")
	if mlhEventsURL == "" {
		fmt.Println("env var MLH_EVENTS_URL must be set")
		return
	}

	fmt.Printf("running channel generation against Mattermost at %s\n", mattermostURL)

	fmt.Println("getting events from MLH API")
	mlhResponse, err := http.Get(mlhEventsURL)
	if err != nil {
		fmt.Printf("exiting - failed getting MLH events, error=%v\n", err.Error())
		return
	}
	defer mlhResponse.Body.Close()
	if mlhResponse.StatusCode != http.StatusOK {
		fmt.Printf("exiting - got a %v status code when expected 200 from MLH API\n", mlhResponse.StatusCode)
		return
	}

	var getEventsAPIResponse *MLHGetEventsAPIResponse
	json.NewDecoder(mlhResponse.Body).Decode(&getEventsAPIResponse)

	if getEventsAPIResponse.Data == nil || len(getEventsAPIResponse.Data) == 0 {
		fmt.Println("exiting - no events returned")
		return
	}

	events := getEventsAPIResponse.Data
	fmt.Printf("found %v events\n", len(events))

	client := model.NewAPIv4Client(mattermostURL)
	client.SetToken(apiToken)

	channelsMap := map[string]*model.Channel{}
	numChannelsInLastPage := pageSize
	for page := 0; numChannelsInLastPage >= pageSize; page++ {
		fmt.Printf("getting page %v of MM channels\n", page)

		channelsPage, resp := client.GetPublicChannelsForTeam(teamID, page, pageSize, "")
		if resp.Error != nil {
			fmt.Printf("exiting - got a %v status code when expected 200 from MM channels GET, page=%v, error=%v\n", mlhResponse.StatusCode, page, resp.Error.Error())
			return
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("exiting - got a %v status code when expected 200 from MM channels GET, page=%v\n", mlhResponse.StatusCode, page)
			return
		}

		for _, channel := range channelsPage {
			if len(channel.Name) >= 4 && channel.Name[:4] == channelNamePrefix {
				eventID := strings.Replace(channel.Name, channelNamePrefix, "", 1)
				channelsMap[eventID] = channel
			}
		}

		numChannelsInLastPage = len(channelsPage)
	}

	fmt.Printf("found %v existing event channels\n", len(channelsMap))

	createdCount := 0
	for _, event := range events {
		_, ok := channelsMap[event.ID]
		if ok {
			fmt.Printf("skipping event with ID %v as channel already exists\n", event.ID)
			continue
		}

		attributes := event.Attributes

		fmt.Printf("creating channel for new event with ID %s and title %s \n", event.ID, attributes.Title)

		channel := &model.Channel{
			Name:        fmt.Sprintf("%s%s", channelNamePrefix, event.ID),
			DisplayName: cleanChannelDisplayName(attributes.Title),
			TeamId:      teamID,
			Type:        model.CHANNEL_OPEN,
			Header:      fmt.Sprintf("Start Date: %v at %v (%v) | Location: %v, %v | Host: [%v](%v) | ID: %v", attributes.StartDate, attributes.StartTime, attributes.TimeZone, attributes.Location.City, attributes.Location.Country, attributes.Host.Name, attributes.Host.Website, event.ID),
		}

		createdChannel, resp := client.CreateChannel(channel)
		if resp.Error != nil {
			fmt.Printf("failed to create channel for event with ID %v error=%v\n", event.ID, resp.Error.Error())
			continue
		}

		message := fmt.Sprintf(`## Welcome to [%v](%v)!

#### Event Schedule

%v

#### More About the Host

%v
		`, attributes.Title, event.Links.View, attributes.Schedule, attributes.Host.Description)

		post := &model.Post{
			ChannelId: createdChannel.Id,
			IsPinned:  true,
			Message:   message,
		}

		_, resp = client.CreatePost(post)
		if resp.Error != nil {
			fmt.Printf("failed to create post for event with ID %v error=%v\n", event.ID, resp.Error.Error())
			continue
		}

		createdCount++
	}

	fmt.Printf("created %v new event channels\n", createdCount)
	fmt.Println("success")
}

func cleanChannelDisplayName(s string) string {
	if len(s) > model.CHANNEL_NAME_MAX_LENGTH {
		return s[:model.CHANNEL_NAME_MAX_LENGTH]
	}
	return s
}
