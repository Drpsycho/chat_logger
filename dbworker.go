package main

import (
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"bytes"
	"fmt"
	"time"
	"github.com/Drpsycho/now"
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

func Convert2txt(){
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("Events")).Cursor()
		min := []byte(strconv.FormatInt(now.BeginningOfMonth().Unix(), 10))
		max := []byte(strconv.FormatInt(time.Now().Unix(), 10))
		// Iterate over the 90's.
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}

		return nil
	})
}