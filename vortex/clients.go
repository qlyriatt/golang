package main

import (
	"context"
	"time"
)

type Client struct {
	ID          int64
	ClientName  string
	Version     int
	Image       string
	CPU         string
	Memory      string
	Priority    float64
	NeedRestart bool
	SpawnedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Функции для взаимодействия с БД посредством структуры Client

func addClient(c *Client) {
	g_db.Exec(context.Background(),
		`INSERT INTO clients (
		client_id,
		VWAP,
		TWAP,
		HFT)
	
	VALUES (
		$1
		DEFAULT,
		DEFAULT,
		DEFAULT)`, c.ID)
}

func UpdateClient(c *Client) {
	g_db.Exec(context.Background(),
		`UPDATE clients SET
		VWAP = DEFAULT,
		TWAP = DEFAULT,
		HFT = DEFAULT

		WHERE client_id = $1`, c.ID)
}

func DeleteClient(c *Client) {
	g_db.Exec(context.Background(),
		`DELETE FROM clients
		
		WHERE client_id = $1`, c.ID)
}

func updateAlgorithmStatus(c *Client, algo string, on bool) {
	g_db.Exec(context.Background(),
		`UPDATE clients SET
		$1 = $2
		
		WHERE client_id = $3`, algo, on, c.ID)
}
