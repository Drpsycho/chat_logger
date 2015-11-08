package main

import (
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
)

var workDB db.DB

func initDB(pathToDB string) {
	fmt.Println("create database")
	DB, err := db.OpenDB(pathToDB)
	if err != nil {
		panic(err)
	}
	workDB = *DB
}

func saveMsg(pathToDB string, msg chan chanMsg) {
	initDB(pathToDB)

	var colsName []string

	fmt.Println("we got collection:")
	for _, name := range workDB.AllCols() {
		fmt.Println(name)
		colsName = append(colsName, name)
	}

	for true {
		t := <-msg

		if !CheckNameDB(colsName, t) {
			if err := workDB.Create(t.channelName); err != nil {
				fmt.Println(err)
			}
			colsName = append(colsName, t.channelName)
		}

		channelInDB := workDB.Use(t.channelName)

		if !IsMsgExist("timestamp", t.timestamp.String(), t.channelName) {
			fmt.Println("Message added")
			_, err := channelInDB.Insert(map[string]interface{}{
				"name":      t.author,
				"text":      t.text,
				"timestamp": t.timestamp.String()})
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Message already exist")
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
	fmt.Println("index is found ? ", indexFound)
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

func CheckNameDB(colsName []string, t chanMsg) bool {
	for k := range colsName {
		if colsName[k] == t.channelName {
			// fmt.Println("we already got a channel")
			return true
		}
	}
	return false
}
