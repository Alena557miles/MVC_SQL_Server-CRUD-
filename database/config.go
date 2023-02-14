package database

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   `root`,
		Passwd: `password`,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "my-app",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func PingDB(db *sql.DB) {
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
