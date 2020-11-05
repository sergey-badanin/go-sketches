package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	runWebServer()
}

func testSave() {
	db := createDb()
	defer db.Close()

	saveMessage("Message 01", db)
}

func runWebServer() {

	html, err := loadPage("send.html")
	check(err)

	db := createDb()
	check(err)
	defer db.Close()

	msgSaver := func(msg string) {
		saveMessage(msg, db)
	}

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		sendHandler(w, r, string(html), msgSaver)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendHandler(w http.ResponseWriter, r *http.Request, html string, save func(string)) {
	switch method := r.Method; method {
	case "GET":
		handleSendGet(w, r, html)
	case "POST":
		handleSendPost(w, r, save)
	}
}

func handleSendGet(w http.ResponseWriter, r *http.Request, html string) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func handleSendPost(w http.ResponseWriter, r *http.Request, save func(string)) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	body := r.Form["body"]

	var resultBuffer bytes.Buffer
	for _, v := range body {
		if len(v) > 0 {
			resultBuffer.WriteString(" " + v)
		}
	}

	msg := resultBuffer.String()
	fmt.Printf("String:%v, length: %v\n", msg, len(msg))
	save(msg)

	http.Redirect(w, r, "/send", http.StatusFound)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadPage(fileName string) ([]byte, error) {
	filename := fileName
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return html, nil
}

func saveMessage(msg string, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO message(created_at, text) VALUES ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now(), msg)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func createDb() *sql.DB {
	dataSource := "postgres://postgres:postgresql@localhost:5432/composed"
	db, err := sql.Open("pgx", dataSource)
	check(err)
	err = db.Ping()
	check(err)
	return db
}

func runDbQuery() {
	db, err := sql.Open("pgx", "postgres://postgres:postgresql@localhost:5432/db01")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	var col1 string
	var col2 string
	err = db.QueryRow("select * from t1").Scan(&col1, &col2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(col1, col2)
}
