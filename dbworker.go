package main

import (
	"github.com/boltdb/bolt"
	"log"
	"strconv"
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
