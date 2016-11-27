// Golang URL Minifier
// Creates tiny URLS from a given URL.
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/23/2016

package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// The following 4 structs cumulatively define our JSON payload's data structure
type Charts struct {
	Cols []Columns `json:"cols"`
	Rws  []Rows    `json:"rows"`
}

type Columns struct {
	Label      string `json:"label"`
	ColumnType string `json:"type"`
}

type Rows struct {
	C []C2 `json:"c"`
}

type C2 struct {
	V string `json:"v"`
	F string `json:"f"`
}

// HTTP handler for /stats/ path
func (udb *urlDB) statsHandler(w http.ResponseWriter, r *http.Request) {
	var rws []Rows
	var cols = []Columns{
		{
			Label:      "X",
			ColumnType: "date",
		},
		{
			Label:      "Daily Total",
			ColumnType: "number",
		},
	}

	// Select all records from the database, split the datetime into a date and group by days.  Count the number of requests per day.
	// This needs to be changed soon to allow a date range, and possibly to also to fill in empty days if they exist.
	rows, err := udb.db.Query("SELECT COUNT(*) as count, t_stamp::DATE as ts FROM url_map GROUP BY ts ORDER BY ts")
	checkDBErr(err)

	for rows.Next() {
		var t time.Time
		var count int
		err := rows.Scan(&count, &t)
		checkDBErr(err)
		const layout = "Jan 2, 2006"
		date := t.Format(layout)

		// Convert the date into a string for the JS structure:  "Date(YYYY,M,D)".
		// Javascript months are zero indexed (stupid) so we must subtract one.  Thanks to my lovely wife for catching that bug
		datestr := fmt.Sprintf("Date(%d,%d,%d)", t.Year(), t.Month()-1, t.Day())
		dailytotal := strconv.FormatInt(int64(count), 10)

		rws = append(rws, Rows{
			C: []C2{
				{
					V: datestr,
					F: date,
				},
				{
					V: dailytotal,
					F: dailytotal,
				},
			},
		})
	}

	// Finish building our data structure
	chart := &Charts{Cols: cols, Rws: rws}

	// Encode into JSON and pass to our stats.html template
	cht, _ := json.Marshal(chart)

	// Build template data structure
	p := &Page{
		Title: "Stats",
		Content: struct {
			JS interface{}
		}{
			template.HTML(string(cht)),
		},
	}

	renderTemplate(w, "stats", p)
}
