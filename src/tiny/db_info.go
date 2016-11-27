// Golang URL Minifier
// Database Configuration
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/14/2016

/*
POSTGRESQL Schema

CREATE TABLE url_map (
	id SERIAL PRIMARY KEY,
	path VARCHAR(16),
	url TEXT,
	t_stamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	ip VARCHAR(39)
);
*/

package main

import (
	"database/sql"
	"fmt"
)

type urlDB struct {
	db     *sql.DB
	APIKey string
}

const APIKey = "1234567890ABCD" // Used for an internal system to allow limitless URLs to be processed

// Database configuration

const (
	DB_USER     = "username"
	DB_PASSWORD = "password"
	DB_NAME     = "database"
)

// Global DB error checking
func checkDBErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (udb *urlDB) open() error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	udb.db, err = sql.Open("postgres", dbinfo)
	checkDBErr(err)
	return err
}
