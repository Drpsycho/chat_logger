/*

rtop-bot - remote system monitoring bot

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
	fmt.Println("slack-bot v.", VERSION)
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

	token1 := string("xoxp-4588019148-4634100797-12065127504-531e321d4c")

	msg := make(chan chanMsg, 100)
	go doSlack(token1, msg)
	go saveMsg("./db",msg)

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
