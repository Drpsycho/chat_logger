package main

import (
	"flag"
	"fmt"
	"os"
	// "time"
)

const (
	VERSION = "0.1"
)

//----------------------------------------------------------------------------

var token = flag.String("token", "", "Token for slack")

type chanMsg struct {
	author      string
	text        string
	timestamp   int64
	channelName string
	channelId   string
}

//----------------------------------------------------------------------------

func usage() {
	fmt.Println("chat_logger v.", VERSION)
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *token == "" {
		usage()
	}

	msg := make(chan chanMsg, 100)
	InitDB()

	go SaveMsg(msg)
	go GetAllSlackMsg(*token, msg)

	var inputs string

	for {
		fmt.Scanln(&inputs)
		fmt.Println("For quit enter 'q'")
		if inputs == "q" {
			fmt.Println("done")
			CloseDB()
			os.Exit(0)
		}
	}
}
