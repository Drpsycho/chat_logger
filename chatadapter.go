package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strconv"
	"strings"
	"time"
)

type Message slack.Message

var api *slack.Client

var mapNames map[string]string
var mapChannels map[string]string

func FillChannelList() {
	mapChannels = make(map[string]string)
	channels, _ := api.GetChannels(false)
	for k := range channels {
		mapChannels[channels[k].Name] = channels[k].ID
	}
}

func FillUserName() {
	mapNames = make(map[string]string)
	users, _ := api.GetUsers()
	for k := range users {
		mapNames[users[k].ID] = users[k].Name
	}
}

func GetAllSlackMsg(token string, mesg chan chanMsg) {
	for {

		{
			api = slack.New(token)

			FillChannelList()
			FillUserName()

			SaveChannels(mapChannels)

			params := slack.HistoryParameters{Count: 1000}
			fmt.Println("Get all message and save it in DB")
			for channelname, channelid := range mapChannels {
				history, _ := api.GetChannelHistory(channelid, params)

				for messages := range history.Messages {
					msg := history.Messages[messages]
					str_timestamp := strings.Split(msg.Timestamp, ".")
					unixIntValue, _ := strconv.ParseInt(str_timestamp[0], 10, 64)
					timeStamp := unixIntValue

					mesg <- chanMsg{author: mapNames[msg.User],
						text:        msg.Text,
						timestamp:   timeStamp,
						channelId:   channelid,
						channelName: channelname}
				}
			}
		}

		fmt.Println("All message have been saved")
		fmt.Println("sleep 24 hour")
		time.Sleep(24 * time.Hour)
	}

}
