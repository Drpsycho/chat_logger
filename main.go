package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
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

var addr = flag.String("addr", "localhost:8080", "http service address")

var homeTempl = template.Must(template.ParseFiles("/home/drpsycho/Yandex.Disk/site/sandbox.html"))
var upgrader = websocket.Upgrader{} // use default options

func home(w http.ResponseWriter, r *http.Request) {
	homeTempl.Execute(w, "ws://"+r.Host+"/reqdb")
}

type inputmsg struct {
	To     string
	Action string
	Data   string
}

type outputmsg struct {
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
}

func reqdb(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		dat := inputmsg{}

		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}

		if dat.To == "server" {
			switch dat.Action {
			case "init":
				channels := GetChannelsName()
				m := outputmsg{
					Channels: channels,
				}
				log.Println(m)
				res, _ := json.Marshal(m)
				c.WriteMessage(1, res)
				break
			case "getmsg":
				str_timestamp := strings.Split(dat.Data, ".")
				log.Println(str_timestamp)

				msg_transfer := make(chan string, 100)

				go GetMsgByTime(str_timestamp[0], str_timestamp[1], str_timestamp[2], msg_transfer)
				for {
					tmp := <-msg_transfer
					if tmp == "done" {
						break
					}
					m := outputmsg{
						Data: tmp,
					}
					log.Println(m)
					res, _ := json.Marshal(m)
					c.WriteMessage(1, res)
				}
				break
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	if *token == "" {
		usage()
	}

	msg := make(chan chanMsg, 100)
	InitDB()

	go SaveMsg(msg)
	go GetAllSlackMsg(*token, msg)

	http.HandleFunc("/reqdb", reqdb)
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))

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
