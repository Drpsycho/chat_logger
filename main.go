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
var initdb = flag.Bool("init", false, "It's for first run, initialization DB")

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
	if !*initdb {

	}

	msg := make(chan chanMsg, 100)
	go doSlack(*token, msg)
	go saveMsg("./db", msg)

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
