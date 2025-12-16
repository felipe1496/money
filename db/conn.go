package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func Conn(connString string) (*sql.DB, error) {
	dsn := connString

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
