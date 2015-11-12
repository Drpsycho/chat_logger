package main

import (
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
)

var workDB db.DB
var colsName []string

func initDB(pathToDB string) {
	fmt.Println("create database")
	DB, err := db.OpenDB(pathToDB)
	if err != nil {
		panic(err)
	}
	workDB = *DB

	for _, name := range workDB.AllCols() {
		fmt.Println(name)
		colsName = append(colsName, name)
	}
}

func SaveMsgSafe(msg chan chanMsg, quit chan bool) {
	for {
		select {
		case inmsg := <-msg:
			if !CheckNameDB(colsName, inmsg.channelName) {
				if err := workDB.Create(inmsg.channelName); err != nil {
					fmt.Println(err)
				}
				colsName = append(colsName, inmsg.channelName)
			}

			channelInDB := workDB.Use(inmsg.channelName)

			if !IsMsgExist("timestamp", inmsg.timestamp.String(), inmsg.channelName) {
				//				fmt.Println("Message added")
				_, err := channelInDB.Insert(map[string]interface{}{
					"name":      inmsg.author,
					"text":      inmsg.text,
					"timestamp": inmsg.timestamp.String()})
				if err != nil {
					panic(err)
				}
			} else {
				//				fmt.Println("Message already exist")
			}
		case <-quit:
			return
		}
	}
}

func SaveMsg(msg chan chanMsg) {
	for {
		inmsg := <-msg
		if !CheckNameDB(colsName, inmsg.channelName) {
			if err := workDB.Create(inmsg.channelName); err != nil {
				fmt.Println(err)
			}
			colsName = append(colsName, inmsg.channelName)
		}

		channelInDB := workDB.Use(inmsg.channelName)

		if !IsMsgExist("timestamp", inmsg.timestamp.String(), inmsg.channelName) {
			//			fmt.Println("Message added")
			_, err := channelInDB.Insert(map[string]interface{}{
				"name":      inmsg.author,
				"text":      inmsg.text,
				"timestamp": inmsg.timestamp.String()})
			if err != nil {
				panic(err)
			}
		} else {
			//			fmt.Println("Message already exist")
		}
	}
}

func IsMsgExist(key, value, collectionName string) bool {
	// return false
	query := map[string]interface{}{
		"eq": value,
		"in": []interface{}{key},
	}

	coll := workDB.Use(collectionName)

	indexFound := false
	for _, path := range coll.AllIndexes() {
		if path[0] == key {
			indexFound = true
			break
		}
	}
	//	fmt.Println("index is found ? ", indexFound)
	if !indexFound {
		if err := coll.Index([]string{key}); err != nil {
			panic(err)
		}
	}

	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(query, coll, &queryResult); err != nil {
		panic(err)
	} else {
		if len(queryResult) > 0 {
			return true
		}
	}
	return false
}

func CheckNameDB(colsName []string, t string) bool {
	for k := range colsName {
		if colsName[k] == t {
			// fmt.Println("we already got a channel")
			return true
		}
	}
	return false
}
