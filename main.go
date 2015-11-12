package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	VERSION = "0.1"
)

//----------------------------------------------------------------------------

var token = flag.String("token", "", "Token for slack")

// var initdb = flag.Bool("init", false, "It's for first run, initialization DB")

type chanMsg struct {
	author      string
	text        string
	timestamp   time.Time
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
	// if !*initdb {

	// }
	go func() {
		var inputs string
		q := true

		for q {
			fmt.Scanln(&inputs)
			fmt.Println("For quit enter 'q'")
			fmt.Println("You enter ", inputs)
			if inputs == "q" {
				q = false
			}
		}
		fmt.Println("done")

		os.Exit(0)
	}()

	initDB("./db")

	msg := make(chan chanMsg, 100)

	go SaveMsg(msg)
	GetAllSlackMsg(*token, msg)
}
