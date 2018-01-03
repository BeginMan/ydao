package main

import (
	"database/sql"
	"log"
	"path"
	"time"
	"encoding/csv"

	"github.com/apcera/termtables"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"os"
	"strconv"
)

type Word struct {
	id      int       `json:"id"`
	word    string    `json:"word"`
	en      bool      `json:"en"`
	created time.Time `json:"created"`
}

var (
	dbpath = path.Join(userHomeDir(), ".ydao.db")
)


func getConn() *sql.DB {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func queryWords(db *sql.DB) []Word {

	var result []Word

	rows, err := db.Query("SELECT * FROM words;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := Word{}
		rows.Scan(&item.id, &item.word, &item.en, &item.created)
		result = append(result, item)
	}
	return result
}

func prettyWords() {
	db := getConn()
	res := queryWords(db)
	table := termtables.CreateTable()
	table.AddHeaders("ID", "Word", "英文", "Created")

	for _, item := range res {
		table.AddRow(item.id, item.word, item.en, item.created)
	}
	db.Close()
	fmt.Println(table.Render())
}

func cleanHistory()  {
	db := getConn()
	db.Exec("DELETE FROM words")
}

func dumpHistory(){
	file, err := os.Create("history.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	db := getConn()
	res := queryWords(db)
	defer db.Close()
	for _, item := range res {
		w.Write([]string{item.word, strconv.FormatBool(item.en), item.created.String()})
	}
}

func init() {
	db := getConn()
	sql_table := `
    CREATE TABLE IF NOT EXISTS words(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word VARCHAR(250) NULL,
        en BOOLEAN NULL,
        created TimeStamp NOT NULL DEFAULT (datetime('now','localtime'))
    );
    `
	db.Exec(sql_table)
}

func insert(t Word, db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO words(word, en) values(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.word, t.en)
	if err != nil {
		return err
	}
	return nil
}

