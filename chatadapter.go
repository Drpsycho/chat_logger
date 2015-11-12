package main

import (
	"github.com/Drpsycho/now"
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
		//		fmt.Println(channels[k].Name, channels[k].ID)
		mapChannels[channels[k].Name] = channels[k].ID
	}
}

func FillUserName() {
	mapNames = make(map[string]string)
	users, _ := api.GetUsers()
	for k := range users {
		//		fmt.Println(users[k].Name, users[k].ID)
		mapNames[users[k].ID] = users[k].Name
	}
}

func ConvertToString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func GetAllSlackMsg(token string, mesg chan chanMsg, quit chan bool) {
	api = slack.New(token)

	FillChannelList()
	FillUserName()

	params := slack.HistoryParameters{Count: 1000}

	for channelname, channelid := range mapChannels {
		history, _ := api.GetChannelHistory(channelid, params)

		for messages := range history.Messages {
			msg := history.Messages[messages]
			str_timestamp := strings.Split(msg.Timestamp, ".")
			unixIntValue, _ := strconv.ParseInt(str_timestamp[0], 10, 64)
			timeStamp := time.Unix(unixIntValue, 0)
			mesg <- chanMsg{author: mapNames[msg.User],
				text:        msg.Text,
				timestamp:   timeStamp,
				channelId:   channelid,
				channelName: channelname}
		}
	}
	quit <- true
}

func GetlackMsgEveryDay(token string, mesg chan chanMsg) {
	api = slack.New(token)

	FillChannelList()
	FillUserName()
	for {
		beginDayToday := now.BeginningOfDay()
		beginDayBefore := now.OneDayBefore(beginDayToday)

		params := slack.HistoryParameters{
			Oldest: ConvertToString(beginDayBefore),
			Latest: ConvertToString(beginDayToday)}

		for channelname, channelid := range mapChannels {
			history, _ := api.GetChannelHistory(channelid, params)

			for messages := range history.Messages {
				msg := history.Messages[messages]
				str_timestamp := strings.Split(msg.Timestamp, ".")
				unixIntValue, _ := strconv.ParseInt(str_timestamp[0], 10, 64)
				timeStamp := time.Unix(unixIntValue, 0)
				mesg <- chanMsg{author: mapNames[msg.User],
					text:        msg.Text,
					timestamp:   timeStamp,
					channelId:   channelid,
					channelName: channelname}
			}
		}
		time.Sleep(24 * time.Hour)
	}
}
