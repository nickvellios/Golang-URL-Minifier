// Golang URL Minifier
// Creates tiny URLS from a given URL.
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/23/2016

package main

import (
	"net/http"
	"log"
	_ "github.com/lib/pq"
)

var templateDir = "/root/go/tiny/bin/templates/"
var baseURL = "http://r8r.org/"

func main() {
	db = openDB()
	defer db.Close()

	// HTTP Routing
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/generate/", generateHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}