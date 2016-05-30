package main

import (
	"bytes"
	"github.com/Drpsycho/now"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
	"time"
)

var db *bolt.DB
var err error

func InitDB() {
	db, err = bolt.Open(".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

func SaveInBucket(msg chanMsg) {
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(msg.channelName))
		if err != nil {
			log.Fatal(err)
			return nil
		}
		b.Put([]byte(strconv.FormatInt(msg.timestamp, 10)), []byte("["+msg.author+"]: "+msg.text))
		return nil
	})
}

func SaveMsg(msg chan chanMsg) {
	for {
		SaveInBucket(<-msg)
	}
}

func SaveChannels(channels map[string]string) {
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("channels"))
		if err != nil {
			log.Fatal(err)
			return nil
		}

		for channelname, channelid := range channels {
			b.Put([]byte(channelid), []byte(channelname))
		}
		return nil
	})
}

func GetChannelsName() []string {
	var channels []string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		b.ForEach(func(k, v []byte) error {
			channels = append(channels, string(v))
			return nil
		})
		return nil
	})
	return channels
}

func GetMsgByTime(channel string, newest string, latest string, msg_transfer chan string) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(channel)).Cursor()

		for k, v := c.Seek([]byte(latest)); k != nil && bytes.Compare(k, []byte(newest)) <= 0; k, v = c.Next() {
			unixIntValue, _ := strconv.ParseInt(string(k), 10, 64)
			date := time.Unix(unixIntValue, 0)
			msg_transfer <- date.Format("2006-01-02 15:04:05") + " " + string(v) + "\n"
		}

		return nil
	})
	msg_transfer <- "done"
}

func WriteMsgToDisk() {
	time.Sleep(1 * time.Minute)
	for {
		channels := GetChannelsName()
		path2folder := "txt/"
		for ch := range channels {
			_path := path2folder + channels[ch] + "_" + strconv.Itoa(now.BeginningOfYear().Year()) + ".txt"
			file, err := os.OpenFile(_path, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				panic(err)
			}
			latest := strconv.FormatInt(now.BeginningOfYear().Unix(), 10)
			newest := strconv.FormatInt(time.Now().Unix(), 10)
			msg_transfer := make(chan string, 100)
			go GetMsgByTime(channels[ch], newest, latest, msg_transfer)
			for {
				tmp := <-msg_transfer
				if tmp == "done" {
					break
				}

				if _, err = file.WriteString(tmp); err != nil {
					panic(err)
				}

			}
					if err := file.Close(); err != nil {
					panic(err)
				}

	}
		time.Sleep(3 * time.Hour)
	}
}
