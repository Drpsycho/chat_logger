package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
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
		tmp := <-msg
		SaveInBucket(tmp)
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
			msg_transfer <- date.String() + " " + string(v)
			// fmt.Printf("%s: %s\n", k, v)
		}

		return nil
	})
	msg_transfer <- "done"
}
