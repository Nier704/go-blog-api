package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
	Date   string `json:"date"`
}

type MessageDTO struct {
	Author string `json:"author"`
	Text   string `json:"text"`
	Date   string `json:"date"`
}

type MessageList []Message

var db *sql.DB

func Router(r chi.Router) {
	if err := openDBConn(); err != nil {
		log.Fatal(err)
	}

	r.Get("/", allMessages)
	r.Post("/", sendMessage)
}

func allMessages(w http.ResponseWriter, r *http.Request) {
	results, err := db.Query("SELECT * FROM message")
	if err != nil {
		log.Fatal(err)
	}

	var msgs MessageList

	for results.Next() {
		var msg Message
		if err = results.Scan(&msg.Id, &msg.Author, &msg.Text, &msg.Date); err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, msg)
	}

	json.NewEncoder(w).Encode(msgs)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var new_msg MessageDTO
	err := json.NewDecoder(r.Body).Decode(&new_msg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO message (author, text, date) VALUES ('%s', '%s', '%s')", new_msg.Author, new_msg.Text, new_msg.Date))
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(new_msg)
}

func openDBConn() (err error) {
	db, err = sql.Open("mysql", "root:srlucky123@tcp(127.0.0.1:3306)/")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS blogdb")
	if err != nil {
		return err
	}

	_, err = db.Exec("USE blogdb")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS message (id INTEGER PRIMARY KEY AUTO_INCREMENT, author VARCHAR(255) NOT NULL, text VARCHAR(255) NOT NULL, date VARCHAR(255) NOT NULL)")
	if err != nil {
		return err
	}

	fmt.Println("database connection established")
	return nil
}
