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

var db *sql.DB

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

// Open a connection to our database and return a pointer to the instance of the DB connection
func openDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkDBErr(err)
	return db
}
