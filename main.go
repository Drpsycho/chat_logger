package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	VERSION = "0.2"
)

var token = flag.String("token", "", "Token for slack")

type chanMsg struct {
	author      string
	text        string
	timestamp   int64
	channelName string
	channelId   string
}

func main() {
	fmt.Println("chat_logger v.", VERSION)
	flag.Parse()
	log.SetFlags(0)
	if *token == "" {
		usage()
	}

	msg := make(chan chanMsg, 100)

	InitDB()

	go SaveMsg(msg)
	go GetAllSlackMsg(*token, msg)
	go WriteMsgToDisk()

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

func usage() {
	flag.Usage()
	os.Exit(1)
}
