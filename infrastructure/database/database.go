package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewConnectionPostgres(host, port, user, password, name string) (*sql.DB, error) {
	connSTR := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name,
	)

	db, err := sql.Open("postgres", connSTR)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("SET TIME ZONE 'UTC'")
	if err != nil {
		return nil, err
	}

	var timezone string
	err = db.QueryRow("SHOW TIME ZONE").Scan(&timezone)
	if err != nil {
		log.Printf("Warning: Error getting timezone: %s", err)
	} else if timezone != "UTC" {
		log.Printf("Warning: Timezone is not UTC: %s", timezone)
	}

	return db, nil
}
