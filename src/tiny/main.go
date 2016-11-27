// Golang URL Minifier
// Creates tiny URLS from a given URL.
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/23/2016

package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

var templateDir = "/root/go/tiny/bin/templates/"
var baseURL = "http://r8r.org"

const (
	readTimeout  = time.Duration(1 * time.Second)
	writeTimeout = readTimeout
)

func main() {
	gdb := &urlDB{}
	gdb.open()
	defer gdb.db.Close()

	// HTTP Routing
	http.HandleFunc("/", gdb.rootHandler)
	http.HandleFunc("/generate/", gdb.generateHandler)
	http.HandleFunc("/stats/", gdb.statsHandler)

	srv := http.Server{
		Addr:         ":80",
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
