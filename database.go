package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func connectToDatabase() *pgx.Conn {
	dburl, err := os.ReadFile("auth.txt")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(context.Background(), string(dburl))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

const (
	d = 0 //drop
	c = 1 //create
	r = 2 //refresh
)

func populateDatabase(action int) {

	var sql [5]string

	sql[0] = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		phone VARCHAR(20),
		zip VARCHAR(20),
		city VARCHAR(50),
		address VARCHAR(100),
		region VARCHAR(50),
		email VARCHAR(50)
	);`
	sql[1] = `CREATE TABLE IF NOT EXISTS transactions (
		transaction VARCHAR(100) PRIMARY KEY,
		request_id VARCHAR(20),
		currency VARCHAR(20),
		provider VARCHAR(20),
		amount INT,
		payment_dt INT,
		bank VARCHAR(20),
		delivery_cost INT,
		goods_total INT,
		custom_fee INT
	);`
	sql[2] = `CREATE TABLE IF NOT EXISTS items (
		chrt_id INT PRIMARY KEY,
		track_number VARCHAR(100),
		price INT,
		rid VARCHAR(100),
		name VARCHAR(100),
		sale INT,
		size VARCHAR(20),
		total_price INT,
		nm_id INT,
		brand VARCHAR(100),
		status INT
	);`
	sql[3] = `CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(100) PRIMARY KEY,
		track_number VARCHAR(100),
		entry VARCHAR(20),
		delivery INT,
		payment VARCHAR(100),
		items INT [],
		locale VARCHAR(10),
		internal_signature VARCHAR(10),
		customer_id VARCHAR(20),
		delivery_service VARCHAR(20),
		shardkey VARCHAR(10),
		sm_id INT,
		date_created TIMESTAMP,
		FOREIGN KEY(delivery) REFERENCES users(id),
		FOREIGN KEY(payment) REFERENCES transactions(transaction)
	);`
	sql[4] = `CREATE TABLE IF NOT EXISTS cache (
		id SERIAL PRIMARY KEY,
		message JSON
	);`
	// !!
	//FOREIGN KEY(EACH ELEMENT OF items) REFERENCES items(chrt_id)

	if action == d || action == r {
		_, err := db.Exec(context.Background(), "TRUNCATE TABLE transactions, users, orders, items;")
		if err != nil {
			log.Fatal(err)
		}
	}
	if action == c || action == r {
		for i := range sql {
			_, err := db.Exec(context.Background(), sql[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func insertIntoDatabase(order *Order) {

	// insert user
	user := order.Delivery
	row := db.QueryRow(context.Background(), `INSERT INTO users(
		name,
		phone,
		zip,
		city,
		address,
		region,
		email)

		VALUES(`+
		"'"+user.Name+"', "+
		"'"+user.Phone+"', "+
		"'"+user.Zip+"', "+
		"'"+user.City+"', "+
		"'"+user.Address+"', "+
		"'"+user.Region+"', "+
		"'"+user.Email+"') RETURNING id;")

	var user_id int // used in inserting order
	err := row.Scan(&user_id)
	if err != nil {
		log.Fatal(err)
	}

	// insert transaction
	t := order.Payment
	_, err = db.Exec(context.Background(), `INSERT INTO transactions(
		transaction,
		request_id,
		currency,
		provider,
		amount,
		payment_dt,
		bank,
		delivery_cost,
		goods_total,
		custom_fee)
		
		VALUES(`+
		"'"+t.Transaction+"', "+
		"'"+t.Request_id+"', "+
		"'"+t.Currency+"', "+
		"'"+t.Provider+"', "+
		strconv.Itoa(t.Amount)+", "+
		strconv.FormatInt(t.Payment_dt, 10)+", "+
		"'"+t.Bank+"', "+
		strconv.Itoa(t.Delivery_cost)+", "+
		strconv.Itoa(t.Goods_total)+", "+
		strconv.Itoa(t.Custom_fee)+");")

	if err != nil {
		log.Fatal(err)
	}

	// insert items
	for _, item := range order.Items {
		_, err := db.Exec(context.Background(), `INSERT INTO items(
		chrt_id,
		track_number,
		price,
		rid,
		name,
		sale,
		size,
		total_price,
		nm_id,
		brand,
		status)
		
		VALUES(`+
			strconv.Itoa(item.Chrt_id)+", "+
			"'"+item.Track_number+"', "+
			strconv.Itoa(item.Price)+", "+
			"'"+item.Rid+"', "+
			"'"+item.Name+"', "+
			strconv.Itoa(item.Sale)+", "+
			"'"+item.Size+"', "+
			strconv.Itoa(item.Total_price)+", "+
			strconv.Itoa(item.Nm_id)+", "+
			"'"+item.Brand+"', "+
			strconv.Itoa(item.Status)+");")

		if err != nil {
			log.Fatal(err)
		}
	}

	orderSql := `INSERT INTO orders(
		order_uid,
		track_number,
		entry,
		delivery,
		payment,
		items,
		locale,
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created)

		VALUES(` +
		"'" + order.Order_uid + "', " +
		"'" + order.Track_number + "', " +
		"'" + order.Entry + "', " +
		strconv.Itoa(user_id) + ", " +
		"'" + t.Transaction + "', " + "ARRAY["

	for i, item := range order.Items {
		orderSql += strconv.Itoa(item.Chrt_id)
		if i != len(order.Items)-1 {
			orderSql += ", "
		}
	}
	orderSql += "],"

	_, err = db.Exec(context.Background(), orderSql+
		"'"+order.Locale+"', "+
		"'"+order.Internal_signature+"', "+
		"'"+order.Customer_id+"', "+
		"'"+order.Delivery_service+"', "+
		"'"+order.Shardkey+"', "+
		strconv.Itoa(order.Sm_id)+", "+
		"'"+order.Date_created+"');")
	if err != nil {
		log.Fatal(err)
	}
}

func restoreCache() {
	rows, err := db.Query(context.Background(), "SELECT message FROM cache ORDER BY id DESC LIMIT "+strconv.Itoa(CACHE_RESTORE_LIMIT)+";")
	if err != nil {
		log.Fatal(err)
	}

	var tmp []byte
	for rows.Next() {
		rows.Scan(&tmp)
		cache = append(cache, tmp)
	}
}

func getRandomMessage() Order {
	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalf("getRandomMessage: %v", err)
	}
	var msg Order
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Fatalf("getRandomMessage: %v", err)
	}
	msg.Order_uid = randomString(12)
	msg.Track_number = randomString(7)
	msg.Items[0].Chrt_id = randomInt(1000000000)
	if rand.Intn(2) == 1 {
		msg.Items = append(msg.Items, msg.Items[0])
		msg.Items[1].Chrt_id = randomInt(1000000)
	}
	msg.Payment.Transaction = randomString(10)
	msg.Delivery.Name = randomString(6) + " " + randomString(6)
	msg.Date_created = time.Now().Format("2006-01-02T15:04:05Z07")
	return msg
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(l int) string {
	rand.Seed(time.Now().Unix())

	s := make([]rune, l)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomInt(l int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(l)
}
