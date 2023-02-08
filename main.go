package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/stan.go"
)

func checkArg(arg string) bool {
	for _, v := range os.Args[1:] {
		if v == arg {
			return true
		}
	}
	return false
}

func checkArgComplex(arg string) (int, error) {
	var t = false
	for _, v := range os.Args[1:] {
		if t {
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			return i, nil
		}
		if v == arg {
			t = true
		}
	}

	if t {
		return 0, errors.New("no id provided")
	}
	return 0, nil
}

// restored from psql cache
const CACHE_RESTORE_LIMIT = 5

var db *pgx.Conn
var cache [][]byte

func main() {

	db = connectToDatabase()
	defer db.Close(context.Background())

	restoreCache()

	// db args
	if checkArg("-d") {
		populateDatabase(d)
		log.Print("Truncated Database")
		return
	} else if checkArg("-c") {
		populateDatabase(c)
		log.Print("Created Database")
		return
	} else if checkArg("-r") {
		populateDatabase(r)
		log.Print("Refreshed Database")
		return
	}

	// cache
	if checkArg("-clear") {
		_, err := db.Exec(context.Background(), "TRUNCATE TABLE cache")
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Truncated Cache")
		return
	}

	// db insert
	if checkArg("-i") {
		msg := getRandomMessage()
		insertIntoDatabase(&msg)
		return
	}

	// nats
	if checkArg("-m") {
		publishMessage()
		return
	}

	if id, err := checkArgComplex("-id"); err != nil {
		log.Fatal(err)
	} else if id != 0 {
		row := db.QueryRow(context.Background(), "SELECT message FROM cache WHERE id ="+strconv.Itoa(id)+";")

		var tmp []byte
		err := row.Scan(&tmp)
		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile("data.json", tmp, 0644)
		return
	}

	serv, err := stan.Connect("test-cluster", "0")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connection OK")
	}
	defer serv.Close()
	_, err = serv.Subscribe("wb", messageCatcher, stan.DurableName("test"))
	if err != nil {
		log.Fatal(err)
	}

	for {

	}
}
