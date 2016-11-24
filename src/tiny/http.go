// Golang URL Minifier
// Creates tiny URLS from a given URL.
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/23/2016

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"math"
	"net"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles(
	templateDir+"index.html",
	templateDir+"header.html",
	templateDir+"footer.html"))

type Tiny struct {
	URL       string
	Path      string
	IP        string
	Timestamp string
	ID        string
}

type URLResponse struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// There may be URLs longer, but to avoid attacks we don't want them.
const (
	MAX_URL = 1024
)

// Send a JSON result back to the client with given status code
func writeResponse(w http.ResponseWriter, code int, url, error string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	ur := &URLResponse{url, error}
	err := json.NewEncoder(w).Encode(ur)

	return err
}

func renderTemplate(w http.ResponseWriter, tmpl string, data []map[string]string) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HTTP handler for /generate/ which is the API.  Takes a long URL and generates a tiny URL.
func (udb *urlDB) generateHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	// Ignore super long URLs
	if len(url) > MAX_URL {
		writeResponse(w, 413, "", "URL exceeds maximum length of 1024 characters")
		return
	}

	// 301/302 redirects fail without a valid URL.  Our site frondend checks for this and adds a http:// prefix if needed, but we will check and do the same for API requests
	reg, _ := regexp.Compile(`^(http|https|ftp)+(://)`)
	if !reg.MatchString(url) {
		url = "http://" + url
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("User IP is not IP:port", r.RemoteAddr)
	}

	tiny := &Tiny{URL: url, IP: ip}

	// Limit to 10 requests per hour per IP
	if tiny.throttleCheck(udb.db) {
		tiny.save(udb.db)
		writeResponse(w, 200, baseURL+tiny.Path, "")
	} else {
		writeResponse(w, 429, "", "You're doing that too often.  Slow down")
	}
}

// HTTP handler for / path
func (udb *urlDB) rootHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/"):]

	if len(url) > 0 && url != "favicon.ico" {
		t := &Tiny{Path: url}
		t.load(udb.db)
		fmt.Println("Redirecting to: ", t.URL)
		http.Redirect(w, r, t.URL, 302)

		return
	}

	renderTemplate(w, "index", nil)
}

// Save URL to DB, get unique ID, generate tiny path from the ID, update the DB.
func (t *Tiny) save(db *sql.DB) int {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO url_map(path, url, ip) VALUES($1,$2,$3) returning id;", "", t.URL, t.IP).Scan(&lastInsertId)
	checkDBErr(err)

	stmt, err := db.Prepare("UPDATE url_map SET path=$1 WHERE id=$2")
	checkDBErr(err)
	path := generateCode(lastInsertId)
	t.Path = path
	_, err = stmt.Exec(t.Path, lastInsertId)
	checkDBErr(err)

	return lastInsertId
}

// Load URL from DB.
func (t *Tiny) load(db *sql.DB) {
	rows, err := db.Query("SELECT url FROM url_map WHERE path = $1", t.Path)
	checkDBErr(err)

	for rows.Next() {
		err := rows.Scan(&t.URL)
		checkDBErr(err)
	}
}

// Check if user has used service more than 10 times in an hour.
func (t *Tiny) throttleCheck(db *sql.DB) bool {
	rows, err := db.Query("SELECT COUNT(*) as count FROM url_map WHERE ip = $1 AND t_stamp > CURRENT_TIMESTAMP - INTERVAL '1 hour'", t.IP)
	checkDBErr(err)

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		checkDBErr(err)
	}

	return count < 10
}

// Generate a unique A-Z/0-9 based URL from a given number input.
func generateCode(number int) string {
	var out []byte
	codes := []byte("abcdefghjkmnpqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ")

	for number > 53 {
		key := number % 54
		number = int(math.Floor(float64(number)/54) - 1)
		out = append(out, []byte(codes[key : key+1])[0])
	}

	return string(append(out, codes[number]))
}
