package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func transaction(db *sql.DB) {
	// query := `
	// CREATE TABLE IF NOT EXISTS accounts (
	//     user_id SERIAL PRIMARY KEY,
	//     balance INT
	// )`
	// db.Exec(query)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO accounts (user_id, balance) VALUES ($1, $2)", 1, 100)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", 50, 1)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE user_id = $2", 50, 2)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Транзакція успішно виконана")
}
