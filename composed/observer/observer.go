package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var requestID uint64
var requestChannel = make(map[uint64](chan []string))

func main() {
	html, err := loadPage("observer.html")
	check(err)

	db := createDb()
	defer db.Close()

	stmt := prepreSelect(db)
	defer stmt.Close()

	go pollMessages(stmt)

	http.HandleFunc("/observe", func(w http.ResponseWriter, r *http.Request) {
		handleObserveGet(w, r, string(html))
	})
	http.HandleFunc("/poll", func(w http.ResponseWriter, r *http.Request) {
		handlePoll(w, r, stmt)
	})
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func loadPage(fileName string) ([]byte, error) {
	filename := fileName
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return html, nil
}

func handlePoll(w http.ResponseWriter, r *http.Request, stmt *sql.Stmt) {
	id := atomic.AddUint64(&requestID, 1)

	ch := make(chan []string, 100)
	defer close(ch)
	requestChannel[id] = ch

	fmt.Printf("Lock on polling for reqID %v\n", id)

	messages := <-ch

	for _, v := range messages {
		fmt.Fprintln(w, v)
	}

	fmt.Printf("Unlock on polling for reqID %v\n", id)
	requestChannel[id] = nil
}

func handleObserveGet(w http.ResponseWriter, r *http.Request, html string) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
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

func prepreSelect(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT text FROM message WHERE created_at >= $1 ORDER BY created_at DESC")
	check(err)
	return stmt
}

func execQuery(stmt *sql.Stmt, olderThan *time.Time) []string {
	rows, err := stmt.Query(olderThan)
	check(err)
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var m string
		err := rows.Scan(&m)
		check(err)
		messages = append(messages, m)
	}
	check(rows.Err())
	return messages
}

func pollMessages(stmt *sql.Stmt) {
	observeOlderThan := time.Now()
	var messages []string

	for {
		messages = nil
		for messages == nil {
			time.Sleep(50 * time.Microsecond)
			messages = execQuery(stmt, &observeOlderThan)
		}

		for k, v := range requestChannel {
			fmt.Printf("Lock on DB polling for reqID %v\n", k)
			if v != nil {
				v <- messages
			}
			fmt.Printf("Unlock on DB polling for reqID %v\n", k)
		}

		observeOlderThan = time.Now()
	}
}
