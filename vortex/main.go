package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// подключение к базе данных
var g_db *pgx.Conn

// подключение к интерфейсу Deployer
var g_deploy Deployer

func main() {
	const db_url = ""

	db, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer db.Close(context.Background())

	g_db = db
	g_deploy = NewDeployer()

	for {
		Update()
		time.Sleep(time.Minute * 5)
	}
}
