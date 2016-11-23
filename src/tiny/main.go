// Golang URL Minifier
// Creates tiny URLS from a given URL.
// https://github.com/nickvellios/Golang-URL-Minifier
// Nick Vellios
// 11/23/2016

package main

import (
	"fmt"
	"math"
	"net/http"
	"log"
	"net"
	"html/template"
	"encoding/json"
	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseFiles(
	"./bin/templates/index.html",
	"./bin/templates/header.html",
	"./bin/templates/footer.html"))

type Tiny struct {
	URL string
	Path string
	IP string
	Timestamp string
	ID string
}

type TinyJS struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

func init() {
}

func main() {
	db = openDB()
	defer db.Close()

	// HTTP Routing
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/generate/", generateHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

// Send a JSON result back to the client with given status code
func writeResponse(w http.ResponseWriter, code int, url, error string) error {
	pkg := map[string]string{"url": url, "error": error}
	js, err := json.Marshal(pkg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)

	return err
}

func renderTemplate(w http.ResponseWriter, tmpl string, data []map[string]string) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("User IP is not IP:port", r.RemoteAddr)
	}

	tiny := &Tiny{URL: url, IP: ip}

	if tiny.throttleCheck() {
		tiny.save()
		writeResponse(w, 200, "http://r8r.org/" + tiny.Path, "")
	} else {
		writeResponse(w, 200, "", "You're doing that too often.  Slow down")
	}
}

func (t *Tiny) save() int {
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

func (t *Tiny) load() {
	rows, err := db.Query("SELECT url FROM url_map WHERE path = $1", t.Path)
	checkDBErr(err)

	for rows.Next() {
		err := rows.Scan(&t.URL)
		checkDBErr(err)
	}
}

// Check if user has used service more than 10 times in an hour.
func (t *Tiny) throttleCheck() bool {
	rows, err := db.Query("SELECT COUNT(*) as count FROM url_map WHERE ip = $1 AND t_stamp > CURRENT_TIMESTAMP - INTERVAL '1 hour'", t.IP)
	checkDBErr(err)

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		checkDBErr(err)
	}

	return count<10
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/"):]

	if len(url) > 0 && url != "favicon.ico" {
		t := &Tiny{Path: url}
		t.load()
		fmt.Println("Redirecting to: ", t.URL)
		http.Redirect(w, r, t.URL, 302)

		return
	}

	renderTemplate(w, "index", nil)
}

func generateCode(number int) string {
	var out []byte
	codes := []byte("abcdefghjkmnpqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ");

	for number > 53 {
		key := number % 54
		number = int(math.Floor(float64(number) / 54) - 1)
		out = append(out, []byte(codes[key:key+1])[0])
	}

	return string(append(out, codes[number]))
}