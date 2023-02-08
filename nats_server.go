package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

func publishMessage() {
	serv, err := stan.Connect("test-cluster", "1")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connection OK")
	}
	defer serv.Close()

	var data []byte
	data, err = json.Marshal(getRandomMessage())
	if err != nil {
		log.Fatal(err)
	}
	err = serv.Publish("wb", data)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Published message to NATS Server")
}

func messageCatcher(msg *stan.Msg) {

	var tmp Order
	json.Unmarshal(msg.Data, &tmp)

	if tmp.Order_uid == "" {
		log.Println("Empty order_uid in message")
		return
	}

	cache = append(cache, msg.Data)
	row := db.QueryRow(context.Background(),
		"INSERT INTO cache(message) VALUES('"+string(msg.Data)+"') RETURNING id;")
	var psql_id int
	err := row.Scan(&psql_id)
	if err != nil {
		log.Fatal("messageCatcher: ", err)
	}
	insertIntoDatabase(&tmp)
	log.Printf("Message caught, cache id: %v, psql_cache id: %v", len(cache), psql_id)
}
