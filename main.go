package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Person struct {
	ID         string
	Name       string `json:"name"`
	Age        int64  `json:"age"`
	Occupation string `json:"occupation"`
}

func main() {
	// Setup a connection to postgres
	connStr := "postgres://postgres:pass@localhost:5432/gopgtest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	// Create a database
	createProductTable(db)

	fmt.Println("Redis trial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ping)

	personID := uuid.NewString()
	jsonStirng, err := json.Marshal(Person{
		ID:         personID,
		Name:       "Abdo",
		Age:        30,
		Occupation: "SWE",
	})
	if err != nil {
		fmt.Printf("failed to marshal: %s", err.Error())
	}
	personKey := fmt.Sprintf("person:%s", personID)
	err = client.Set(context.Background(), personKey, jsonStirng, 0).Err()
	if err != nil {
		fmt.Printf("Failed to set value in redis instance: %s\n", err.Error())
	}

	val, err := client.Get(context.Background(), personKey).Result()
	if err != nil {
		fmt.Printf("failed to get value from redis: %s", err.Error())
	}

	fmt.Printf("value retrieved = %s", val)
}

func createProductTable(db *sql.DB) {
	qurey := `
	create table if not exists product (
		id serial primary key,
		name varchar(100) not null,
		price numeric(6,2) not null,
		available boolean,
		created timestamp default NOW()
	)`

	_, err := db.Exec(qurey)
	if err != nil {
		log.Fatal(err)
	}
}
