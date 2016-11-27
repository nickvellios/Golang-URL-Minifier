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

type urlMapRow struct {
	path    string
	url     string
	ip      string
	t_stamp string
}

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
	rows, err := udb.db.Query("SELECT COUNT(*) as count, t_stamp::DATE as ts FROM url_map GROUP BY ts ORDER BY ts")
	checkDBErr(err)

	for rows.Next() {
		var t time.Time
		var count int
		err := rows.Scan(&count, &t)
		checkDBErr(err)
		const layout = "Jan 2, 2006"
		date := t.Format(layout)

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

	chart := &Charts{Cols: cols, Rws: rws}

	cht, _ := json.Marshal(chart)
	vars := map[string]interface{}{
		"js": template.HTML(string(cht)),
	}
	err = templates.ExecuteTemplate(w, "stats", vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
